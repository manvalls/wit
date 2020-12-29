package wit

import "golang.org/x/net/html"

// Parent applies given delta to the parent element
type Parent struct {
	Delta
}

// Apply applies the delta to the provided elements
func (p Parent) Apply(root *html.Node, nodes []*html.Node) {
	childNodes := make([]*html.Node, 0, len(nodes))

	for _, node := range nodes {
		m := node.Parent
		if m != nil {
			childNodes = append(childNodes, m)
		}
	}

	p.Delta.Apply(root, childNodes)
}

// MarshalJSON marshals the delta to JSON format
func (p Parent) MarshalJSON() ([]byte, error) {
	return []byte("[" + parentLabelJSON + deltaToCSV(p.Delta) + "]"), nil
}
