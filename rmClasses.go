package wit

import (
	"strconv"

	"golang.org/x/net/html"
)

// RmClasses removes provided classes from matching elements
type RmClasses struct {
	Classes string
}

// Apply applies the delta to the provided elements
func (r RmClasses) Apply(d Document) {
	class := r.Classes
	classesToRm := parseClass(class)

	for _, node := range d.nodes {
		if node.Type != html.ElementNode {
			continue
		}

		for i, att := range node.Attr {
			if att.Namespace != "" {
				continue
			}

			if att.Key == "class" {
				parsed := parseClass(att.Val)

				for key, value := range classesToRm {
					if value {
						delete(parsed, key)
					} else {
						parsed[key] = true
					}
				}

				att.Val = buildClass(parsed)
				node.Attr[i] = att
				break
			}
		}
	}
}

// MarshalJSON marshals the delta to JSON format
func (r RmClasses) MarshalJSON() ([]byte, error) {
	return []byte("[" + rmClassesLabelJSON + "," + strconv.Quote(r.Classes) + "]"), nil
}
