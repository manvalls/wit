package wit

import (
	"context"
	"io"
	"strings"
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

// RunHTML runs the given function under the given context and writer, returning a delta
func RunHTML(parentCtx context.Context, callback func(*io.PipeWriter)) Delta {
	reader, writer := io.Pipe()
	ctx, cancel := context.WithCancel(parentCtx)

	go callback(writer)
	go func() {
		<-ctx.Done()
		reader.Close()
	}()

	return Delta{htmlPipeType, &deltaHTMLPipe{reader, cancel}}
}

// Nil represents an effectless delta
var Nil = Delta{}

// Remove removes from the document matching elements
var Remove = Delta{removeType, &deltaRemove{}}

// Clear empties matching elements
var Clear = Delta{clearType, &deltaClear{}}

// HTML sets the inner HTML of the matching elements
func HTML(html string) Delta {
	return Delta{htmlType, &deltaHTML{strings.NewReader(html)}}
}

// HTMLReader sets the inner HTML of the matching elements
func HTMLReader(reader io.Reader) Delta {
	return Delta{htmlType, &deltaHTML{reader}}
}

// HTMLFile sets the inner HTML of the matching elements
func HTMLFile(file string) Delta {
	return Delta{htmlFileType, &deltaHTMLFile{file}}
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

// Replace replaces matching elements with empty fragments and applies the
// provided deltas to them
func Replace(deltas ...Delta) Delta {
	return Delta{replaceType, &deltaReplace{List(deltas...)}}
}

// Append creates a fragment at the end of the matching elements and
// applies the provided deltas to it
func Append(deltas ...Delta) Delta {
	return Delta{appendType, &deltaAppend{List(deltas...)}}
}

// Prepend creates a fragment at the beginning of the matching elements and
// applies the provided deltas to it
func Prepend(deltas ...Delta) Delta {
	return Delta{prependType, &deltaPrepend{List(deltas...)}}
}

// InsertAfter creates a fragment after the matching elements and
// applies the provided deltas to it
func InsertAfter(deltas ...Delta) Delta {
	return Delta{insertAfterType, &deltaInsertAfter{List(deltas...)}}
}

// InsertBefore creates a fragment before the matching elements and
// applies the provided deltas to it
func InsertBefore(deltas ...Delta) Delta {
	return Delta{insertBeforeType, &deltaInsertBefore{List(deltas...)}}
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
func JS(url string) Delta {
	return Delta{jsType, &deltaJS{url}}
}

// AsyncJS loads the provided script asynchronously
func AsyncJS(url string) Delta {
	return Delta{asyncJSType, &deltaAsyncJS{url}}
}

// CSS loads the provided script synchronously
func CSS(url string) Delta {
	return Delta{cssType, &deltaCSS{url}}
}

// AsyncCSS loads the provided script asynchronously
func AsyncCSS(url string) Delta {
	return Delta{asyncCSSType, &deltaAsyncCSS{url}}
}

// Call calls a JavaScript function with provided parameters, when it becomes available
func Call(path []string, args map[string]string) Delta {
	return Delta{callType, &deltaCall{path, args}}
}

// Jump discards all deltas present and future and applies the given delta to the document
func Jump(delta Delta) Delta {
	return Delta{jumpType, &deltaJump{delta}}
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
