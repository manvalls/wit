package wit

import (
	"strconv"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

// All applies given delta to all matching elements
type All struct {
	Selector
	Delta
}

// Apply applies the delta to the provided elements
func (a All) Apply(d Document) {
	childNodes := make([]*html.Node, 0, len(d.nodes))

	for _, node := range d.nodes {
		for _, match := range cascadia.QueryAll(node, a.Selector) {
			childNodes = append(childNodes, match)
		}
	}

	d.nodes = childNodes
	a.Delta.Apply(d)
}

// MarshalJSON marshals the delta to JSON format
func (a All) MarshalJSON() ([]byte, error) {
	return []byte("[" + selectorAllLabelJSON + "," + strconv.Quote(a.Selector.String()) + deltaToCSV(a.Delta) + "]"), nil
}
