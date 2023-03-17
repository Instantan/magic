package magic_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/Instantan/magic"
)

func TestTemplate(t *testing.T) {
	template, err := magic.ParseFile("example/index.html")
	if err != nil {
		t.Error(err)
	}
	template.Apply(map[string]string{
		"name": "Felix",
	})

	buf := new(strings.Builder)
	template.Write(buf)

	t.Log(buf.String())
}

func TestHandler(t *testing.T) {

	template := magic.Must(magic.ParseFile("example/index.html"))

	mux := http.NewServeMux()

	mux.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {

		template.Execute(w, map[string]string{
			"name": "Felix",
		})
	})

}
