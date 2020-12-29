package wit

import (
	"strconv"

	"golang.org/x/net/html"
)

// Prepend prepends the provided HTML to matching elements
type Prepend struct {
	HTMLSource
}

// Apply applies the delta to the provided elements
func (p Prepend) Apply(root *html.Node, nodes []*html.Node) {
	for _, node := range nodes {
		if node.Type != html.ElementNode {
			continue
		}

		children := p.HTMLSource.Nodes(node)

		if node.FirstChild != nil {
			for _, child := range children {
				node.InsertBefore(child, node.FirstChild)
			}
		} else {
			for _, child := range children {
				node.AppendChild(child)
			}
		}

	}
}

// MarshalJSON marshals the delta to JSON format
func (p Prepend) MarshalJSON() ([]byte, error) {
	return []byte("[" + prependLabelJSON + "," + strconv.Quote(p.HTMLSource.String()) + "]"), nil
}
