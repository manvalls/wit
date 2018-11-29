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
	d := List(actions...).Delta()
	if d.typeID == 0 {
		return d
	}

	return Delta{selectorType, &deltaSelector{s, d}}
}

// All applies the given actions to all matching elements
func (s Selector) All(actions ...Action) Action {
	d := List(actions...).Delta()
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
