package magic

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"net/http"
	"regexp"
	"strings"
)

//go:embed script/magic.min.js
var magicMinScript []byte
var magicMinScriptWithTags []byte
var gzippedMagicMinScript []byte
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
		gzippedMagicMinScript = buf.Bytes()
	}
}

func injectLiveScript(templ string) string {
	return regexpHeadTag.ReplaceAllStringFunc(templ, func(s string) string {
		s = strings.Replace(s, ">", ">{{magic:inject}}", 1)
		return s
	})
}

func ServeMagicScript(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/javascript")
	w.Header().Set("Cache-Control", "public, max-age=604800, immutable")
	if acceptsGzip(r) {
		w.Header().Add("Content-Encoding", "gzip")
		w.Header().Add("Vary", "Accept-Encoding")
		w.Write(gzippedMagicMinScript)
		return
	}
	w.Write(magicMinScript)
}

func injectedScripts(urls []string) string {
	s := ""
	for _, url := range urls {
		s = s + "<script magic:inject src=\"" + url + "\"/>"
	}
	return s
}

func injectedInlineScripts(inline []string) string {
	s := ""
	for _, inline := range inline {
		s = s + "<script magic:inject>" + inline + "</script>"
	}
	return s
}
