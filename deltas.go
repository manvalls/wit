package wit

import (
	"context"
	"io"
)

// Delta represents a document change
type Delta struct {
	delta interface{}
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

type deltaSelector struct {
	selector Selector
	delta    Delta
}

type deltaSelectorAll struct {
	selector Selector
	delta    Delta
}

// - Removers

type removeDelta struct{}

// - Content modifiers

type deltaHTML struct {
	reader io.Reader
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

type deltaToggleClass struct {
	class string
}

// - Placement

type deltaReplaceWith struct {
	selector Selector
}

type deltaReplaceWithClone struct {
	selector Selector
}

type deltaAppendFrom struct {
	selector Selector
}

type deltaAppendCloneFrom struct {
	selector Selector
}

type deltaPrependFrom struct {
	selector Selector
}

type deltaPrependCloneFrom struct {
	selector Selector
}

type deltaInsertAfterFrom struct {
	selector Selector
}

type deltaInsertCloneAfterFrom struct {
	selector Selector
}

type deltaInsertBeforeFrom struct {
	selector Selector
}

type deltaInsertCloneBeforeFrom struct {
	selector Selector
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

// Flow control

type deltaJump struct {
	delta Delta
}

// Low level request changes

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
	reader io.Reader
}
