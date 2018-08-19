package wit

import (
	"context"
	"errors"
	"io"

	"golang.org/x/net/html"
)

// List groups a list of deltas together
func List(deltas ...Delta) Delta {
	filteredDeltas := make([]Delta, 0, len(deltas))
	for _, delta := range deltas {
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

// Root applies given deltas to the root of the document
func Root(deltas ...Delta) Delta {
	d := List(deltas...)
	if d.typeID == 0 {
		return d
	}

	return Delta{rootType, &deltaRoot{d}}
}

// Run runs the given function under the given context, returning a delta
func Run(parentCtx context.Context, callback func(ctx context.Context) Delta) Delta {
	ctx, cancel := context.WithCancel(parentCtx)
	channel := make(chan Delta)

	go func() {
		select {
		case channel <- callback(ctx):
		case <-ctx.Done():
		}

		close(channel)
	}()

	return Delta{channelType, &deltaChannel{channel, cancel}}
}

// RunChannel runs the given function under the given context and channel, returning a delta
func RunChannel(parentCtx context.Context, callback func(ctx context.Context, ch chan<- Delta) Delta) Delta {
	ctx, cancel := context.WithCancel(parentCtx)
	channel := make(chan Delta)
	go callback(ctx, channel)
	return Delta{channelType, &deltaChannel{channel, cancel}}
}

// Nil represents an effectless delta
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
func HTML(factory Factory) Delta {
	return Delta{htmlType, &deltaHTML{factory}}
}

// Parent applies provided deltas to the parent of matching elements
func Parent(deltas ...Delta) Delta {
	d := List(deltas...)
	if d.typeID == 0 {
		return d
	}

	return Delta{parentType, &deltaParent{d}}
}

// FirstChild applies provided deltas to the first child of matching elements
func FirstChild(deltas ...Delta) Delta {
	d := List(deltas...)
	if d.typeID == 0 {
		return d
	}

	return Delta{firstChildType, &deltaFirstChild{d}}
}

// LastChild applies provided deltas to the last child of matching elements
func LastChild(deltas ...Delta) Delta {
	d := List(deltas...)
	if d.typeID == 0 {
		return d
	}

	return Delta{lastChildType, &deltaLastChild{d}}
}

// PrevSibling applies provided deltas to the previous sibling of matching elements
func PrevSibling(deltas ...Delta) Delta {
	d := List(deltas...)
	if d.typeID == 0 {
		return d
	}

	return Delta{prevSiblingType, &deltaPrevSibling{d}}
}

// NextSibling applies provided deltas to the previous sibling of matching elements
func NextSibling(deltas ...Delta) Delta {
	d := List(deltas...)
	if d.typeID == 0 {
		return d
	}

	return Delta{nextSiblingType, &deltaNextSibling{d}}
}

// Replace replaces matching elements with the provided HTML
func Replace(html Factory) Delta {
	return Delta{replaceType, &deltaReplace{html}}
}

// Append adds the provided HTML at the end of matching elements
func Append(html Factory) Delta {
	return Delta{appendType, &deltaAppend{html}}
}

// Prepend adds the provided HTML at the beginning of matching elements
func Prepend(html Factory) Delta {
	return Delta{prependType, &deltaPrepend{html}}
}

// InsertAfter inserts the provided HTML after matching elements
func InsertAfter(html Factory) Delta {
	return Delta{insertAfterType, &deltaInsertAfter{html}}
}

// InsertBefore inserts the provided HTML before matching elements
func InsertBefore(html Factory) Delta {
	return Delta{insertBeforeType, &deltaInsertBefore{html}}
}

// AddAttr adds the provided attributes to the matching elements
func AddAttr(attr map[string]string) Delta {
	if len(attr) == 0 {
		return Nil
	}

	return Delta{addAttrType, &deltaAddAttr{attr}}
}

// SetAttr sets the attributes of the matching elements
func SetAttr(attr map[string]string) Delta {
	return Delta{setAttrType, &deltaSetAttr{attr}}
}

// RmAttr removes the provided attributes from the matching elements
func RmAttr(attrs ...string) Delta {
	if len(attrs) == 0 {
		return Nil
	}

	return Delta{rmAttrType, &deltaRmAttr{attrs}}
}

// AddStyles adds the provided styles to the matching elements
func AddStyles(styles map[string]string) Delta {
	if len(styles) == 0 {
		return Nil
	}

	return Delta{addStylesType, &deltaAddStyles{styles}}
}

// RmStyles removes the provided styles from the matching elements
func RmStyles(styles ...string) Delta {
	if len(styles) == 0 {
		return Nil
	}

	return Delta{rmStylesType, &deltaRmStyles{styles}}
}

// AddClass adds the provided class to the matching elements
func AddClass(class string) Delta {
	return Delta{addClassType, &deltaAddClass{class}}
}

// RmClass adds the provided class to the matching elements
func RmClass(class string) Delta {
	return Delta{rmClassType, &deltaRmClass{class}}
}

// Call calls a JavaScript function with provided parameters, when it becomes available
func Call(path []string, args map[string]string) Delta {
	return Delta{callType, &deltaCall{path, args}}
}

// Error stops the delta flow and throws the given error
func Error(err error) Delta {
	return Delta{errorType, &deltaError{err}}
}

// RunSync runs the given function synchronously, applying returned delta
func RunSync(handler func() Delta) Delta {
	return Delta{runSyncType, &deltaRunSync{handler}}
}

// ErrEnd is a generic error used to signal that the normal flow should be stopped
var ErrEnd = errors.New("Delta flow was ended")

// End throws the ErrEnd error
var End = Error(ErrEnd)
