package wit

import "golang.org/x/net/html"

// NextSibling applies given delta to the next sibling
type NextSibling struct {
	Delta
}

// Apply applies the delta to the provided elements
func (ns NextSibling) Apply(d Document) {
	childNodes := make([]*html.Node, 0, len(d.nodes))

	for _, node := range d.nodes {
		m := node.NextSibling
		for m != nil && m.Type != html.ElementNode {
			m = m.NextSibling
		}

		if m != nil {
			childNodes = append(childNodes, m)
		}
	}

	d.nodes = childNodes
	ns.Delta.Apply(d)
}

// MarshalJSON marshals the delta to JSON format
func (ns NextSibling) MarshalJSON() ([]byte, error) {
	return []byte("[" + nextSiblingLabelJSON + deltaToCSV(ns.Delta) + "]"), nil
}
