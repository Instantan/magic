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
	node := parseView(prepareDynamicValues(minifyView(template)))
	makeViewDynamic(&node)
	return func(props Props) Node {
		return node
	}
}

func minifyView(template string) string {
	template, err := minifier.String("text/html", template)
	if err != nil {
		panic(err)
	}
	return template
}

func prepareDynamicValues(template string) string {
	re := regexp.MustCompile(`(?m)<.*=.*{{.*}}.*.*>`)
	template = re.ReplaceAllStringFunc(template, func(s string) string {
		print(s)
		s = strings.ReplaceAll(s, "{{", "{%")
		s = strings.ReplaceAll(s, "}}", "%}")
		return s
	})
	re = regexp.MustCompile(`(?m){{\s*end\s*}}`)
	template = re.ReplaceAllString(template, "</m>")
	re = regexp.MustCompile(`((?m){{\s*(range|if)\s.*}})`)
	template = re.ReplaceAllStringFunc(template, func(s string) string {
		s = strings.ReplaceAll(s, "{{", "<m expr=\"")
		s = strings.ReplaceAll(s, "}}", "\">")
		return s
	})
	re = regexp.MustCompile(`((?m){{\s*)`)
	template = re.ReplaceAllString(template, "<m expr=\"")
	re = regexp.MustCompile(`((?m)\s*}})`)
	template = re.ReplaceAllString(template, "\"/>")
	return template
}

func parseView(template string) Node {
	tkn := html.NewTokenizer(strings.NewReader(template))
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
	for {
		tt := tkn.Next()
		if tt == html.ErrorToken || tt == html.EndTagToken {
			break
		}
		if tt != html.StartTagToken && tt != html.SelfClosingTagToken && tt != html.TextToken {
			continue
		}
		nodes = append(nodes, parseNode(tkn))
	}
	return nodes
}

func parseNode(tkn *html.Tokenizer) Node {
	n := Node{}
	token := tkn.Token()
	if token.Type == html.StartTagToken || token.Type == html.SelfClosingTagToken {
		n.Tag = token.Data
		for i := range token.Attr {
			n.Attributes = append(n.Attributes, Attribute{
				Name:  token.Attr[i].Key,
				Value: parseTemplate(token.Attr[i].Val),
			})
		}
		if token.Type == html.SelfClosingTagToken {
			return n
		}
	} else if token.Type == html.ErrorToken || token.Type == html.EndTagToken {
		return n
	} else if token.Type == html.TextToken {
		n.Data = token.Data
		return n
	} else {
		tkn.Next()
		return parseNode(tkn)
	}
	n.Children = parseNodes(tkn)
	return n
}

func makeViewDynamic(node *Node) {
	for i := range node.Attributes {
		makeAttributeDynamic(&node.Attributes[i])
	}
	expr := strings.TrimSpace(getNodeExpr(*node))
	if expr == "" {
		for i := range node.Children {
			makeViewDynamic(&node.Children[i])
		}
		return
	}
	exprParts := removeEmptyExprParts(strings.Split(expr, " "))
	node.Attributes = []Attribute{}
	node.Data = buildFnFromExprPartsAndChildren(exprParts, node.Children)
}

func makeAttributeDynamic(attr *Attribute) {
	template := attr.Value.(string)
	if !strings.Contains(template, "{%") || !strings.Contains(template, "%}") {
		return
	}
	attr.Value = parseTemplate(template)
}

func getNodeExpr(node Node) string {
	if node.Tag == "m" && node.Data != nil {
		expr := node.Data.(string)
		for i := range node.Attributes {
			if node.Attributes[i].Name == "expr" {
				expr = node.Attributes[i].Value.(string)
				break
			}
		}
		return expr
	}
	return ""
}

func parseTemplate(data string) RxValue {
	return data
}

func removeEmptyExprParts(exprParts []string) []string {
	newParts := []string{}
	for i := range exprParts {
		if strings.TrimSpace(exprParts[i]) == "" {
			continue
		}
		newParts = append(newParts, exprParts[i])
	}
	return newParts
}

func buildFnFromExprPartsAndChildren(parts []string, children []Node) RxValue {
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 2 {
		switch parts[0] {
		case "range":
			return buildRangeFnFromChildren(parts[1], children)
		case "if":
			return buildIfFnFromChildren(parts[1], children)
		}
	} else if len(parts) == 1 {
		return buildValueFn(parts[0])
	}
	return ""
}

func buildRangeFnFromChildren(selector string, children []Node) RxValue {
	return func() []Node {
		return []Node{}
	}
}

func buildIfFnFromChildren(selector string, children []Node) RxValue {
	return func() []Node {
		return []Node{}
	}
}

func buildValueFn(selector string) RxValue {
	return func(data any) string {
		return "BLA"
	}
}
