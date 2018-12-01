package wit

import (
	"sync"

	"github.com/andybalholm/cascadia"
)

// Selector wraps a CSS selector
type Selector struct {
	*selectorInfo
}

// S wraps a CSS selector in a Selector object
func S(selector string) Selector {
	return Selector{&selectorInfo{selector, sync.Mutex{}, nil}}
}

// Head matches the head element
var Head = S("head")

// Body matches the body element
var Body = S("body")

// One applies the given actions to the first matching element
func (s Selector) One(actions ...Action) Action {
	deltas := extractDeltas(actions)
	if len(deltas) == 0 {
		return Nil
	}

	return Delta{selectorType, &deltaInfo{selector: s, deltas: deltas}}
}

// All applies the given actions to all matching elements
func (s Selector) All(actions ...Action) Action {
	deltas := extractDeltas(actions)
	if len(deltas) == 0 {
		return Nil
	}

	return Delta{selectorAllType, &deltaInfo{selector: s, deltas: deltas}}
}

type selectorInfo struct {
	selectorText string
	mutex        sync.Mutex
	sel          cascadia.Selector
}

func (s *selectorInfo) selector() cascadia.Selector {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.sel != nil {
		return s.sel
	}

	selector, err := cascadia.Compile(s.selectorText)
	if err != nil {
		return nil
	}

	s.sel = selector
	return selector
}
