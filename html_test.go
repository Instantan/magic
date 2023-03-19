package magic

import "testing"

func TestInjectDataIntoHTML(t *testing.T) {
	html := []byte(`
		<html lang="de">
		<header>

		</header>
		jkflawndajkl
		<kjlmndjeklfm
	`)

	res := injectDataIntoHTML(html, func() []byte {
		return []byte("BLA")
	}, func() []byte {
		return []byte("BLAAAAA")
	})
	t.Log(string(res))
}
