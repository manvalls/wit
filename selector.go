package wit

// Selector wraps a CSS selector
type Selector struct {
	*selectorInfo
}

// S wraps a CSS selector in a Selector object
func S(selector string) Selector {
	return Selector{&selectorInfo{selector}}
}

// One applies the given delta to the first matching element
func (s Selector) One(deltas ...Delta) Delta {
	if len(deltas) == 1 {
		return Delta{selectorType, &deltaSelector{s, deltas[0]}}
	}

	return Delta{selectorType, &deltaSelector{s, List(deltas...)}}
}

// All applies the given delta to all matching elements
func (s Selector) All(deltas ...Delta) Delta {
	if len(deltas) == 1 {
		return Delta{selectorAllType, &deltaSelectorAll{s, deltas[0]}}
	}

	return Delta{selectorAllType, &deltaSelectorAll{s, List(deltas...)}}
}

type selectorInfo struct {
	selector string
}
