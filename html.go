package wit

import (
	"io"
	"strings"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const witCall = "(function(){var a=window.wit=window.wit||{},g=[];if(!a.call){a.call=function(c,b,d){var a=window;d=d||document.scripts[document.scripts.length-1].parentNode;var e;for(e=0;e<c.length;e++){var f=a;a=a[c[e]];if(!a){g.push([c,b,d]);return}}a.call(f,b,d)};var f=a.run;a.run=function(){var c=g,b;g=[];for(b=0;b<c.length;b++)try{a.call(c[b][0],c[b][1],c[b][2])}catch(d){setTimeout(function(){throw d;},0)}f&&f()}}})();"

var headSelector = cascadia.MustCompile("head")

var baseDocument, _ = html.Parse(strings.NewReader("<!DOCTYPE html><html><head></head><body></body></html>"))

type htmlContext struct {
	root        *html.Node
	loadWitCall bool
}

type htmlRenderer struct {
	root *html.Node
}

// NewHTMLRenderer returns a new renderer which will render HTML
func NewHTMLRenderer(delta Delta) (Renderer, error) {
	nodes := clone([]*html.Node{baseDocument})
	c := &htmlContext{
		root: nodes[0],
	}

	err := applyDelta(c, nodes, delta)
	if err != nil {
		return nil, err
	}

	if c.loadWitCall {
		head := headSelector.MatchFirst(c.root)
		if head != nil {
			script := &html.Node{
				Type:      html.ElementNode,
				DataAtom:  atom.Script,
				Data:      "script",
				Namespace: "",
			}

			if head.FirstChild != nil {
				head.InsertBefore(script, head.FirstChild)
			} else {
				head.AppendChild(script)
			}

			script.AppendChild(&html.Node{
				Type: html.TextNode,
				Data: witCall,
			})
		}
	}

	return &htmlRenderer{c.root}, nil
}

func (r *htmlRenderer) Render(w io.Writer) error {
	return html.Render(w, r.root)
}

func applyDelta(c *htmlContext, nodes []*html.Node, delta Delta) (err error) {

	switch delta.typeID {

	case sliceType:
		deltas := delta.delta.(*deltaSlice).deltas

		for _, childDelta := range deltas {
			if err != nil {
				Discard(childDelta)
			} else {
				err = applyDelta(c, nodes, childDelta)
			}
		}

	case channelType:
		d := delta.delta.(*deltaChannel)

		channel := d.channel
		cancel := d.cancel

		for childDelta := range channel {
			if err != nil {
				Discard(childDelta)
			} else {
				err = applyDelta(c, nodes, childDelta)
				if err != nil {
					cancel()
				}
			}
		}

	case rootType:
		return applyDelta(c, []*html.Node{c.root}, delta.delta.(*deltaRoot).delta)

	case selectorType:
		d := delta.delta.(*deltaSelector)
		selector := d.selector.selector()
		childNodes := make([]*html.Node, 0, len(nodes))

		if selector != nil {
			for _, node := range nodes {
				m := selector.MatchFirst(node)
				if m != nil {
					childNodes = append(childNodes, m)
				}
			}
		}

		return applyDelta(c, childNodes, d.delta)

	case selectorAllType:
		d := delta.delta.(*deltaSelectorAll)
		selector := d.selector.selector()
		childNodes := make([]*html.Node, 0, len(nodes))

		if selector != nil {
			for _, node := range nodes {
				ms := selector.MatchAll(node)
				for _, m := range ms {
					childNodes = append(childNodes, m)
				}
			}
		}

		return applyDelta(c, childNodes, d.delta)

	case parentType:
		d := delta.delta.(*deltaParent)
		childNodes := make([]*html.Node, 0, len(nodes))

		for _, node := range nodes {
			m := node.Parent
			if m != nil {
				childNodes = append(childNodes, m)
			}
		}

		return applyDelta(c, childNodes, d.delta)

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

		return applyDelta(c, childNodes, d.delta)

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

		return applyDelta(c, nodes, d.delta)

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

		return applyDelta(c, childNodes, d.delta)

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

		return applyDelta(c, nodes, d.delta)

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

			children := delta.delta.(*deltaHTML).factory.Nodes(node)
			for _, child := range children {
				node.AppendChild(child)
			}
		}

	case replaceType:

		for _, node := range nodes {
			if node.Type != html.ElementNode || node.Parent == nil {
				continue
			}

			children := delta.delta.(*deltaHTML).factory.Nodes(node)
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

			children := delta.delta.(*deltaHTML).factory.Nodes(node)

			for _, child := range children {
				node.AppendChild(child)
			}
		}

	case prependType:

		for _, node := range nodes {
			if node.Type != html.ElementNode {
				continue
			}

			children := delta.delta.(*deltaHTML).factory.Nodes(node)

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

			children := delta.delta.(*deltaHTML).factory.Nodes(node)

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

			children := delta.delta.(*deltaHTML).factory.Nodes(node)
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
		attr := delta.delta.(*deltaSetAttr).attr
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
		styles := delta.delta.(*deltaRmStyles).styles

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

	case callType:
		d := delta.delta.(*deltaCall)
		c.loadWitCall = true

		for _, node := range nodes {
			if node.Type != html.ElementNode {
				continue
			}

			script := &html.Node{
				Type:      html.ElementNode,
				DataAtom:  atom.Script,
				Data:      "script",
				Namespace: "",
			}

			node.AppendChild(script)

			script.AppendChild(&html.Node{
				Type: html.TextNode,
				Data: "wit.call(" + strSliceToJSON(d.path) + "," + strMapToJSON(d.arguments) + ");",
			})
		}

	case errorType:
		err = delta.delta.(*deltaError).err

	case runSyncType:
		f := delta.delta.(*deltaRunSync).handler
		return applyDelta(c, nodes, f())

	}

	return
}
