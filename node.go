package magic

type Node struct {
	Tag        string
	Attributes []Attribute
	Children   []Node
}

type Attribute struct {
	Name  string
	Value string
}

func createHtmlTag(tag string, attributes []Attribute, children []Node) Node {
	return Node{
		Tag:        tag,
		Attributes: attributes,
		Children:   children,
	}
}

func createAttribute(name string, value string) Attribute {
	return Attribute{
		Name:  name,
		Value: value,
	}
}

func (node Node) HTML() string {
	attributes := ""
	for _, attr := range node.Attributes {
		attributes += attr.Name + "=" + "\"" + attr.Value + "\" "
	}
	children := ""
	for _, child := range node.Children {
		children += child.HTML()
	}
	return "<" + node.Tag + " " + attributes + ">" + children + "</" + node.Tag + ">"
}
