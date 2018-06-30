package wit

import (
	"net/http"
	"strings"

	"github.com/andybalholm/cascadia"
	"github.com/manvalls/wit/util"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const witCall = "(function(){var a=window.wit=window.wit||{},g=[];if(!a.call){a.call=function(c,b,d){var a=window;d=d||document.scripts[document.scripts.length-1].parentNode;var e;for(e=0;e<c.length;e++){var f=a;a=a[c[e]];if(!a){g.push([c,b,d]);return}}a.call(f,b,d)};var f=a.run;a.run=function(){var c=g,b;g=[];for(b=0;b<c.length;b++)try{a.call(c[b][0],c[b][1],c[b][2])}catch(d){setTimeout(function(){throw d;},0)}f&&f()}}})();"

const loadCSS = "/* loadCSS. [c]2017 Filament Group, Inc. MIT License */(function(e){var d={};d.c=function(){try{var a=e.document.createElement(\"link\").f.supports(\"preload\")}catch(b){a=!1}return function(){return a}}();d.a=function(a){function b(){a.media=c}var c=a.media||\"all\";a.addEventListener?a.addEventListener(\"load\",b):a.attachEvent&&a.attachEvent(\"onload\",b);setTimeout(function(){a.rel=\"stylesheet\";a.media=\"only x\"});setTimeout(b,3E3)};d.b=function(){for(var a=e.document.getElementsByTagName(\"link\"),b=0;b<a.length;b++){var c=a[b];\"preload\"!==c.rel||\"style\"!==c.getAttribute(\"as\")||c.getAttribute(\"data-loadcss\")||(c.setAttribute(\"data-loadcss\",!0),d.a(c))}};d.c()||d.b()})(window);"

var headSelector = cascadia.MustCompile("head")

var baseDocument, _ = html.Parse(strings.NewReader("<!DOCTYPE html><html><head></head><body></body></html>"))

type htmlContext struct {
	root            *html.Node
	nodes           []*html.Node
	loadCSSPolyfill *html.Node
	isWitCallLoaded bool
}

// WriteHTML writes the result of applying the provided delta to an empty
// document as formatted HTML
func WriteHTML(w http.ResponseWriter, delta Delta) error {
	nodes := util.Clone([]*html.Node{baseDocument})

	c := applyDelta(&htmlContext{
		root:  nodes[0],
		nodes: nodes,
	}, delta)

	if c.root != nil {
		return html.Render(w, c.root)
	}

	return nil
}

func applyDelta(c *htmlContext, delta Delta) (next *htmlContext) {

	next = c

	switch delta.typeID {

	case sliceType:
		deltas := delta.delta.(*deltaSlice).deltas

		for _, childDelta := range deltas {
			if c.root != next.root {
				discardDelta(childDelta)
			} else {
				next = applyDelta(next, childDelta)
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
				next = applyDelta(next, childDelta)
				if c.root != next.root {
					cancel()
				}
			}
		}

	case rootType:
		childCtx := *c
		childCtx.nodes = []*html.Node{c.root}
		return applyDelta(&childCtx, delta.delta.(*deltaRoot).delta)

	case selectorType:
		d := delta.delta.(*deltaSelector)
		selector := d.selector.selector()

		if selector != nil {
			childNodes := make([]*html.Node, 0, len(c.nodes))
			for _, node := range c.nodes {
				m := selector.MatchFirst(node)
				if m != nil {
					childNodes = append(childNodes, m)
				}
			}

			if len(childNodes) > 0 {
				childCtx := *c
				childCtx.nodes = childNodes
				return applyDelta(&childCtx, d.delta)
			}
		}

	case selectorAllType:
		d := delta.delta.(*deltaSelectorAll)
		selector := d.selector.selector()

		if selector != nil {
			childNodes := make([]*html.Node, 0, len(c.nodes))
			for _, node := range c.nodes {
				ms := selector.MatchAll(node)
				for _, m := range ms {
					childNodes = append(childNodes, m)
				}
			}

			if len(childNodes) > 0 {
				childCtx := *c
				childCtx.nodes = childNodes
				return applyDelta(&childCtx, d.delta)
			}
		}

	case parentType:
		d := delta.delta.(*deltaParent)
		childNodes := make([]*html.Node, 0, len(c.nodes))

		for _, node := range c.nodes {
			m := node.Parent
			if m != nil {
				childNodes = append(childNodes, m)
			}
		}

		if len(childNodes) > 0 {
			childCtx := *c
			childCtx.nodes = childNodes
			return applyDelta(&childCtx, d.delta)
		}

	case firstChildType:
		d := delta.delta.(*deltaFirstChild)
		childNodes := make([]*html.Node, 0, len(c.nodes))

		for _, node := range c.nodes {
			m := node.FirstChild
			for m != nil && m.Type != html.ElementNode {
				m = m.NextSibling
			}

			if m != nil {
				childNodes = append(childNodes, m)
			}
		}

		if len(childNodes) > 0 {
			childCtx := *c
			childCtx.nodes = childNodes
			return applyDelta(&childCtx, d.delta)
		}

	case lastChildType:
		d := delta.delta.(*deltaLastChild)
		childNodes := make([]*html.Node, 0, len(c.nodes))

		for _, node := range c.nodes {
			m := node.LastChild
			for m != nil && m.Type != html.ElementNode {
				m = m.PrevSibling
			}

			if m != nil {
				childNodes = append(childNodes, m)
			}
		}

		if len(childNodes) > 0 {
			childCtx := *c
			childCtx.nodes = childNodes
			return applyDelta(&childCtx, d.delta)
		}

	case prevSiblingType:
		d := delta.delta.(*deltaPrevSibling)
		childNodes := make([]*html.Node, 0, len(c.nodes))

		for _, node := range c.nodes {
			m := node.PrevSibling
			for m != nil && m.Type != html.ElementNode {
				m = m.PrevSibling
			}

			if m != nil {
				childNodes = append(childNodes, m)
			}
		}

		if len(childNodes) > 0 {
			childCtx := *c
			childCtx.nodes = childNodes
			return applyDelta(&childCtx, d.delta)
		}

	case nextSiblingType:
		d := delta.delta.(*deltaNextSibling)
		childNodes := make([]*html.Node, 0, len(c.nodes))

		for _, node := range c.nodes {
			m := node.NextSibling
			for m != nil && m.Type != html.ElementNode {
				m = m.NextSibling
			}

			if m != nil {
				childNodes = append(childNodes, m)
			}
		}

		if len(childNodes) > 0 {
			childCtx := *c
			childCtx.nodes = childNodes
			return applyDelta(&childCtx, d.delta)
		}

	case removeType:

		for _, node := range c.nodes {
			parent := node.Parent
			if parent != nil {
				parent.RemoveChild(node)
			}
		}

	case clearType:

		for _, node := range c.nodes {
			if node.Type == html.ElementNode {
				for node.FirstChild != nil {
					node.RemoveChild(node.FirstChild)
				}
			}
		}

	case htmlType:

		childNodes := delta.delta.(*deltaHTML).factory.Nodes()
		for _, node := range c.nodes {
			if node.Type != html.ElementNode {
				continue
			}

			children := childNodes
			if len(c.nodes) > 1 {
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
		for _, node := range c.nodes {
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
		for _, node := range c.nodes {
			if node.Type != html.ElementNode || node.Parent == nil {
				continue
			}

			children := childNodes
			if len(c.nodes) > 1 {
				children = util.Clone(children)
			}

			for _, child := range children {
				node.Parent.InsertBefore(child, node)
			}

			node.Parent.RemoveChild(node)
		}

	case appendType:

		childNodes := delta.delta.(*deltaAppend).factory.Nodes()
		for _, node := range c.nodes {
			if node.Type != html.ElementNode {
				continue
			}

			children := childNodes
			if len(c.nodes) > 1 {
				children = util.Clone(children)
			}

			for _, child := range children {
				node.AppendChild(child)
			}
		}

	case prependType:

		childNodes := delta.delta.(*deltaPrepend).factory.Nodes()
		for _, node := range c.nodes {
			if node.Type != html.ElementNode {
				continue
			}

			children := childNodes
			if len(c.nodes) > 1 {
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
		for _, node := range c.nodes {
			if node.Type != html.ElementNode || node.Parent == nil {
				continue
			}

			children := childNodes
			if len(c.nodes) > 1 {
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
		for _, node := range c.nodes {
			if node.Type != html.ElementNode || node.Parent == nil {
				continue
			}

			children := childNodes
			if len(c.nodes) > 1 {
				children = util.Clone(children)
			}

			for _, child := range children {
				node.Parent.InsertBefore(child, node)
			}
		}

	case addAttrType:
		attr := delta.delta.(*deltaAddAttr).attr
		for _, node := range c.nodes {
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
		for _, node := range c.nodes {
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

		for _, node := range c.nodes {
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

		for _, node := range c.nodes {
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

		for _, node := range c.nodes {
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

		for _, node := range c.nodes {
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

		for _, node := range c.nodes {
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

	case jsType:
		url := delta.delta.(*deltaJS).url
		head := headSelector.MatchFirst(c.root)

		if head != nil {
			head.AppendChild(&html.Node{
				Type:      html.ElementNode,
				DataAtom:  atom.Script,
				Data:      "script",
				Namespace: "",
				Attr: []html.Attribute{
					html.Attribute{
						Key: "src",
						Val: url,
					},
				},
			})
		}

	case asyncJSType:
		url := delta.delta.(*deltaAsyncJS).url
		head := headSelector.MatchFirst(c.root)

		if head != nil {
			head.AppendChild(&html.Node{
				Type:      html.ElementNode,
				DataAtom:  atom.Script,
				Data:      "script",
				Namespace: "",
				Attr: []html.Attribute{
					html.Attribute{
						Key: "src",
						Val: url,
					},
					html.Attribute{
						Key: "async",
					},
				},
			})
		}

	case cssType:
		url := delta.delta.(*deltaCSS).url
		head := headSelector.MatchFirst(c.root)

		if head != nil {
			head.AppendChild(&html.Node{
				Type:      html.ElementNode,
				DataAtom:  atom.Link,
				Data:      "link",
				Namespace: "",
				Attr: []html.Attribute{
					html.Attribute{
						Key: "rel",
						Val: "stylesheet",
					},
					html.Attribute{
						Key: "type",
						Val: "text/css",
					},
					html.Attribute{
						Key: "href",
						Val: url,
					},
				},
			})
		}

	case asyncCSSType:
		url := delta.delta.(*deltaAsyncCSS).url

		if c.loadCSSPolyfill == nil {
			head := headSelector.MatchFirst(c.root)
			if head == nil {
				return
			}

			c.loadCSSPolyfill = &html.Node{
				Type:      html.ElementNode,
				DataAtom:  atom.Script,
				Data:      "script",
				Namespace: "",
			}

			head.AppendChild(c.loadCSSPolyfill)

			c.loadCSSPolyfill.AppendChild(&html.Node{
				Type: html.TextNode,
				Data: loadCSS,
			})
		}

		parent := c.loadCSSPolyfill.Parent
		if parent != nil {
			parent.InsertBefore(&html.Node{
				Type:      html.ElementNode,
				DataAtom:  atom.Link,
				Data:      "link",
				Namespace: "",
				Attr: []html.Attribute{
					html.Attribute{
						Key: "rel",
						Val: "preload",
					},
					html.Attribute{
						Key: "href",
						Val: url,
					},
					html.Attribute{
						Key: "as",
						Val: "style",
					},
					html.Attribute{
						Key: "onload",
						Val: "this.onload=null;this.rel='stylesheet'",
					},
				},
			}, c.loadCSSPolyfill)

			noscript := &html.Node{
				Type:      html.ElementNode,
				DataAtom:  atom.Noscript,
				Data:      "noscript",
				Namespace: "",
			}

			parent.InsertBefore(noscript, c.loadCSSPolyfill)
			noscript.AppendChild(&html.Node{
				Type:      html.ElementNode,
				DataAtom:  atom.Link,
				Data:      "link",
				Namespace: "",
				Attr: []html.Attribute{
					html.Attribute{
						Key: "rel",
						Val: "stylesheet",
					},
					html.Attribute{
						Key: "type",
						Val: "text/css",
					},
					html.Attribute{
						Key: "href",
						Val: url,
					},
				},
			})
		}

	case callType:
		d := delta.delta.(*deltaCall)

		if !c.isWitCallLoaded {
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

		for _, node := range c.nodes {
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
				Data: "wit.call(" + pathToJSON(d.path) + "," + argsToJSON(d.arguments) + ")",
			})
		}

	case jumpType:
		d := delta.delta.(*deltaJump).delta
		nodes := util.Clone([]*html.Node{baseDocument})

		return applyDelta(&htmlContext{
			root:  nodes[0],
			nodes: nodes,
		}, d)

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

func escape(str string) string {
	return strings.Replace(
		strings.Replace(str, "\\", "\\\\", -1),
		"\"",
		"\\\"",
		-1,
	)
}

func pathToJSON(arr []string) string {
	result := "["

	for i, str := range arr {
		if i != 0 {
			result += ","
		}

		result += "\"" + escape(str) + "\""
	}

	result += "]"
	return result
}

func argsToJSON(args map[string]string) string {
	result := "{"
	i := 0

	for key, value := range args {
		if i != 0 {
			result += ","
		}

		result += "\"" + escape(key) + "\":\"" + escape(value) + "\""
		i++
	}

	result += "}"
	return result
}

func discardDelta(delta Delta) {
}
