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

// One applies the given delta to the first matching element
func (s Selector) One(deltas ...Delta) Delta {
	d := List(deltas...)
	if d.typeID == 0 {
		return d
	}

	return Delta{selectorType, &deltaSelector{s, d}}
}

// All applies the given delta to all matching elements
func (s Selector) All(deltas ...Delta) Delta {
	d := List(deltas...)
	if d.typeID == 0 {
		return d
	}

	return Delta{selectorAllType, &deltaSelectorAll{s, d}}
}

type selectorInfo struct {
	selectorText string
	sync.Mutex
	cascadia.Selector
}

func (s *selectorInfo) selector() cascadia.Selector {
	s.Lock()
	defer s.Unlock()

	if s.Selector != nil {
		return s.Selector
	}

	selector, err := cascadia.Compile(s.selectorText)
	if err != nil {
		return nil
	}

	s.Selector = selector
	return selector
}
