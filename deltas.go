package wit

// Delta represents a document change
type Delta struct {
	typeID uint
	info   *deltaInfo
}

// Command encapsulates a delta
type Command interface {
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

// IsNil checks if the given command has nil effect
func IsNil(command Command) bool {
	return command == nil || command.Delta().typeID == 0
}
