package wit

import "golang.org/x/net/html"

// Delta represents a page change
type Delta interface {
	Apply(root *html.Node, nodes []*html.Node)
	MarshalJSON() ([]byte, error)
}
