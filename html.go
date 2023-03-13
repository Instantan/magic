package magic

import (
	"bytes"
	"io/fs"

	"golang.org/x/net/html"
)

func readAndTokenizeHTMLFile(fsys fs.FS, filename string) (*html.Tokenizer, error) {
	data, err := fs.ReadFile(fsys, filename)
	if err != nil {
		return nil, err
	}
	return html.NewTokenizer(bytes.NewReader(data)), nil
}
