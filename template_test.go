package magic_test

import (
	"log"
	"testing"

	"github.com/Instantan/magic"
)

func TestCompileTemplate(t *testing.T) {
	tmpl := magic.CompileTemplate(`hello {{range .bla}}{{.}}haha{{end}}.`)
	res := tmpl.Exec(map[string]any{
		"bla": []int{1, 2, 3, 4},
	})
	log.Printf("%#v", res)
}
