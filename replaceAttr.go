package wit

import "golang.org/x/net/html"

// ReplaceAttr replaces the attributes of matching elements
type ReplaceAttr struct {
	Attributes map[string]string
}

// Apply applies the delta to the provided elements
func (r ReplaceAttr) Apply(root *html.Node, nodes []*html.Node) {
	attr := r.Attributes
	for _, node := range nodes {
		if node.Type != html.ElementNode {
			continue
		}

		i := 0
		nodeAttr := make([]html.Attribute, len(attr))
		for key, value := range attr {
			nodeAttr[i] = html.Attribute{
				Key: key,
				Val: value,
			}

			i++
		}

		node.Attr = nodeAttr
	}
}

// MarshalJSON marshals the delta to JSON format
func (r ReplaceAttr) MarshalJSON() ([]byte, error) {
	return []byte("[" + replaceAttrLabelJSON + "," + strMapToJSON(r.Attributes) + "]"), nil
}
