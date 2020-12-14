package wit

import "golang.org/x/net/html"

// Clear empties matching elements
type Clear struct{}

// Empty returns whether or not this delta is empty
func (c Clear) Empty() bool {
	return false
}

// Flatten returns a new delta with redundant information removed
func (c Clear) Flatten() Delta {
	return c
}

// Apply applies the delta to the provided elements
func (c Clear) Apply(root *html.Node, nodes []*html.Node) {
	for _, node := range nodes {
		if node.Type == html.ElementNode {
			for node.FirstChild != nil {
				node.RemoveChild(node.FirstChild)
			}
		}
	}
}

// MarshalJSON marshals the delta to JSON format
func (c Clear) MarshalJSON() ([]byte, error) {
	return []byte("[" + clearLabelJSON + "]"), nil
}
