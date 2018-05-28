package wit

import (
	"context"
	"io"
	"strings"
)

// List groups a list of deltas together
func List(deltas ...Delta) Delta {
	return Delta{&deltaSlice{deltas}}
}

// Root applies given deltas to the root of the document
func Root(deltas ...Delta) Delta {
	return Delta{&deltaRoot{List(deltas...)}}
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

	return Delta{&deltaChannel{channel, cancel}}
}

// RunChannel runs the given function under the given context and channel, returning a delta
func RunChannel(parentCtx context.Context, callback func(context.Context, chan<- Delta) Delta) Delta {
	ctx, cancel := context.WithCancel(parentCtx)
	channel := make(chan Delta)
	go callback(ctx, channel)
	return Delta{&deltaChannel{channel, cancel}}
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

	return Delta{&deltaHTMLPipe{reader, cancel}}
}

// Remove removes from the document matching elements
var Remove = Delta{&deltaRemove{}}

// HTML sets the inner HTML of the matching elements
func HTML(html string) Delta {
	return Delta{&deltaHTML{strings.NewReader(html)}}
}

// HTMLReader sets the inner HTML of the matching elements
func HTMLReader(reader io.Reader) Delta {
	return Delta{&deltaHTML{reader}}
}

// Text sets the inner text of the matching elements
func Text(txt string) Delta {
	return Delta{&deltaText{txt}}
}

// Replace empties matching elements and applies the provided deltas to them
func Replace(deltas ...Delta) Delta {
	return Delta{&deltaReplace{List(deltas...)}}
}

// Append creates a fragment at the end of the matching elements and
// applies the provided deltas to it
func Append(deltas ...Delta) Delta {
	return Delta{&deltaAppend{List(deltas...)}}
}

// Prepend creates a fragment at the beginning of the matching elements and
// applies the provided deltas to it
func Prepend(deltas ...Delta) Delta {
	return Delta{&deltaPrepend{List(deltas...)}}
}

// InsertAfter creates a fragment after the matching elements and
// applies the provided deltas to it
func InsertAfter(deltas ...Delta) Delta {
	return Delta{&deltaInsertAfter{List(deltas...)}}
}

// InsertBefore creates a fragment before the matching elements and
// applies the provided deltas to it
func InsertBefore(deltas ...Delta) Delta {
	return Delta{&deltaInsertBefore{List(deltas...)}}
}

// AddAttr adds the provided attributes to the matching elements
func AddAttr(attr map[string]string) Delta {
	return Delta{&deltaAddAttr{attr}}
}

// SetAttr sets the attributes of the matching elements
func SetAttr(attr map[string]string) Delta {
	return Delta{&deltaSetAttr{attr}}
}

// RmAttr removes the provided attributes from the matching elements
func RmAttr(attrs ...string) Delta {
	return Delta{&deltaRmAttr{attrs}}
}

// AddStyles adds the provided styles to the matching elements
func AddStyles(styles map[string]string) Delta {
	return Delta{&deltaAddStyles{styles}}
}

// RmStyles removes the provided styles from the matching elements
func RmStyles(styles ...string) Delta {
	return Delta{&deltaRmStyles{styles}}
}

// AddClass adds the provided class to the matching elements
func AddClass(class string) Delta {
	return Delta{&deltaAddClass{class}}
}

// RmClass adds the provided class to the matching elements
func RmClass(class string) Delta {
	return Delta{&deltaRmClass{class}}
}

// ReplaceWith replaces matching elements with that matching the provided selector
func ReplaceWith(selector Selector) Delta {
	return Delta{&deltaReplaceWith{selector}}
}

// ReplaceWithClone replaces matching elements with a clone of that matching the provided selector
func ReplaceWithClone(selector Selector) Delta {
	return Delta{&deltaReplaceWithClone{selector}}
}

// AppendFrom appends the element matching the provided selector to matching elements
func AppendFrom(selector Selector) Delta {
	return Delta{&deltaAppendFrom{selector}}
}

// AppendCloneFrom appends a clone of the element matching the provided selector to matching elements
func AppendCloneFrom(selector Selector) Delta {
	return Delta{&deltaAppendCloneFrom{selector}}
}

// PrependFrom prepends the element matching the provided selector to matching elements
func PrependFrom(selector Selector) Delta {
	return Delta{&deltaPrependFrom{selector}}
}

// PrependCloneFrom prepends a clone of the element matching the provided selector to matching elements
func PrependCloneFrom(selector Selector) Delta {
	return Delta{&deltaPrependCloneFrom{selector}}
}

// InsertAfterFrom inserts the element matching the provided selector after matching elements
func InsertAfterFrom(selector Selector) Delta {
	return Delta{&deltaInsertAfterFrom{selector}}
}

// InsertCloneAfterFrom inserts a clone of the element matching the provided selector
// after matching elements
func InsertCloneAfterFrom(selector Selector) Delta {
	return Delta{&deltaInsertCloneAfterFrom{selector}}
}

// InsertBeforeFrom inserts the element matching the provided selector before matching elements
func InsertBeforeFrom(selector Selector) Delta {
	return Delta{&deltaInsertBeforeFrom{selector}}
}

// InsertCloneBeforeFrom inserts a clone of the element matching the provided selector
// before matching elements
func InsertCloneBeforeFrom(selector Selector) Delta {
	return Delta{&deltaInsertCloneBeforeFrom{selector}}
}

// JS loads the provided script synchronously
func JS(url string) Delta {
	return Delta{&deltaJS{url}}
}

// AsyncJS loads the provided script asynchronously
func AsyncJS(url string) Delta {
	return Delta{&deltaAsyncJS{url}}
}

// CSS loads the provided script synchronously
func CSS(url string) Delta {
	return Delta{&deltaCSS{url}}
}

// AsyncCSS loads the provided script asynchronously
func AsyncCSS(url string) Delta {
	return Delta{&deltaAsyncCSS{url}}
}

// Call calls a JavaScript function with provided parameters, when it becomes available
func Call(path []string, args map[string]string) Delta {
	return Delta{&deltaCall{path, args}}
}

// Jump discards all deltas present and future and applies the given delta to the document
func Jump(delta Delta) Delta {
	return Delta{&deltaJump{delta}}
}

// Redirect discards future deltas and redirects to a different URL
func Redirect(location string, code int) Delta {
	return Delta{&deltaRedirect{location, code}}
}

// AddHeaders adds some headers to the response
func AddHeaders(headers map[string]string) Delta {
	return Delta{&deltaAddHeaders{headers}}
}

// SetHeaders sets the headers of the response
func SetHeaders(headers map[string]string) Delta {
	return Delta{&deltaSetHeaders{headers}}
}

// RmHeaders removes haders from the response
func RmHeaders(headers []string) Delta {
	return Delta{&deltaRmHeaders{headers}}
}

// Answer discards future deltas and sends the provided raw response
func Answer(reader io.ReadCloser) Delta {
	return Delta{&deltaAnswer{reader}}
}
