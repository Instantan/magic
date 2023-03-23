package internal

import "testing"

func TestInjectDataIntoHTML(t *testing.T) {
	html := []byte(`
		<html lang="de">
		<header>

		</header>
		jkflawndajkl
		<kjlmndjeklfm
	`)

	res := InjectDataIntoHTML(html, func() []byte {
		return []byte("BLA")
	}, func() []byte {
		return []byte("BLAAAAA")
	})
	t.Log(string(res))
}

func TestReplaceTemplateBracesInHTMLInnerTextWithComponent(t *testing.T) {

	html := `
		<html lang="de">
		<header {{bla}}>
		{{blub}}
		</header>
		jkflawndajkl
		{{blub}}
		<kjlmndjeklfm
	`

	nhtml := ReplaceTemplateBracesInHTMLInnerTextWithComponent(html)
	t.Log(nhtml)
}
