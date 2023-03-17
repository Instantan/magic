package magic

import (
	"bytes"
	"errors"
	"io"

	"golang.org/x/net/html"
)

func (t *Template) parseHTML() ([]html.Token, error) {
	tokenizer := html.NewTokenizer(bytes.NewReader(t.data))
	tokens := []html.Token{}
	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			if errors.Is(tokenizer.Err(), io.EOF) {
				return tokens, tokenizer.Err()
			}
			return nil, tokenizer.Err()
		case html.StartTagToken:
			token := tokenizer.Token()
			if name, _ := tokenizer.TagName(); string(name) == "html" {
				token.Attr = append(token.Attr, html.Attribute{
					Key: "data-ssr",
					Val: "",
				})
			}
			tokens = append(tokens, token)
		default:
			tokens = append(tokens, tokenizer.Token())
		}
	}
}
