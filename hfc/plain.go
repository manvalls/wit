package hfc

import (
	"bytes"
	"io"
	"strings"

	"github.com/manvalls/wit"
	"golang.org/x/net/html"
)

type plainFactory struct {
	factory func() io.Reader
}

func (f *plainFactory) HTML() io.Reader {
	return f.factory()
}

func (f *plainFactory) Nodes(ctx *html.Node) []*html.Node {
	nodes, err := html.ParseFragment(f.factory(), ctx)
	if err != nil {
		return []*html.Node{}
	}

	return nodes
}

// FromReaderFactory returns a plain HTMLFactory from a function expected to
// return valid HTML
func FromReaderFactory(factory func() io.Reader) wit.HTMLFactory {
	return &plainFactory{factory}
}

// FromString returns a plain HTMLFactory from a string
func FromString(html string) wit.HTMLFactory {
	return &plainFactory{func() io.Reader {
		return strings.NewReader(html)
	}}
}

// FromBytes returns a plain HTMLFactory from a byte slice
func FromBytes(html []byte) wit.HTMLFactory {
	return &plainFactory{func() io.Reader {
		return bytes.NewReader(html)
	}}
}

// FromHandler returns a plain HTMLFactory from a handler function
func FromHandler(handler func(io.Writer)) wit.HTMLFactory {
	return &plainFactory{func() io.Reader {
		r, w := io.Pipe()

		go func() {
			handler(w)
			w.Close()
		}()

		return r
	}}
}
