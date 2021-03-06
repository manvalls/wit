package wit

import (
	"strconv"

	"golang.org/x/net/html"
)

// Replace replaces matching elements with the provided content
type Replace struct {
	HTMLSource
}

// Apply applies the delta to the provided elements
func (r Replace) Apply(d Document) {
	for _, node := range d.nodes {
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
