package wit

import (
	"strconv"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

// First applies given delta to the first matching element
type First struct {
	Selector
	Delta
}

// Apply applies the delta to the provided elements
func (f First) Apply(d Document) {
	childNodes := make([]*html.Node, 0, len(d.nodes))

	for _, node := range d.nodes {
		match := cascadia.Query(node, f.Selector)
		if match != nil {
			childNodes = append(childNodes, match)
		}
	}

	d.nodes = childNodes
	f.Delta.Apply(d)
}

// MarshalJSON marshals the delta to JSON format
func (f First) MarshalJSON() ([]byte, error) {
	return []byte("[" + selectorLabelJSON + "," + strconv.Quote(f.Selector.String()) + deltaToCSV(f.Delta) + "]"), nil
}
