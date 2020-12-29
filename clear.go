package wit

import "golang.org/x/net/html"

// Clear empties matching elements
type Clear struct{}

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
