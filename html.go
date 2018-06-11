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
			m := node.PrevSibling
			if m != nil {
				childNodes = append(childNodes, m)
			}
		}

		if len(childNodes) > 0 {
			return applyDelta(root, childNodes, d.delta)
		}

	}

	return root, false
}

func discardDelta(delta Delta) {

}
