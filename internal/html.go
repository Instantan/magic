package internal

import (
	"strings"

	"golang.org/x/net/html"
)

type ByteInjector func() []byte

func InjectDataIntoHTML(data []byte, injectHTMLAttributes, injectHeaderChilds ByteInjector) []byte {
	saw := []byte{}
	injectHTMLAttributesAt := 0
	sawHTML := false
	injectHeaderChildsAt := 0
	sawHeader := false
	for i := range data {
		saw = append(saw, byte(data[i]))
		if !sawHTML && sawHTMLOpenTag(saw) {
			sawHTML = true
			injectHTMLAttributesAt = i + 1
		}
		if !sawHeader && sawHeadTag(saw) {
			sawHeader = true
			injectHeaderChildsAt = i + 1
		}
		if sawHTML && sawHeader {
			toInjectHTMLAttributes := injectHTMLAttributes()
			toInjectHeaderChils := injectHeaderChilds()
			data = []byte(string(data[:injectHTMLAttributesAt]) + " " + string(toInjectHTMLAttributes) + string(data[injectHTMLAttributesAt:injectHeaderChildsAt]) + string(toInjectHeaderChils) + string(data[injectHeaderChildsAt:]))
			return data
		}
	}
	return data

}

func sawHTMLOpenTag(saw []byte) bool {
	needsToSee := "<html"
	return len(saw) >= len(needsToSee) && strings.HasSuffix(strings.ToLower(string(saw)), needsToSee)
}

func sawHeadTag(saw []byte) bool {
	needsToSee := "<head>"
	return len(saw) >= len(needsToSee) && strings.HasSuffix(strings.ToLower(string(saw)), needsToSee)
}

func BuildLiveTemplateFromRawTemplate(template string) string {
	tkn := html.NewTokenizer(strings.NewReader(template))
	result := strings.Builder{}

	for {
		tt := tkn.Next()
		switch tt {
		case html.ErrorToken:
			return result.String()
		case html.StartTagToken:
			t := tkn.Token()
			magicAttrs := []html.Attribute{}
			for i := range t.Attr {
				attr := t.Attr[i]
				if strings.HasPrefix(attr.Key, "magic-") || !strings.Contains(attr.Val, "{{") {
					continue
				}
				magicAttrs = append(magicAttrs, html.Attribute{
					Namespace: attr.Namespace,
					Key:       "magic-" + attr.Key,
					Val:       convTemplateIntoMagicTemplate(attr.Val),
				})
			}
			t.Attr = append(t.Attr, magicAttrs...)
			result.WriteString(t.String())
		case html.TextToken:
			t := tkn.Token()
			if strings.Contains(t.Data, "{{") {
				result.WriteString(convTextTokenIntoMagicValue(t))
				continue
			}
			result.WriteString(t.String())
		default:
			t := tkn.Token()
			result.WriteString(t.String())
		}
	}
}

func convTextTokenIntoMagicValue(t html.Token) string {
	data := strings.TrimSpace(t.String())
	return "<m magic-value=\"" + convTemplateIntoMagicTemplate(data) + "\">" + data + "</m>"
}

func convTemplateIntoMagicTemplate(template string) string {
	return strings.ReplaceAll(strings.ReplaceAll(template, "{{", "{§"), "}}", "§}")
}
