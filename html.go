package wit

import (
	"io"
	"strings"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

var headSelector = cascadia.MustCompile("head")

var baseDocument, _ = html.Parse(strings.NewReader("<!DOCTYPE html><html><head></head><body></body></html>"))

type htmlContext struct {
	root *html.Node
}

type htmlRenderer struct {
	root *html.Node
}

// NewHTMLRenderer returns a new renderer which will render HTML
func NewHTMLRenderer(command Command) Renderer {
	nodes := clone([]*html.Node{baseDocument})
	c := &htmlContext{
		root: nodes[0],
	}

	if !IsNil(command) {
		applyDelta(c, nodes, command.Delta())
	}

	return &htmlRenderer{c.root}
}

func (r *htmlRenderer) Render(w io.Writer) error {
	return html.Render(w, r.root)
}

func applyDelta(c *htmlContext, nodes []*html.Node, delta Delta) {

	switch delta.typeID {

	case sliceType:
		for _, childDelta := range delta.info.deltas {
			applyDelta(c, nodes, childDelta)
		}

	case rootType:
		childNodes := []*html.Node{c.root}
		for _, childDelta := range delta.info.deltas {
			applyDelta(c, childNodes, childDelta)
		}

	case selectorType:
		selector := delta.info.selector.selector()
		childNodes := make([]*html.Node, 0, len(nodes))

		if selector != nil {
			for _, node := range nodes {
				child := node.FirstChild
				for child != nil {
					m := selector.MatchFirst(child)
					if m != nil {
						childNodes = append(childNodes, m)
						break
					}

					child = child.NextSibling
				}
			}
		}

		for _, childDelta := range delta.info.deltas {
			applyDelta(c, childNodes, childDelta)
		}

	case selectorAllType:
		selector := delta.info.selector.selector()
		childNodes := make([]*html.Node, 0, len(nodes))

		if selector != nil {
			for _, node := range nodes {
				child := node.FirstChild
				for child != nil {
					ms := selector.MatchAll(node)
					for _, m := range ms {
						childNodes = append(childNodes, m)
					}

					child = child.NextSibling
				}
			}
		}

		for _, childDelta := range delta.info.deltas {
			applyDelta(c, childNodes, childDelta)
		}

	case parentType:
		childNodes := make([]*html.Node, 0, len(nodes))

		for _, node := range nodes {
			m := node.Parent
			if m != nil {
				childNodes = append(childNodes, m)
			}
		}

		for _, childDelta := range delta.info.deltas {
			applyDelta(c, childNodes, childDelta)
		}

	case firstChildType:
		childNodes := make([]*html.Node, 0, len(nodes))

		for _, node := range nodes {
			m := node.FirstChild
			for m != nil && m.Type != html.ElementNode {
				m = m.NextSibling
			}

			if m != nil {
				childNodes = append(childNodes, m)
			}
		}

		for _, childDelta := range delta.info.deltas {
			applyDelta(c, childNodes, childDelta)
		}

	case lastChildType:
		childNodes := make([]*html.Node, 0, len(nodes))

		for _, node := range nodes {
			m := node.LastChild
			for m != nil && m.Type != html.ElementNode {
				m = m.PrevSibling
			}

			if m != nil {
				childNodes = append(childNodes, m)
			}
		}

		for _, childDelta := range delta.info.deltas {
			applyDelta(c, childNodes, childDelta)
		}

	case prevSiblingType:
		childNodes := make([]*html.Node, 0, len(nodes))

		for _, node := range nodes {
			m := node.PrevSibling
			for m != nil && m.Type != html.ElementNode {
				m = m.PrevSibling
			}

			if m != nil {
				childNodes = append(childNodes, m)
			}
		}

		for _, childDelta := range delta.info.deltas {
			applyDelta(c, childNodes, childDelta)
		}

	case nextSiblingType:
		childNodes := make([]*html.Node, 0, len(nodes))

		for _, node := range nodes {
			m := node.NextSibling
			for m != nil && m.Type != html.ElementNode {
				m = m.NextSibling
			}

			if m != nil {
				childNodes = append(childNodes, m)
			}
		}

		for _, childDelta := range delta.info.deltas {
			applyDelta(c, childNodes, childDelta)
		}

	case removeType:

		for _, node := range nodes {
			parent := node.Parent
			if parent != nil {
				parent.RemoveChild(node)
			}
		}

	case clearType:

		for _, node := range nodes {
			if node.Type == html.ElementNode {
				for node.FirstChild != nil {
					node.RemoveChild(node.FirstChild)
				}
			}
		}

	case htmlType:

		for _, node := range nodes {
			if node.Type != html.ElementNode {
				continue
			}

			for node.FirstChild != nil {
				node.RemoveChild(node.FirstChild)
			}

			children := delta.info.factory.Nodes(node)
			for _, child := range children {
				node.AppendChild(child)
			}
		}

	case replaceType:

		for _, node := range nodes {
			if node.Type != html.ElementNode || node.Parent == nil {
				continue
			}

			children := delta.info.factory.Nodes(node.Parent)
			for _, child := range children {
				node.Parent.InsertBefore(child, node)
			}

			node.Parent.RemoveChild(node)
		}

	case appendType:

		for _, node := range nodes {
			if node.Type != html.ElementNode {
				continue
			}

			children := delta.info.factory.Nodes(node)
			for _, child := range children {
				node.AppendChild(child)
			}
		}

	case prependType:

		for _, node := range nodes {
			if node.Type != html.ElementNode {
				continue
			}

			children := delta.info.factory.Nodes(node)

			if node.FirstChild != nil {
				for _, child := range children {
					node.InsertBefore(child, node.FirstChild)
				}
			} else {
				for _, child := range children {
					node.AppendChild(child)
				}
			}

		}

	case insertAfterType:

		for _, node := range nodes {
			if node.Type != html.ElementNode || node.Parent == nil {
				continue
			}

			children := delta.info.factory.Nodes(node)

			if node.NextSibling != nil {
				for _, child := range children {
					node.Parent.InsertBefore(child, node.NextSibling)
				}
			} else {
				for _, child := range children {
					node.Parent.AppendChild(child)
				}
			}
		}

	case insertBeforeType:

		for _, node := range nodes {
			if node.Type != html.ElementNode || node.Parent == nil {
				continue
			}

			children := delta.info.factory.Nodes(node)
			for _, child := range children {
				node.Parent.InsertBefore(child, node)
			}
		}

	case addAttrType:
		attr := delta.info.strMap
		for _, node := range nodes {
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

	case setAttrType:
		attr := delta.info.strMap
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

	case rmAttrType:
		attr := delta.info.strList
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

	case addStylesType:
		styles := delta.info.strMap

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

	case rmStylesType:
		styles := delta.info.strList

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

	case addClassType:
		class := delta.info.class
		classesToAdd := parseClass(class)

		for _, node := range nodes {
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

	case rmClassType:
		class := delta.info.class
		classesToRm := parseClass(class)

		for _, node := range nodes {
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

	return
}
