package wit

import (
	"strconv"
)

// HTML sets the inner HTML of matching elements
type HTML struct {
	HTMLSource
}

// Apply applies the delta to the provided elements
func (h HTML) Apply(d Document) {
	for _, node := range d.nodes {
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
