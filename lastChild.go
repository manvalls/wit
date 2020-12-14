package wit

import "golang.org/x/net/html"

// LastChild applies given delta to the last child element
type LastChild struct {
	Delta
}

// Flatten returns a new delta with redundant information removed
func (lc LastChild) Flatten() Delta {
	return LastChild{lc.Delta.Flatten()}
}

// Apply applies the delta to the provided elements
func (lc LastChild) Apply(root *html.Node, nodes []*html.Node) {
	childNodes := make([]*html.Node, 0, len(nodes))

	for _, node := range nodes {
		m := node.LastChild
		for m != nil && m.Type != html.ElementNode {
			m = m.PrevSibling
		}

		if m != nil {
			childNodes = append(childNodes, m)
		}
	}

	lc.Delta.Apply(root, childNodes)
}

// MarshalJSON marshals the delta to JSON format
func (lc LastChild) MarshalJSON() ([]byte, error) {
	return []byte("[" + lastChildLabelJSON + deltaToCSV(lc.Delta) + "]"), nil
}
