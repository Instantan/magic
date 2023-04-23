package magic

import "testing"

func TestParseTemplate(t *testing.T) {
	n := View[any](`
		<h1>1</h1>
		<h2 class="{{bla}}">
			<p>bla</p>
			<p>bla<p>

			{{ range data }}
				{{ . }}
			{{ end }}
		</h2>
	`)(nil)
	t.Log(n.HTML())
}
