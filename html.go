package wit

import (
	"net/http"
	"strings"

	"github.com/manvalls/wit/util"
	"golang.org/x/net/html"
)

var baseDocument, _ = html.Parse(strings.NewReader("<!DOCTYPE html><html><head></head><body></body></html>"))

// WriteHTML writes the result of applying the provided delta to an empty
// document as formatted HTML
func WriteHTML(w http.ResponseWriter, delta Delta) error {
	nodes := util.Clone([]*html.Node{baseDocument})
	root := nodes[0]
	newRoot, _ := applyDelta(root, nodes, delta)
	return html.Render(w, newRoot)
}

func applyDelta(root *html.Node, nodes []*html.Node, delta Delta) (newRoot *html.Node, discardNext bool) {

	newRoot = root
	discardNext = false

	switch delta.typeID {

	case sliceType:
		deltas := delta.delta.(*deltaSlice).deltas

		for _, childDelta := range deltas {
			if discardNext {
				discardDelta(childDelta)
			} else {
				newRoot, discardNext = applyDelta(root, nodes, childDelta)
			}
		}

	case channelType:
		d := delta.delta.(*deltaChannel)

		channel := d.channel
		cancel := d.cancel

		for childDelta := range channel {
			if discardNext {
				discardDelta(childDelta)
			} else {
				newRoot, discardNext = applyDelta(root, nodes, childDelta)
				if discardNext {
					cancel()
				}
			}
		}

	case rootType:
		return applyDelta(root, []*html.Node{root}, delta.delta.(*deltaRoot).delta)

	case selectorType:
		d := delta.delta.(*deltaSelector)
		selector := d.selector.selector()

		if selector != nil {
			childNodes := make([]*html.Node, 0, len(nodes))
			for _, node := range nodes {
				m := selector.MatchFirst(node)
				if m != nil {
					childNodes = append(childNodes, m)
				}
			}

			if len(childNodes) > 0 {
				return applyDelta(root, childNodes, d.delta)
			}
		}

	case selectorAllType:
		d := delta.delta.(*deltaSelectorAll)
		selector := d.selector.selector()

		if selector != nil {
			childNodes := make([]*html.Node, 0, len(nodes))
			for _, node := range nodes {
				ms := selector.MatchAll(node)
				for _, m := range ms {
					childNodes = append(childNodes, m)
				}
			}

			if len(childNodes) > 0 {
				return applyDelta(root, childNodes, d.delta)
			}
		}

	case parentType:
		d := delta.delta.(*deltaParent)
		childNodes := make([]*html.Node, 0, len(nodes))

		for _, node := range nodes {
			m := node.Parent
			if m != nil {
				childNodes = append(childNodes, m)
			}
		}

		if len(childNodes) > 0 {
			return applyDelta(root, childNodes, d.delta)
		}

	case firstChildType:
		d := delta.delta.(*deltaFirstChild)
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

		if len(childNodes) > 0 {
			return applyDelta(root, childNodes, d.delta)
		}

	case lastChildType:
		d := delta.delta.(*deltaLastChild)
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

		if len(childNodes) > 0 {
			return applyDelta(root, childNodes, d.delta)
		}

	case prevSiblingType:
		d := delta.delta.(*deltaPrevSibling)
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

		if len(childNodes) > 0 {
			return applyDelta(root, childNodes, d.delta)
		}

	case nextSiblingType:
		d := delta.delta.(*deltaNextSibling)
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

		if len(childNodes) > 0 {
			return applyDelta(root, childNodes, d.delta)
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

		childNodes := delta.delta.(*deltaHTML).factory.Nodes()
		for _, node := range nodes {
			if node.Type != html.ElementNode {
				continue
			}

			children := childNodes
			if len(nodes) > 1 {
				children = util.Clone(children)
			}

			for node.FirstChild != nil {
				node.RemoveChild(node.FirstChild)
			}

			for _, child := range children {
				node.AppendChild(child)
			}
		}

	case textType:

		text := delta.delta.(*deltaText).text
		for _, node := range nodes {
			if node.Type != html.ElementNode {
				continue
			}

			for node.FirstChild != nil {
				node.RemoveChild(node.FirstChild)
			}

			node.AppendChild(&html.Node{
				Type: html.TextNode,
				Data: text,
			})
		}

	case replaceType:

		childNodes := delta.delta.(*deltaReplace).factory.Nodes()
		for _, node := range nodes {
			if node.Type != html.ElementNode || node.Parent == nil {
				continue
			}

			children := childNodes
			if len(nodes) > 1 {
				children = util.Clone(children)
			}

			for _, child := range children {
				node.Parent.InsertBefore(child, node)
			}

			node.Parent.RemoveChild(node)
		}

	case appendType:

		childNodes := delta.delta.(*deltaAppend).factory.Nodes()
		for _, node := range nodes {
			if node.Type != html.ElementNode {
				continue
			}

			children := childNodes
			if len(nodes) > 1 {
				children = util.Clone(children)
			}

			for _, child := range children {
				node.AppendChild(child)
			}
		}

	case prependType:

		childNodes := delta.delta.(*deltaPrepend).factory.Nodes()
		for _, node := range nodes {
			if node.Type != html.ElementNode {
				continue
			}

			children := childNodes
			if len(nodes) > 1 {
				children = util.Clone(children)
			}

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

		childNodes := delta.delta.(*deltaInsertAfter).factory.Nodes()
		for _, node := range nodes {
			if node.Type != html.ElementNode || node.Parent == nil {
				continue
			}

			children := childNodes
			if len(nodes) > 1 {
				children = util.Clone(children)
			}

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

		childNodes := delta.delta.(*deltaInsertBefore).factory.Nodes()
		for _, node := range nodes {
			if node.Type != html.ElementNode || node.Parent == nil {
				continue
			}

			children := childNodes
			if len(nodes) > 1 {
				children = util.Clone(children)
			}

			for _, child := range children {
				node.Parent.InsertBefore(child, node)
			}
		}

	case addAttrType:
		attr := delta.delta.(*deltaAddAttr).attr
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
		attr := delta.delta.(*deltaAddAttr).attr
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
		attr := delta.delta.(*deltaRmAttr).attr
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
		styles := delta.delta.(*deltaAddStyles).styles

		for _, node := range nodes {
			if node.Type != html.ElementNode {
				continue
			}

			found := false
			for _, att := range node.Attr {
				if att.Namespace != "" {
					continue
				}

				if att.Key == "style" {
					parsed := parseStyle(att.Val)
					for key, value := range styles {
						parsed[key] = value
					}

					att.Val = buildStyle(parsed)
					found = true
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
		styles := delta.delta.(*deltaRmStyles).styles

		for _, node := range nodes {
			if node.Type != html.ElementNode {
				continue
			}

			for _, att := range node.Attr {
				if att.Namespace != "" {
					continue
				}

				if att.Key == "style" {
					parsed := parseStyle(att.Val)
					for _, s := range styles {
						delete(parsed, s)
					}

					att.Val = buildStyle(parsed)
				}
			}
		}

	case addClassType:
		class := delta.delta.(*deltaAddClass).class
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
		class := delta.delta.(*deltaRmClass).class
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

func parseStyle(style string) map[string]string {
	styleMap := map[string]string{}

	key := ""
	value := ""
	lookingForKey := true
	lookingForValue := false
	fillingKey := false
	fillingValue := false
	inSingleQuoteString := false
	inDoubleQuoteString := false
	escapedChar := false

	for _, r := range style {
	start:
		if fillingValue {
			if escapedChar {
				escapedChar = false
				value += string(r)
				continue
			}

			switch r {
			case '\\':
				escapedChar = true
				value += string(r)
				continue
			case '\'':
				if inSingleQuoteString {
					inSingleQuoteString = false
				} else {
					inSingleQuoteString = true
				}

				value += string(r)
				continue
			case '"':
				if inDoubleQuoteString {
					inDoubleQuoteString = false
				} else {
					inDoubleQuoteString = true
				}

				value += string(r)
				continue
			default:
				if inSingleQuoteString || inDoubleQuoteString {
					value += string(r)
					continue
				}
			}

			switch r {
			case ';':
				if key != "" {
					styleMap[key] = value
				}

				key = ""
				value = ""

				lookingForKey = true
				lookingForValue = false
				fillingKey = false
				fillingValue = false
			default:
				value += string(r)
			}

		} else {
			switch r {
			case ' ', '\t', '\r', '\n', '\f':
				fillingKey = false
			case ':':
				value = ""
				lookingForKey = false
				lookingForValue = true
				fillingKey = false
				fillingValue = false
			default:
				if lookingForKey || fillingKey {
					fillingKey = true
					lookingForKey = false
					key += string(r)
				} else if lookingForValue || fillingValue {
					fillingValue = true
					lookingForValue = false
					goto start
				}
			}
		}

	}

	if fillingValue && key != "" {
		styleMap[key] = value
	}

	return styleMap
}

func buildStyle(style map[string]string) string {
	attr := ""
	for key, value := range style {
		attr += key + ": " + value + ";"
	}

	return attr
}

func parseClass(class string) map[string]bool {
	currentClass := ""
	classes := map[string]bool{}

	flush := func() {
		if currentClass != "" {
			classes[currentClass] = true
			currentClass = ""
		}
	}

	for _, r := range class {
		switch r {
		case ' ', '\t', '\r', '\n', '\f':
			flush()
		default:
			currentClass += string(r)
		}
	}

	flush()
	return classes
}

func buildClass(classes map[string]bool) string {
	class := ""
	i := 0

	for key, value := range classes {
		if i != 0 {
			class += " "
		}

		if value {
			class += key
		}
	}

	return class
}

func discardDelta(delta Delta) {
}
