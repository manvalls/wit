package wit

import (
	"io"

	"golang.org/x/net/html"
)

func extractDeltas(actions []Action) []Delta {
	filteredDeltas := make([]Delta, 0, len(actions))
	for _, action := range actions {
		delta := action.Delta()

		switch delta.typeID {
		case 0:
		case sliceType:
			childDeltas := delta.info.deltas
			for _, childDelta := range childDeltas {
				filteredDeltas = append(filteredDeltas, childDelta)
			}
		default:
			filteredDeltas = append(filteredDeltas, delta)
		}
	}

	return filteredDeltas
}

// List groups a list of actions together
func List(actions ...Action) Action {
	filteredDeltas := extractDeltas(actions)

	switch len(filteredDeltas) {
	case 0:
		return Nil
	case 1:
		return filteredDeltas[0]
	default:
		return Delta{sliceType, &deltaInfo{deltas: filteredDeltas}}
	}
}

// Root applies given actions to the root of the document
func Root(actions ...Action) Action {
	deltas := extractDeltas(actions)
	if len(deltas) == 0 {
		return Nil
	}

	return Delta{rootType, &deltaInfo{deltas: deltas}}
}

// Nil represents an effectless action
var Nil = Delta{}

// Remove removes from the document matching elements
var Remove = Delta{removeType, nil}

// Clear empties matching elements
var Clear = Delta{clearType, nil}

// Factory builds HTML documents on demand
type Factory interface {
	HTML() io.Reader
	Nodes(context *html.Node) []*html.Node
}

// HTML sets the inner HTML of the matching elements
func HTML(factory Factory) Action {
	return Delta{htmlType, &deltaInfo{factory: factory}}
}

// Parent applies provided actions to the parent of matching elements
func Parent(actions ...Action) Action {
	deltas := extractDeltas(actions)
	if len(deltas) == 0 {
		return Nil
	}

	return Delta{parentType, &deltaInfo{deltas: deltas}}
}

// FirstChild applies provided actions to the first child of matching elements
func FirstChild(actions ...Action) Action {
	deltas := extractDeltas(actions)
	if len(deltas) == 0 {
		return Nil
	}

	return Delta{firstChildType, &deltaInfo{deltas: deltas}}
}

// LastChild applies provided actions to the last child of matching elements
func LastChild(actions ...Action) Action {
	deltas := extractDeltas(actions)
	if len(deltas) == 0 {
		return Nil
	}

	return Delta{lastChildType, &deltaInfo{deltas: deltas}}
}

// PrevSibling applies provided actions to the previous sibling of matching elements
func PrevSibling(actions ...Action) Action {
	deltas := extractDeltas(actions)
	if len(deltas) == 0 {
		return Nil
	}

	return Delta{prevSiblingType, &deltaInfo{deltas: deltas}}
}

// NextSibling applies provided actions to the previous sibling of matching elements
func NextSibling(actions ...Action) Action {
	deltas := extractDeltas(actions)
	if len(deltas) == 0 {
		return Nil
	}

	return Delta{nextSiblingType, &deltaInfo{deltas: deltas}}
}

// Replace replaces matching elements with the provided HTML
func Replace(factory Factory) Action {
	return Delta{replaceType, &deltaInfo{factory: factory}}
}

// Append adds the provided HTML at the end of matching elements
func Append(factory Factory) Action {
	return Delta{appendType, &deltaInfo{factory: factory}}
}

// Prepend adds the provided HTML at the beginning of matching elements
func Prepend(factory Factory) Action {
	return Delta{prependType, &deltaInfo{factory: factory}}
}

// InsertAfter inserts the provided HTML after matching elements
func InsertAfter(factory Factory) Action {
	return Delta{insertAfterType, &deltaInfo{factory: factory}}
}

// InsertBefore inserts the provided HTML before matching elements
func InsertBefore(factory Factory) Action {
	return Delta{insertBeforeType, &deltaInfo{factory: factory}}
}

// AddAttr adds the provided attributes to the matching elements
func AddAttr(attr map[string]string) Action {
	if len(attr) == 0 {
		return Nil
	}

	return Delta{addAttrType, &deltaInfo{strMap: attr}}
}

// SetAttr sets the attributes of the matching elements
func SetAttr(attr map[string]string) Action {
	return Delta{setAttrType, &deltaInfo{strMap: attr}}
}

// RmAttr removes the provided attributes from the matching elements
func RmAttr(attrs ...string) Action {
	if len(attrs) == 0 {
		return Nil
	}

	return Delta{rmAttrType, &deltaInfo{strList: attrs}}
}

// AddStyles adds the provided styles to the matching elements
func AddStyles(styles map[string]string) Action {
	if len(styles) == 0 {
		return Nil
	}

	return Delta{addStylesType, &deltaInfo{strMap: styles}}
}

// RmStyles removes the provided styles from the matching elements
func RmStyles(styles ...string) Action {
	if len(styles) == 0 {
		return Nil
	}

	return Delta{rmStylesType, &deltaInfo{strList: styles}}
}

// AddClass adds the provided class to the matching elements
func AddClass(class string) Action {
	return Delta{addClassType, &deltaInfo{class: class}}
}

// RmClass adds the provided class to the matching elements
func RmClass(class string) Action {
	return Delta{rmClassType, &deltaInfo{class: class}}
}
