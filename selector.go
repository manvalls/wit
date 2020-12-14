package wit

import (
	"sync"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

// Selector wraps a CSS selector
type Selector interface {
	String() string
	cascadia.Matcher
}

type selector struct {
	selector string
	mutex    sync.Mutex
	matcher  cascadia.Matcher
}

func (s *selector) String() string {
	return s.selector
}

func (s *selector) Match(n *html.Node) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.matcher == nil {
		matcher, err := cascadia.ParseGroupWithPseudoElements(s.selector)
		if err != nil {
			return false
		}

		s.matcher = matcher
	}

	return s.matcher.Match(n)
}

// S wraps a CSS selector in a Selector object
func S(s string) Selector {
	return &selector{s, sync.Mutex{}, nil}
}

// Head matches the head element
var Head = S("head")

// Body matches the body element
var Body = S("body")
