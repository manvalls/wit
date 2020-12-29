package wit

import (
	"strconv"

	"golang.org/x/net/html"
)

// InsertAfter inserts the provided HTML after matching elements
type InsertAfter struct {
	HTMLSource
}

// Apply applies the delta to the provided elements
func (i InsertAfter) Apply(root *html.Node, nodes []*html.Node) {
	for _, node := range nodes {
		if node.Type != html.ElementNode || node.Parent == nil {
			continue
		}

		children := i.HTMLSource.Nodes(node)

		if node.NextSibling != nil {
			for _, child := range children {
				node.Parent.InsertBefore(child, node.NextSibling)
			}
		} else {
			for _, child := range children {
				node.Parent.AppendChild(child)
			}
		}
	}
}

// MarshalJSON marshals the delta to JSON format
func (i InsertAfter) MarshalJSON() ([]byte, error) {
	return []byte("[" + insertAfterLabelJSON + "," + strconv.Quote(i.HTMLSource.String()) + "]"), nil
}
