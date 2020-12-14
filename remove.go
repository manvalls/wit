package wit

import "golang.org/x/net/html"

// Remove removes matching elements
type Remove struct{}

// Empty returns whether or not this delta is empty
func (r Remove) Empty() bool {
	return false
}

// Flatten returns a new delta with redundant information removed
func (r Remove) Flatten() Delta {
	return r
}

// Apply applies the delta to the provided elements
func (r Remove) Apply(root *html.Node, nodes []*html.Node) {
	for _, node := range nodes {
		parent := node.Parent
		if parent != nil {
			parent.RemoveChild(node)
		}
	}
}

// MarshalJSON marshals the delta to JSON format
func (r Remove) MarshalJSON() ([]byte, error) {
	return []byte("[" + removeLabelJSON + "]"), nil
}
