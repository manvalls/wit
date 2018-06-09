package html

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/manvalls/wit"
	"github.com/manvalls/wit/util"
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
	return util.Clone(f.nodes)
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
