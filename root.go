package wit

import "golang.org/x/net/html"

// Root applies given delta to the root of the document
type Root struct {
	Delta
}

// Flatten returns a new delta with redundant information removed
func (r Root) Flatten() Delta {
	return Root{r.Delta.Flatten()}
}

// Apply applies the delta to the provided elements
func (r Root) Apply(root *html.Node, nodes []*html.Node) {
	r.Delta.Apply(root, []*html.Node{root})
}

// MarshalJSON marshals the delta to JSON format
func (r Root) MarshalJSON() ([]byte, error) {
	return []byte("[" + rootLabelJSON + deltaToCSV(r.Delta) + "]"), nil
}
