package wit

// Delta represents a document change
type Delta struct {
	typeID uint
	delta  interface{}
}

// Action encapsulates a delta
type Action interface {
	Delta() Delta
}

// Delta returns this delta itself
func (d Delta) Delta() Delta {
	return d
}

// - Delta groups

type deltaSlice struct {
	deltas []Delta
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
	factory Factory
}

type deltaText struct {
	text string
}

type deltaReplace struct {
	factory Factory
}

type deltaAppend struct {
	factory Factory
}

type deltaPrepend struct {
	factory Factory
}

type deltaInsertAfter struct {
	factory Factory
}

type deltaInsertBefore struct {
	factory Factory
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

// IsNil checks if the given action has nil effect
func IsNil(action Action) bool {
	return action.Delta().typeID == 0
}
