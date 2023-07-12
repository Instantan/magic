package magic

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"regexp"
	"strings"
)

//go:embed script/magic.min.js
var magicMinScript []byte
var magicMinScriptWithTags []byte
var regexpHeadTag = regexp.MustCompile("<head.*>")

func init() {
	magicMinScriptWithTags = append(append([]byte("<script magic:inject>"), magicMinScript...), []byte("</script>")...)
	{
		buf := &bytes.Buffer{}
		writer, err := gzip.NewWriterLevel(buf, gzip.BestCompression)
		if err != nil {
			panic(err)
		}
		if _, err = writer.Write(magicMinScript); err != nil {
			panic(err)
		}
		if err = writer.Close(); err != nil {
			panic(err)
		}
	}
}

func injectLiveScript(templ string) string {
	return regexpHeadTag.ReplaceAllStringFunc(templ, func(s string) string {
		s = strings.Replace(s, ">", ">{{magic:inject}}", 1)
		return s
	})
}
