package magic_test

import (
	"fmt"
	"testing"

	"github.com/Instantan/magic"
)

func TestView(t *testing.T) {

	testView := magic.View(`
		<h1 class="{{ .vis }}">
			{{range .posts}}
				{{.}}
			{{end}}
			<p>Test</p>
		</h1>
	`)

	// data := testView.Render(map[string]any{
	// 	"vis":   "blabla",
	// 	"posts": []int{1, 2, 3},
	// 	"view2": view2.String(),
	// 	// "magicView": func() string {
	// 	// 	return "bla"
	// 	// },
	// })

	fmt.Printf("%v", testView.String())
}
