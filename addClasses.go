package wit

import (
	"strconv"

	"golang.org/x/net/html"
)

// AddClasses adds provided classes to matching elements
type AddClasses struct {
	Classes string
}

// Apply applies the delta to the provided elements
func (a AddClasses) Apply(d Document) {
	class := a.Classes
	classesToAdd := parseClass(class)

	for _, node := range d.nodes {
		if node.Type != html.ElementNode {
			continue
		}

		found := false
		for i, att := range node.Attr {
			if att.Namespace != "" {
				continue
			}

			if att.Key == "class" {
				parsed := parseClass(att.Val)
				for key, value := range classesToAdd {
					parsed[key] = value
				}

				att.Val = buildClass(parsed)
				node.Attr[i] = att
				found = true
				break
			}
		}

		if !found {
			node.Attr = append(node.Attr, html.Attribute{
				Key: "class",
				Val: buildClass(classesToAdd),
			})
		}
	}
}

// MarshalJSON marshals the delta to JSON format
func (a AddClasses) MarshalJSON() ([]byte, error) {
	return []byte("[" + addClassesLabelJSON + "," + strconv.Quote(a.Classes) + "]"), nil
}
