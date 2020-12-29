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
func (f First) Apply(root *html.Node, nodes []*html.Node) {
	childNodes := make([]*html.Node, 0, len(nodes))

	for _, node := range nodes {
		match := cascadia.Query(node, f.Selector)
		if match != nil {
			childNodes = append(childNodes, match)
		}
	}

	f.Delta.Apply(root, childNodes)
}

// MarshalJSON marshals the delta to JSON format
func (f First) MarshalJSON() ([]byte, error) {
	return []byte("[" + selectorLabelJSON + "," + strconv.Quote(f.Selector.String()) + deltaToCSV(f.Delta) + "]"), nil
}
