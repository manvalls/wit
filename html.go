package wit

import (
	"io"
	"net/http"
	"strings"

	"github.com/andybalholm/cascadia"
	"github.com/manvalls/wit/util"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const witCall = "(function(){var a=window.wit=window.wit||{},g=[];if(!a.call){a.call=function(c,b,d){var a=window;d=d||document.scripts[document.scripts.length-1].parentNode;var e;for(e=0;e<c.length;e++){var f=a;a=a[c[e]];if(!a){g.push([c,b,d]);return}}a.call(f,b,d)};var f=a.run;a.run=function(){var c=g,b;g=[];for(b=0;b<c.length;b++)try{a.call(c[b][0],c[b][1],c[b][2])}catch(d){setTimeout(function(){throw d;},0)}f&&f()}}})();"

var headSelector = cascadia.MustCompile("head")

var baseDocument, _ = html.Parse(strings.NewReader("<!DOCTYPE html><html><head></head><body></body></html>"))

type htmlContext struct {
	root        *html.Node
	loadWitCall bool
	deferred    []*deltaWithContext
	status      int
	headers     http.Header
	answer      io.ReadCloser
}

type deltaWithContext struct {
	delta Delta
	root  *html.Node
	nodes []*html.Node
}

// WriteHTML writes the result of applying the provided delta to an empty
// document as formatted HTML
func WriteHTML(w http.ResponseWriter, delta Delta) {
	nodes := util.Clone([]*html.Node{baseDocument})

	c := applyDelta(&htmlContext{
		root:    nodes[0],
		headers: make(http.Header),
	}, nodes, delta)

	for len(c.deferred) > 0 {
		deferred := c.deferred
		c.deferred = nil

		for _, def := range deferred {
			if def.root != c.root {
				discardDelta(def.delta)
			} else {
				c = applyDelta(c, def.nodes, def.delta)
			}
		}
	}

	headers := w.Header()
	for key, value := range c.headers {
		headers[key] = value
	}

	if c.status != 0 {
		w.WriteHeader(c.status)
	}

	if c.root != nil {
		if !c.loadWitCall {
			head := headSelector.MatchFirst(c.root)
			if head == nil {
				return
			}

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

		html.Render(w, c.root)
	} else if c.answer != nil {
		io.Copy(w, c.answer)
		c.answer.Close()
	}
}

func applyDelta(c *htmlContext, nodes []*html.Node, delta Delta) (next *htmlContext) {

	next = c

	switch delta.typeID {

	case sliceType:
		deltas := delta.delta.(*deltaSlice).deltas

		for _, childDelta := range deltas {
			if c.root != next.root {
				discardDelta(childDelta)
			} else {
				next = applyDelta(next, nodes, childDelta)
			}
		}

	case channelType:
		d := delta.delta.(*deltaChannel)

		channel := d.channel
		cancel := d.cancel

		for childDelta := range channel {
			if c.root != next.root {
				discardDelta(childDelta)
			} else {
				next = applyDelta(next, nodes, childDelta)
				if c.root != next.root {
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

	case jumpType:
		d := delta.delta.(*deltaJump).delta
		childNodes := util.Clone([]*html.Node{baseDocument})

		for _, def := range c.deferred {
			discardDelta(def.delta)
		}

		return applyDelta(&htmlContext{
			root:    childNodes[0],
			headers: make(http.Header),
		}, childNodes, d)

	case runSyncType:
		f := delta.delta.(*deltaRunSync).handler
		return applyDelta(c, nodes, f())

	case deferType:
		c.deferred = append(c.deferred, &deltaWithContext{
			root:  c.root,
			nodes: nodes,
			delta: delta.delta.(*deltaDefer).delta,
		})

	case statusType:
		c.status = delta.delta.(*deltaStatus).code

	case addHeadersType:
		headers := delta.delta.(*deltaAddHeaders).headers
		for key, value := range headers {
			for _, h := range value {
				c.headers.Add(key, h)
			}
		}

	case setHeadersType:
		headers := delta.delta.(*deltaSetHeaders).headers
		for key, value := range headers {
			c.headers[key] = value
		}

	case rmHeadersType:
		headers := delta.delta.(*deltaRmHeaders).headers
		for _, header := range headers {
			c.headers.Del(header)
		}

	case answerType:
		c.answer = delta.delta.(*deltaAnswer).reader
		c.root = nil

	}

	return
}
