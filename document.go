package wit

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Document struct {
	root  *html.Node
	nodes []*html.Node
}

var baseNode, _ = html.Parse(strings.NewReader("<!DOCTYPE html><html><head></head><body></body></html>"))

func NewDocument() Document {
	root := cloneNode(baseNode, map[*html.Node]*html.Node{})
	return Document{
		root,
		[]*html.Node{root},
	}
}

func (d Document) Render(w io.Writer) {
	html.Render(w, d.root)
}
