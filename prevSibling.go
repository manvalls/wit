package wit

import "golang.org/x/net/html"

// PrevSibling applies given delta to the previous sibling
type PrevSibling struct {
	Delta
}

// Apply applies the delta to the provided elements
func (ps PrevSibling) Apply(root *html.Node, nodes []*html.Node) {
	childNodes := make([]*html.Node, 0, len(nodes))

	for _, node := range nodes {
		m := node.PrevSibling
		for m != nil && m.Type != html.ElementNode {
			m = m.PrevSibling
		}

		if m != nil {
			childNodes = append(childNodes, m)
		}
	}

	ps.Delta.Apply(root, childNodes)
}

// MarshalJSON marshals the delta to JSON format
func (ps PrevSibling) MarshalJSON() ([]byte, error) {
	return []byte("[" + prevSiblingLabelJSON + deltaToCSV(ps.Delta) + "]"), nil
}
