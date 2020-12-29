package wit

import "golang.org/x/net/html"

// SetStyles sets provided attributes to provided values
type SetStyles struct {
	Styles map[string]string
}

// Apply applies the delta to the provided elements
func (s SetStyles) Apply(root *html.Node, nodes []*html.Node) {
	styles := s.Styles

	for _, node := range nodes {
		if node.Type != html.ElementNode {
			continue
		}

		found := false
		for i, att := range node.Attr {
			if att.Namespace != "" {
				continue
			}

			if att.Key == "style" {
				parsed := parseStyle(att.Val)
				for key, value := range styles {
					parsed[key] = value
				}

				att.Val = buildStyle(parsed)
				node.Attr[i] = att
				found = true
				break
			}
		}

		if !found {
			node.Attr = append(node.Attr, html.Attribute{
				Key: "style",
				Val: buildStyle(styles),
			})
		}
	}

}

// MarshalJSON marshals the delta to JSON format
func (s SetStyles) MarshalJSON() ([]byte, error) {
	return []byte("[" + setStylesLabelJSON + "," + strMapToJSON(s.Styles) + "]"), nil
}
