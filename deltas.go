package wit

import (
	"context"
	"io"
)

// Delta represents a document change
type Delta struct {
	typeID uint
	delta  interface{}
}

// - Delta groups

type deltaSlice struct {
	deltas []Delta
}

type deltaChannel struct {
	channel <-chan Delta
	cancel  context.CancelFunc
}

// - Selectors

type deltaRoot struct {
	delta Delta
}

type deltaSelector struct {
	selector Selector
	delta    Delta
}

type deltaSelectorAll struct {
	selector Selector
	delta    Delta
}

// - Removers

type deltaRemove struct{}

type deltaClear struct{}

// - Content modifiers

type deltaHTML struct {
	reader io.Reader
}

type deltaHTMLPipe struct {
	reader *io.PipeReader
	cancel context.CancelFunc
}

type deltaHTMLFile struct {
	file string
}

type deltaText struct {
	test string
}

type deltaReplace struct {
	delta Delta
}

type deltaAppend struct {
	delta Delta
}

type deltaPrepend struct {
	delta Delta
}

type deltaInsertAfter struct {
	delta Delta
}

type deltaInsertBefore struct {
	delta Delta
}

// - Attributes

type deltaAddAttr struct {
	attr map[string]string
}

type deltaSetAttr struct {
	attr map[string]string
}

type deltaRmAttr struct {
	attr []string
}

type deltaAddStyles struct {
	styles map[string]string
}

type deltaRmStyles struct {
	attr []string
}

type deltaAddClass struct {
	class string
}

type deltaRmClass struct {
	class string
}

// - Loaders

type deltaJS struct {
	url string
}

type deltaAsyncJS struct {
	url string
}

type deltaCSS struct {
	url string
}

type deltaAsyncCSS struct {
	url string
}

// - External scripts

type deltaCall struct {
	path      []string
	arguments map[string]string
}

// - Flow control

type deltaJump struct {
	delta Delta
}

// - Low level request changes

type deltaRedirect struct {
	location string
	code     int
}

type deltaAddHeaders struct {
	headers map[string]string
}

type deltaSetHeaders struct {
	headers map[string]string
}

type deltaRmHeaders struct {
	headers []string
}

type deltaAnswer struct {
	reader io.ReadCloser
}
