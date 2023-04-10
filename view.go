package magic

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type view struct {
	rawTemplate  *MagicTemplate
	liveTemplate *MagicTemplate
}

func View(template string) *view {
	v := view{}

	template = strings.TrimSpace(template)
	template = strings.ReplaceAll(template, "\n", "")
	template = strings.ReplaceAll(template, "\r", "")
	template = strings.ReplaceAll(template, "\t", "")

	v.rawTemplate = CompileTemplate(template)
	v.liveTemplate = CompileTemplate(buildLiveTemplateFromRawTemplate(template))

	return &v
}

func (v *view) Render(data any) string {
	v.rawTemplate.viewResolver = func(name string) string {
		return "bla"
	}
	return v.rawTemplate.Exec(data)
}

func (v *view) String() string {
	return fmt.Sprintf("Raw: %v\nLive: %v", v.rawTemplate.template.Root, v.liveTemplate.template.Root)
}

func buildLiveTemplateFromRawTemplate(template string) string {
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
