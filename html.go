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

	switch delta.typeID {

	case sliceType:
		deltas := delta.delta.(*deltaSlice).deltas
		discard := false

		for _, childDelta := range deltas {
			if discard {
				discardDelta(childDelta)
			} else {
				root, discard = applyDelta(root, nodes, childDelta)
			}
		}

		return root, discard

	case channelType:
		d := delta.delta.(*deltaChannel)

		channel := d.channel
		cancel := d.cancel
		discard := false

		for childDelta := range channel {
			if discard {
				discardDelta(childDelta)
			} else {
				root, discard = applyDelta(root, nodes, childDelta)
				if discard {
					cancel()
				}
			}
		}

		return root, discard

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

	}

	return root, false
}

func discardDelta(delta Delta) {

}
