package wit

import (
	"io"

	"golang.org/x/net/html"
)

// List groups a list of actions together
func List(actions ...Action) Action {
	filteredDeltas := make([]Delta, 0, len(actions))
	for _, action := range actions {
		delta := action.Delta()

		switch delta.typeID {
		case 0:
		case sliceType:
			childDeltas := delta.delta.(*deltaSlice).deltas
			for _, childDelta := range childDeltas {
				filteredDeltas = append(filteredDeltas, childDelta)
			}
		default:
			filteredDeltas = append(filteredDeltas, delta)
		}
	}

	switch len(filteredDeltas) {
	case 0:
		return Nil
	case 1:
		return filteredDeltas[0]
	default:
		return Delta{sliceType, &deltaSlice{filteredDeltas}}
	}
}

// Root applies given actions to the root of the document
func Root(actions ...Action) Action {
	d := List(actions...).Delta()
	if d.typeID == 0 {
		return d
	}

	return Delta{rootType, &deltaRoot{d}}
}

// Nil represents an effectless action
var Nil = Delta{}

// Remove removes from the document matching elements
var Remove = Delta{removeType, &deltaRemove{}}

// Clear empties matching elements
var Clear = Delta{clearType, &deltaClear{}}

// Factory builds HTML documents on demand
type Factory interface {
	HTML() io.Reader
	Nodes(context *html.Node) []*html.Node
}

// HTML sets the inner HTML of the matching elements
func HTML(factory Factory) Action {
	return Delta{htmlType, &deltaHTML{factory}}
}

// Parent applies provided actions to the parent of matching elements
func Parent(actions ...Action) Action {
	d := List(actions...).Delta()
	if d.typeID == 0 {
		return d
	}

	return Delta{parentType, &deltaParent{d}}
}

// FirstChild applies provided actions to the first child of matching elements
func FirstChild(actions ...Action) Action {
	d := List(actions...).Delta()
	if d.typeID == 0 {
		return d
	}

	return Delta{firstChildType, &deltaFirstChild{d}}
}

// LastChild applies provided actions to the last child of matching elements
func LastChild(actions ...Action) Action {
	d := List(actions...).Delta()
	if d.typeID == 0 {
		return d
	}

	return Delta{lastChildType, &deltaLastChild{d}}
}

// PrevSibling applies provided actions to the previous sibling of matching elements
func PrevSibling(actions ...Action) Action {
	d := List(actions...).Delta()
	if d.typeID == 0 {
		return d
	}

	return Delta{prevSiblingType, &deltaPrevSibling{d}}
}

// NextSibling applies provided actions to the previous sibling of matching elements
func NextSibling(actions ...Action) Action {
	d := List(actions...).Delta()
	if d.typeID == 0 {
		return d
	}

	return Delta{nextSiblingType, &deltaNextSibling{d}}
}

// Replace replaces matching elements with the provided HTML
func Replace(html Factory) Action {
	return Delta{replaceType, &deltaReplace{html}}
}

// Append adds the provided HTML at the end of matching elements
func Append(html Factory) Action {
	return Delta{appendType, &deltaAppend{html}}
}

// Prepend adds the provided HTML at the beginning of matching elements
func Prepend(html Factory) Action {
	return Delta{prependType, &deltaPrepend{html}}
}

// InsertAfter inserts the provided HTML after matching elements
func InsertAfter(html Factory) Action {
	return Delta{insertAfterType, &deltaInsertAfter{html}}
}

// InsertBefore inserts the provided HTML before matching elements
func InsertBefore(html Factory) Action {
	return Delta{insertBeforeType, &deltaInsertBefore{html}}
}

// AddAttr adds the provided attributes to the matching elements
func AddAttr(attr map[string]string) Action {
	if len(attr) == 0 {
		return Nil
	}

	return Delta{addAttrType, &deltaAddAttr{attr}}
}

// SetAttr sets the attributes of the matching elements
func SetAttr(attr map[string]string) Action {
	return Delta{setAttrType, &deltaSetAttr{attr}}
}

// RmAttr removes the provided attributes from the matching elements
func RmAttr(attrs ...string) Action {
	if len(attrs) == 0 {
		return Nil
	}

	return Delta{rmAttrType, &deltaRmAttr{attrs}}
}

// AddStyles adds the provided styles to the matching elements
func AddStyles(styles map[string]string) Action {
	if len(styles) == 0 {
		return Nil
	}

	return Delta{addStylesType, &deltaAddStyles{styles}}
}

// RmStyles removes the provided styles from the matching elements
func RmStyles(styles ...string) Action {
	if len(styles) == 0 {
		return Nil
	}

	return Delta{rmStylesType, &deltaRmStyles{styles}}
}

// AddClass adds the provided class to the matching elements
func AddClass(class string) Action {
	return Delta{addClassType, &deltaAddClass{class}}
}

// RmClass adds the provided class to the matching elements
func RmClass(class string) Action {
	return Delta{rmClassType, &deltaRmClass{class}}
}
