package wit

import (
	"context"
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

type deltaParent struct {
	delta Delta
}

type deltaFirstChild struct {
	delta Delta
}

type deltaLastChild struct {
	delta Delta
}

type deltaPrevSibling struct {
	delta Delta
}

type deltaNextSibling struct {
	delta Delta
}

// - Removers

type deltaRemove struct{}

type deltaClear struct{}

// - Content modifiers

type deltaHTML struct {
	factory HTMLFactory
}

type deltaText struct {
	text string
}

type deltaReplace struct {
	factory HTMLFactory
}

type deltaAppend struct {
	factory HTMLFactory
}

type deltaPrepend struct {
	factory HTMLFactory
}

type deltaInsertAfter struct {
	factory HTMLFactory
}

type deltaInsertBefore struct {
	factory HTMLFactory
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
	styles []string
}

type deltaAddClass struct {
	class string
}

type deltaRmClass struct {
	class string
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

type deltaRunSync struct {
	handler func() Delta
}

// IsNil checks if the given delta has nil effect
func IsNil(d Delta) bool {
	return d.typeID == 0
}
