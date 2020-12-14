package wit

import (
	"bytes"

	"golang.org/x/net/html"
)

// HTMLSource represents an HTML source
type HTMLSource interface {
	String() string
	Nodes(context *html.Node) []*html.Node
}

type basicHTMLSource struct {
	reader func() string
}

func (b *basicHTMLSource) String() string {
	return b.reader()
}

func (b *basicHTMLSource) Nodes(ctx *html.Node) []*html.Node {
	nodes, err := html.ParseFragment(bytes.NewReader([]byte(b.reader())), ctx)
	if err != nil {
		return []*html.Node{}
	}

	return nodes
}

// HTMLFromStringFunc builds an HTMLSource from a string function
func HTMLFromStringFunc(reader func() string) HTMLSource {
	return &basicHTMLSource{reader}
}

// HTMLFromString builds an HTMLSource from a string
func HTMLFromString(html string) HTMLSource {
	return &basicHTMLSource{func() string {
		return html
	}}
}
