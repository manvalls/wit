package wit

import (
	"strconv"

	"golang.org/x/net/html"
)

// HTML sets the inner HTML of matching elements
type HTML struct {
	HTMLSource
}

// Empty returns whether or not this delta is empty
func (h HTML) Empty() bool {
	return false
}

// Flatten returns a new delta with redundant information removed
func (h HTML) Flatten() Delta {
	return h
}

// Apply applies the delta to the provided elements
func (h HTML) Apply(root *html.Node, nodes []*html.Node) {
	for _, node := range nodes {
		for node.FirstChild != nil {
			node.RemoveChild(node.FirstChild)
		}

		children := h.HTMLSource.Nodes(node)
		for _, child := range children {
			node.AppendChild(child)
		}
	}
}

// MarshalJSON marshals the delta to JSON format
func (h HTML) MarshalJSON() ([]byte, error) {
	return []byte("[" + htmlLabelJSON + "," + strconv.Quote(h.HTMLSource.String()) + "]"), nil
}
