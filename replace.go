package wit

import (
	"strconv"

	"golang.org/x/net/html"
)

// Replace replaces matching elements with the provided content
type Replace struct {
	HTMLSource
}

// Empty returns whether or not this delta is empty
func (r Replace) Empty() bool {
	return false
}

// Flatten returns a new delta with redundant information removed
func (r Replace) Flatten() Delta {
	return r
}

// Apply applies the delta to the provided elements
func (r Replace) Apply(root *html.Node, nodes []*html.Node) {
	for _, node := range nodes {
		if node.Type != html.ElementNode || node.Parent == nil {
			continue
		}

		children := r.HTMLSource.Nodes(node.Parent)
		for _, child := range children {
			node.Parent.InsertBefore(child, node)
		}

		node.Parent.RemoveChild(node)
	}
}

// MarshalJSON marshals the delta to JSON format
func (r Replace) MarshalJSON() ([]byte, error) {
	return []byte("[" + replaceLabelJSON + "," + strconv.Quote(r.HTMLSource.String()) + "]"), nil
}
