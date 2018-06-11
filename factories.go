package wit

import (
	"context"
	"io"

	"golang.org/x/net/html"
)

// List groups a list of deltas together
func List(deltas ...Delta) Delta {
	return Delta{sliceType, &deltaSlice{deltas}}
}

// Root applies given deltas to the root of the document
func Root(deltas ...Delta) Delta {
	return Delta{rootType, &deltaRoot{List(deltas...)}}
}

// Run runs the given function under the given context, returning a delta
func Run(parentCtx context.Context, callback func(context.Context) Delta) Delta {
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
func RunChannel(parentCtx context.Context, callback func(context.Context, chan<- Delta) Delta) Delta {
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

// HTMLFactory builds HTML documents on demand
type HTMLFactory interface {
	HTML() io.Reader
	Nodes() []*html.Node
}

// HTML sets the inner HTML of the matching elements
func HTML(factory HTMLFactory) Delta {
	return Delta{htmlType, &deltaHTML{factory}}
}

// Text sets the inner text of the matching elements
func Text(txt string) Delta {
	return Delta{textType, &deltaText{txt}}
}

// Parent applies provided deltas to the parent of matching elements
func Parent(deltas ...Delta) Delta {
	return Delta{parentType, &deltaParent{List(deltas...)}}
}

// FirstChild applies provided deltas to the first child of matching elements
func FirstChild(deltas ...Delta) Delta {
	return Delta{firstChildType, &deltaFirstChild{List(deltas...)}}
}

// LastChild applies provided deltas to the last child of matching elements
func LastChild(deltas ...Delta) Delta {
	return Delta{lastChildType, &deltaLastChild{List(deltas...)}}
}

// PrevSibling applies provided deltas to the previous sibling of matching elements
func PrevSibling(deltas ...Delta) Delta {
	return Delta{prevSiblingType, &deltaPrevSibling{List(deltas...)}}
}

// NextSibling applies provided deltas to the previous sibling of matching elements
func NextSibling(deltas ...Delta) Delta {
	return Delta{nextSiblingType, &deltaNextSibling{List(deltas...)}}
}

// Replace replaces matching elements with the provided HTML
func Replace(html HTMLFactory) Delta {
	return Delta{replaceType, &deltaReplace{html}}
}

// Append adds the provided HTML at the end of matching elements
func Append(html HTMLFactory) Delta {
	return Delta{appendType, &deltaAppend{html}}
}

// Prepend adds the provided HTML at the beginning of matching elements
func Prepend(html HTMLFactory) Delta {
	return Delta{prependType, &deltaPrepend{html}}
}

// InsertAfter inserts the provided HTML after matching elements
func InsertAfter(html HTMLFactory) Delta {
	return Delta{insertAfterType, &deltaInsertAfter{html}}
}

// InsertBefore inserts the provided HTML before matching elements
func InsertBefore(html HTMLFactory) Delta {
	return Delta{insertBeforeType, &deltaInsertBefore{html}}
}

// AddAttr adds the provided attributes to the matching elements
func AddAttr(attr map[string]string) Delta {
	return Delta{addAttrType, &deltaAddAttr{attr}}
}

// SetAttr sets the attributes of the matching elements
func SetAttr(attr map[string]string) Delta {
	return Delta{setAttrType, &deltaSetAttr{attr}}
}

// RmAttr removes the provided attributes from the matching elements
func RmAttr(attrs ...string) Delta {
	return Delta{rmAttrType, &deltaRmAttr{attrs}}
}

// AddStyles adds the provided styles to the matching elements
func AddStyles(styles map[string]string) Delta {
	return Delta{addStylesType, &deltaAddStyles{styles}}
}

// RmStyles removes the provided styles from the matching elements
func RmStyles(styles ...string) Delta {
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

// JS loads the provided script synchronously
func JS(key string, url string) Delta {
	return Delta{jsType, &deltaJS{key, url}}
}

// AsyncJS loads the provided script asynchronously
func AsyncJS(key string, url string) Delta {
	return Delta{asyncJSType, &deltaAsyncJS{key, url}}
}

// CSS loads the provided script synchronously
func CSS(key string, url string) Delta {
	return Delta{cssType, &deltaCSS{key, url}}
}

// AsyncCSS loads the provided script asynchronously
func AsyncCSS(key string, url string) Delta {
	return Delta{asyncCSSType, &deltaAsyncCSS{key, url}}
}

// Call calls a JavaScript function with provided parameters, when it becomes available
func Call(path []string, args map[string]string) Delta {
	return Delta{callType, &deltaCall{path, args}}
}

// Jump discards all deltas present and future and applies the given delta to the document
func Jump(delta Delta) Delta {
	return Delta{jumpType, &deltaJump{delta}}
}

// RunSync runs the given function synchronously, applying returned delta
func RunSync(handler func() Delta) Delta {
	return Delta{runSyncType, &deltaRunSync{handler}}
}

// Status sets the status code of the response
func Status(statusCode int) Delta {
	return Delta{statusType, &deltaStatus{statusCode}}
}

// Redirect discards future deltas and redirects to a different URL
func Redirect(location string, code int) Delta {
	return Delta{redirectType, &deltaRedirect{location, code}}
}

// AddHeaders adds some headers to the response
func AddHeaders(headers map[string]string) Delta {
	return Delta{addHeadersType, &deltaAddHeaders{headers}}
}

// SetHeaders sets the headers of the response
func SetHeaders(headers map[string]string) Delta {
	return Delta{setHeadersType, &deltaSetHeaders{headers}}
}

// RmHeaders removes haders from the response
func RmHeaders(headers []string) Delta {
	return Delta{rmHeadersType, &deltaRmHeaders{headers}}
}

// Answer discards future deltas and sends the provided raw response
func Answer(reader io.ReadCloser) Delta {
	return Delta{answerType, &deltaAnswer{reader}}
}
