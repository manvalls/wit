package wit

import "golang.org/x/net/html"

// NextSibling applies given delta to the next sibling
type NextSibling struct {
	Delta
}

// Apply applies the delta to the provided elements
func (ns NextSibling) Apply(root *html.Node, nodes []*html.Node) {
	childNodes := make([]*html.Node, 0, len(nodes))

	for _, node := range nodes {
		m := node.NextSibling
		for m != nil && m.Type != html.ElementNode {
			m = m.NextSibling
		}

		if m != nil {
			childNodes = append(childNodes, m)
		}
	}

	ns.Delta.Apply(root, childNodes)
}

// MarshalJSON marshals the delta to JSON format
func (ns NextSibling) MarshalJSON() ([]byte, error) {
	return []byte("[" + nextSiblingLabelJSON + deltaToCSV(ns.Delta) + "]"), nil
}
