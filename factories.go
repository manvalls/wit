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

	go func() {
		callback(ctx, channel)
	}()

	return Delta{&deltaChannel{channel, cancel}}
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
