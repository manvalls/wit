package wit

import "golang.org/x/net/html"

// Parent applies given delta to the parent element
type Parent struct {
	Delta
}

// Apply applies the delta to the provided elements
func (p Parent) Apply(d Document) {
	childNodes := make([]*html.Node, 0, len(d.nodes))

	for _, node := range d.nodes {
		m := node.Parent
		if m != nil {
			childNodes = append(childNodes, m)
		}
	}

	d.nodes = childNodes
	p.Delta.Apply(d)
}

// MarshalJSON marshals the delta to JSON format
func (p Parent) MarshalJSON() ([]byte, error) {
	return []byte("[" + parentLabelJSON + deltaToCSV(p.Delta) + "]"), nil
}
