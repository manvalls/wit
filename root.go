package wit

import "golang.org/x/net/html"

// Root applies given delta to the root of the document
type Root struct {
	Delta
}

// Apply applies the delta to the provided elements
func (r Root) Apply(d Document) {
	d.nodes = []*html.Node{d.root}
	r.Delta.Apply(d)
}

// MarshalJSON marshals the delta to JSON format
func (r Root) MarshalJSON() ([]byte, error) {
	return []byte("[" + rootLabelJSON + deltaToCSV(r.Delta) + "]"), nil
}
