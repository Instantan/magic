package magic

import "strings"

func injectDataIntoHTML(data []byte, injectHTMLAttributes func() []byte, injectHeaderChilds func() []byte) []byte {
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
		if !sawHeader && sawHeaderTag(saw) {
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

func sawHeaderTag(saw []byte) bool {
	needsToSee := "<head>"
	return len(saw) >= len(needsToSee) && strings.HasSuffix(strings.ToLower(string(saw)), needsToSee)
}
