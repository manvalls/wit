package wit

import "golang.org/x/net/html"

// RmAttr removes provided attributes
type RmAttr struct {
	Attributes []string
}

// Apply applies the delta to the provided elements
func (r RmAttr) Apply(root *html.Node, nodes []*html.Node) {
	attr := r.Attributes
	attrMap := map[string]bool{}
	for _, key := range attr {
		attrMap[key] = true
	}

	for _, node := range nodes {
		if node.Type != html.ElementNode {
			continue
		}

		nodeAttr := make([]html.Attribute, 0, len(node.Attr))
		for _, att := range node.Attr {
			if att.Namespace != "" {
				continue
			}

			if !attrMap[att.Key] {
				nodeAttr = append(nodeAttr, att)
			}
		}

		node.Attr = nodeAttr
	}
}

// MarshalJSON marshals the delta to JSON format
func (r RmAttr) MarshalJSON() ([]byte, error) {
	return []byte("[" + rmAttrLabelJSON + strSliceToQuotedCSV(r.Attributes) + "]"), nil
}
