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

// Flatten returns a new delta with redundant information removed
func (a All) Flatten() Delta {
	return All{a.Selector, a.Delta.Flatten()}
}

// Apply applies the delta to the provided elements
func (a All) Apply(root *html.Node, nodes []*html.Node) {
	childNodes := make([]*html.Node, 0, len(nodes))

	for _, node := range nodes {
		for _, match := range cascadia.QueryAll(node, a.Selector) {
			childNodes = append(childNodes, match)
		}
	}

	a.Delta.Apply(root, childNodes)
}

// MarshalJSON marshals the delta to JSON format
func (a All) MarshalJSON() ([]byte, error) {
	return []byte("[" + selectorAllLabelJSON + "," + strconv.Quote(a.Selector.String()) + deltaToCSV(a.Delta) + "]"), nil
}
