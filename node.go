package magic

type RxValue any

type Node struct {
	Tag        string
	Attributes []Attribute
	Children   []Node
	Data       RxValue
}

type Attribute struct {
	Name  string
	Value RxValue
}

func (node Node) HTML() string {
	if node.Tag == "" {
		return renderRxValueToString(node.Data)
	}
	attributes := ""
	for _, attr := range node.Attributes {
		attributes += attr.Name + "=" + "\"" + renderRxValueToString(attr.Value) + "\" "
	}
	children := ""
	for _, child := range node.Children {
		children += child.HTML()
	}
	if attributes == "" {
		return "<" + node.Tag + ">" + children + "</" + node.Tag + ">"
	}
	return "<" + node.Tag + " " + attributes + ">" + children + "</" + node.Tag + ">"
}

func renderRxValueToString(v RxValue) string {
	switch tv := v.(type) {
	case string:
		return tv
	case func() string:
		return tv()
	default:
		return ""
	}
}
