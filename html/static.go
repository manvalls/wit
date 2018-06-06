package html

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/manvalls/wit"
	"golang.org/x/net/html"
)

type staticFactory struct {
	bytes []byte
	nodes []*html.Node
}

func (f *staticFactory) HTML() io.Reader {
	return bytes.NewReader(f.bytes)
}

func (f *staticFactory) Nodes() []*html.Node {
	clones := make([]*html.Node, len(f.nodes))
	cache := map[*html.Node]*html.Node{}

	for i, node := range f.nodes {
		clones[i] = cloneNode(node, cache)
	}

	return clones
}

func cloneNode(node *html.Node, cache map[*html.Node]*html.Node) *html.Node {
	if node == nil {
		return nil
	}

	if val, ok := cache[node]; ok {
		return val
	}

	newNode := &html.Node{}
	cache[node] = newNode

	newNode.Parent = cloneNode(node.Parent, cache)
	newNode.FirstChild = cloneNode(node.FirstChild, cache)
	newNode.LastChild = cloneNode(node.LastChild, cache)
	newNode.PrevSibling = cloneNode(node.PrevSibling, cache)
	newNode.NextSibling = cloneNode(node.NextSibling, cache)

	newNode.Type = node.Type
	newNode.DataAtom = node.DataAtom
	newNode.Data = node.Data
	newNode.Namespace = node.Namespace

	newNode.Attr = make([]html.Attribute, len(node.Attr))
	copy(newNode.Attr, node.Attr)

	return newNode
}

// Static calls all methods from the provided factory, stores its results in
// memory and returns a new factory which will return copies of the original
// values
func Static(factory wit.HTMLFactory) wit.HTMLFactory {
	bytes := make(chan []byte)
	nodes := make(chan []*html.Node)

	go func() {
		b, err := ioutil.ReadAll(factory.HTML())
		if err != nil {
			b = []byte{}
		}

		bytes <- b
	}()

	go func() {
		nodes <- factory.Nodes()
	}()

	return &staticFactory{<-bytes, <-nodes}
}
