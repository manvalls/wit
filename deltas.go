package wit

// Delta represents a document change
type Delta struct {
	typeID uint
	info   *deltaInfo
}

// Action encapsulates a delta
type Action interface {
	Delta() Delta
}

// Delta returns this delta itself
func (d Delta) Delta() Delta {
	return d
}

type deltaInfo struct {
	deltas   []Delta
	selector Selector
	factory  Factory
	strMap   map[string]string
	strList  []string
	class    string
}

// IsNil checks if the given action has nil effect
func IsNil(action Action) bool {
	return action == nil || action.Delta().typeID == 0
}
