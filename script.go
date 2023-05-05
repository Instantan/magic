package magic

import (
	_ "embed"
	"regexp"
	"strings"
)

//go:embed script/magic.min.js
var magicMinScript []byte

func init() {
	magicMinScript = append(append([]byte("<script>"), magicMinScript...), []byte("</script>")...)
}

var regexpHeadTag = regexp.MustCompile("<head.*>")

func injectLiveScript(templ string) string {
	return regexpHeadTag.ReplaceAllStringFunc(templ, func(s string) string {
		s = strings.Replace(s, ">", ">{{magic:live}}", 1)
		return s
	})
}
