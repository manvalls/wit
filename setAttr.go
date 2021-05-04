package wit

import (
	"golang.org/x/net/html"
)

// SetAttr sets provided attributes to provided values
type SetAttr struct {
	Attributes map[string]string
}

// Apply applies the delta to the provided elements
func (s SetAttr) Apply(d Document) {
	attr := s.Attributes
	for _, node := range d.nodes {
		if node.Type != html.ElementNode {
			continue
		}

		nodeAttr := map[string]int{}
		for i, att := range node.Attr {
			if att.Namespace != "" {
				continue
			}

			if _, ok := nodeAttr[att.Key]; ok {
				continue
			}

			nodeAttr[att.Key] = i
		}

		for key, value := range attr {
			i, ok := nodeAttr[key]
			if ok {
				node.Attr[i] = html.Attribute{
					Key: key,
					Val: value,
				}
			} else {
				node.Attr = append(node.Attr, html.Attribute{
					Key: key,
					Val: value,
				})
			}
		}
	}
}

// MarshalJSON marshals the delta to JSON format
func (s SetAttr) MarshalJSON() ([]byte, error) {
	return []byte("[" + setAttrLabelJSON + "," + strMapToJSON(s.Attributes) + "]"), nil
}
