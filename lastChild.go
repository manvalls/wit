package wit

import "golang.org/x/net/html"

// LastChild applies given delta to the last child element
type LastChild struct {
	Delta
}

// Apply applies the delta to the provided elements
func (lc LastChild) Apply(d Document) {
	childNodes := make([]*html.Node, 0, len(d.nodes))

	for _, node := range d.nodes {
		m := node.LastChild
		for m != nil && m.Type != html.ElementNode {
			m = m.PrevSibling
		}

		if m != nil {
			childNodes = append(childNodes, m)
		}
	}

	d.nodes = childNodes
	lc.Delta.Apply(d)
}

// MarshalJSON marshals the delta to JSON format
func (lc LastChild) MarshalJSON() ([]byte, error) {
	return []byte("[" + lastChildLabelJSON + deltaToCSV(lc.Delta) + "]"), nil
}
