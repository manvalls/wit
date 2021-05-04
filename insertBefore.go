package wit

import (
	"strconv"

	"golang.org/x/net/html"
)

// InsertBefore inserts the provided HTML before matching elements
type InsertBefore struct {
	HTMLSource
}

// Apply applies the delta to the provided elements
func (i InsertBefore) Apply(d Document) {
	for _, node := range d.nodes {
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
