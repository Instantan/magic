package internal

import (
	"strings"
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

func ReplaceTemplateBracesInHTMLInnerTextWithComponent(data string) string {
	isInInnerText := false
	escaped := false
	sawTag := false

	b := strings.Builder{}
	b.Grow(len(data))

	tagContent := ""

	for i := 0; i < len(data); i++ {
		if escaped {
			escaped = false
			b.WriteByte(data[i])
			continue
		}
		switch data[i] {
		case '<':
			isInInnerText = false
		case '>':
			isInInnerText = true
		case '\\':
			escaped = true
		case '{':
			if isInInnerText && len(data)-1 > i+1 && data[i+1] == '{' {
				sawTag = true
				tagContent = ""
				i++
				continue
			}
		case '}':
			if isInInnerText && sawTag && len(data)-1 > i+1 && data[i+1] == '}' {
				sawTag = false
				b.WriteString("<m magic-value=\"")
				b.WriteString(tagContent)
				b.WriteString("\">{{")
				b.WriteString(tagContent)
				b.WriteString("}}</m>")
				i++
				continue
			}
		}
		if isInInnerText && sawTag {
			tagContent += string(data[i])
		} else {
			b.WriteByte(data[i])
		}
	}

	return b.String()
}
