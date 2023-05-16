package magic

import "testing"

func TestParseTemplate(t *testing.T) {
	n := View(`
		<h1>1</h1>
		<h2 class="{{bla}}">
			<p>bla</p>
			<p>bla<p>
		</h2>
	`)(nil)
	t.Log(n)
}
