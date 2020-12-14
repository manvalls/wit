package wit

import "golang.org/x/net/html"

// Delta represents a page change
type Delta interface {
	Empty() bool
	Flatten() Delta
	Apply(root *html.Node, nodes []*html.Node)
	MarshalJSON() ([]byte, error)
}
