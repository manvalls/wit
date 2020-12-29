package wit

import (
	"strconv"

	"golang.org/x/net/html"
)

// Append appends the provided HTML to matching elements
type Append struct {
	HTMLSource
}

// Apply applies the delta to the provided elements
func (a Append) Apply(root *html.Node, nodes []*html.Node) {
	for _, node := range nodes {
		if node.Type != html.ElementNode {
			continue
		}

		children := a.HTMLSource.Nodes(node)
		for _, child := range children {
			node.AppendChild(child)
		}
	}
}

// MarshalJSON marshals the delta to JSON format
func (a Append) MarshalJSON() ([]byte, error) {
	return []byte("[" + appendLabelJSON + "," + strconv.Quote(a.HTMLSource.String()) + "]"), nil
}
