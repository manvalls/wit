package wit

import "golang.org/x/net/html"

// FirstChild applies given delta to the first child element
type FirstChild struct {
	Delta
}

// Apply applies the delta to the provided elements
func (fc FirstChild) Apply(root *html.Node, nodes []*html.Node) {
	childNodes := make([]*html.Node, 0, len(nodes))

	for _, node := range nodes {
		m := node.FirstChild
		for m != nil && m.Type != html.ElementNode {
			m = m.NextSibling
		}

		if m != nil {
			childNodes = append(childNodes, m)
		}
	}

	fc.Delta.Apply(root, nodes)
}

// MarshalJSON marshals the delta to JSON format
func (fc FirstChild) MarshalJSON() ([]byte, error) {
	return []byte("[" + firstChildLabelJSON + deltaToCSV(fc.Delta) + "]"), nil
}
