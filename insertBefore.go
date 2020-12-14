package wit

import (
	"strconv"

	"golang.org/x/net/html"
)

// InsertBefore inserts the provided HTML before matching elements
type InsertBefore struct {
	HTMLSource
}

// Empty returns whether or not this delta is empty
func (i InsertBefore) Empty() bool {
	return false
}

// Flatten returns a new delta with redundant information removed
func (i InsertBefore) Flatten() Delta {
	return i
}

// Apply applies the delta to the provided elements
func (i InsertBefore) Apply(root *html.Node, nodes []*html.Node) {
	for _, node := range nodes {
		if node.Type != html.ElementNode || node.Parent == nil {
			continue
		}

		children := i.HTMLSource.Nodes(node)
		for _, child := range children {
			node.Parent.InsertBefore(child, node)
		}
	}
}

// MarshalJSON marshals the delta to JSON format
func (i InsertBefore) MarshalJSON() ([]byte, error) {
	return []byte("[" + insertBeforeLabelJSON + "," + strconv.Quote(i.HTMLSource.String()) + "]"), nil
}
