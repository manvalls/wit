package wit

import "golang.org/x/net/html"

// RmStyles removes provided CSS properties from matching elements
type RmStyles struct {
	Styles []string
}

// Apply applies the delta to the provided elements
func (r RmStyles) Apply(root *html.Node, nodes []*html.Node) {
	styles := r.Styles

	for _, node := range nodes {
		if node.Type != html.ElementNode {
			continue
		}

		for i, att := range node.Attr {
			if att.Namespace != "" {
				continue
			}

			if att.Key == "style" {
				parsed := parseStyle(att.Val)
				for _, s := range styles {
					delete(parsed, s)
				}

				att.Val = buildStyle(parsed)
				node.Attr[i] = att
				break
			}
		}
	}
}

// MarshalJSON marshals the delta to JSON format
func (r RmStyles) MarshalJSON() ([]byte, error) {
	return []byte("[" + rmStylesLabelJSON + strSliceToQuotedCSV(r.Styles) + "]"), nil
}
