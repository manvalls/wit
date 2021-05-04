package wit

import "golang.org/x/net/html"

// PrevSibling applies given delta to the previous sibling
type PrevSibling struct {
	Delta
}

// Apply applies the delta to the provided elements
func (ps PrevSibling) Apply(d Document) {
	childNodes := make([]*html.Node, 0, len(d.nodes))

	for _, node := range d.nodes {
		m := node.PrevSibling
		for m != nil && m.Type != html.ElementNode {
			m = m.PrevSibling
		}

		if m != nil {
			childNodes = append(childNodes, m)
		}
	}

	d.nodes = childNodes
	ps.Delta.Apply(d)
}

// MarshalJSON marshals the delta to JSON format
func (ps PrevSibling) MarshalJSON() ([]byte, error) {
	return []byte("[" + prevSiblingLabelJSON + deltaToCSV(ps.Delta) + "]"), nil
}
