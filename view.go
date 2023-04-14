package magic

import (
	"regexp"
	"strings"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	mhtml "github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/json"
	"github.com/tdewolff/minify/svg"
	"github.com/tdewolff/minify/xml"
	"golang.org/x/net/html"
)

var minifier *minify.M

func init() {
	minifier = minify.New()
	minifier.AddFunc("text/css", css.Minify)
	minifier.Add("text/html", &mhtml.Minifier{
		KeepDefaultAttrVals: true,
		KeepWhitespace:      true,
	})
	minifier.AddFunc("image/svg+xml", svg.Minify)
	minifier.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	minifier.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	minifier.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
}

type ViewFn[Props any] func(props Props) Node

func View[Props any](template string) ViewFn[Props] {
	node := parseTemplate(minifyTemplate(template))
	return func(props Props) Node {
		return node
	}
}

func minifyTemplate(template string) string {
	template, err := minifier.String("text/html", template)
	if err != nil {
		panic(err)
	}
	return template
}

func parseTemplate(template string) Node {
	tkn := html.NewTokenizer(strings.NewReader(template))
	tt := tkn.Next()
	if tt == html.ErrorToken {
		return Node{}
	}
	nodes := parseNodes(tkn)
	if len(nodes) == 0 {
		return Node{}
	}
	if len(nodes) == 1 {
		return nodes[0]
	}
	return Node{
		Tag:      "m",
		Children: nodes,
	}
}

func parseNodes(tkn *html.Tokenizer) []Node {
	nodes := []Node{}
	node := Node{}
	for {
		tt := tkn.Next()
		switch {
		case tt == html.ErrorToken:
			break
		case tt == html.StartTagToken:
			if node.Tag != "" {
				nodes = append(nodes, node)
				node = Node{}
			}
			token := tkn.Token()
			node.Tag = token.Data
			for i := range token.Attr {
				node.Attributes = append(node.Attributes, Attribute{
					Name:  token.Attr[i].Key,
					Value: token.Attr[i].Val,
				})
			}
		case tt == html.EndTagToken:
			if node.Tag != "" {
				nodes = append(nodes, node)
				node = Node{}
			}
		case tt == html.TextToken:
			if node.Tag != "" {
				nodes = append(nodes, node)
				node = Node{}
			}
			t := tkn.Token()
		}
	}
	if node.Tag != "" {
		nodes = append(nodes, node)
		node = Node{}
	}
	return nodes
}
