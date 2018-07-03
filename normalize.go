package wit

import "net/http"

// Normalize resolves the provided delta to its normalized representation. Jumps
// and keys are resolved locally.
func Normalize(delta Delta, keys []string) Delta {
	return delta
}

// Clean resolves the provided delta to its normalized representation,
// removing header and status information and returning it. Jumps and keys are
// resolved locally.
func Clean(delta Delta, keys []string) CleanDelta {
	return CleanDelta{Delta: delta}
}

// CleanDelta holds the result of a Clean operation
type CleanDelta struct {
	Delta        Delta
	Status       int
	HeadersToRm  []string
	HeadersToSet http.Header
	HeadersToAdd http.Header
}

type normalizationRef struct{}

type normalizationContext struct {
	ref      *normalizationRef
	keys     map[string]bool
	baseKeys map[string]bool
	deferred []*deltaWithRef
}

type deltaWithRef struct {
	delta Delta
	ref   *normalizationRef
}

func normalize(c *normalizationContext, delta Delta) (nextContext *normalizationContext, nextDelta Delta) {
	nextContext = c
	nextDelta = delta

	switch delta.typeID {

	case sliceType:
		var nd Delta
		deltas := delta.delta.(*deltaSlice).deltas
		nextDeltas := make([]Delta, 0, len(deltas))

		for _, childDelta := range deltas {
			if c.ref != nextContext.ref {
				discardDelta(childDelta)
			} else {
				nextContext, nd = normalize(nextContext, childDelta)
				if nd.typeID != 0 && c.ref == nextContext.ref {
					nextDeltas = append(nextDeltas, nd)
				}
			}
		}

		if c.ref == nextContext.ref {
			nextDelta = List(nextDeltas...)
		} else {
			nextDelta = nd
		}

	case channelType:
		var nd Delta
		d := delta.delta.(*deltaChannel)
		nextDeltas := []Delta{}

		channel := d.channel
		cancel := d.cancel

		for childDelta := range channel {
			if c.ref != nextContext.ref {
				discardDelta(childDelta)
			} else {
				nextContext, nd = normalize(nextContext, childDelta)
				if c.ref != nextContext.ref {
					cancel()
				} else if nd.typeID != 0 {
					nextDeltas = append(nextDeltas, nd)
				}
			}
		}

		if c.ref == nextContext.ref {
			nextDelta = List(nextDeltas...)
		} else {
			nextDelta = nd
		}

	case rootType:
		nextContext, nextDelta = normalize(c, delta.delta.(*deltaRoot).delta)
		if c.ref == nextContext.ref {
			nextDelta = Root(nextDelta)
		}

	case selectorType:
		d := delta.delta.(*deltaSelector)
		nextContext, nextDelta = normalize(c, d.delta)
		if c.ref == nextContext.ref {
			nextDelta = d.selector.One(nextDelta)
		}

	case selectorAllType:
		d := delta.delta.(*deltaSelectorAll)
		nextContext, nextDelta = normalize(c, d.delta)
		if c.ref == nextContext.ref {
			nextDelta = d.selector.All(nextDelta)
		}

	case parentType:
		nextContext, nextDelta = normalize(c, delta.delta.(*deltaParent).delta)
		if c.ref == nextContext.ref {
			nextDelta = Parent(nextDelta)
		}

	case firstChildType:
		nextContext, nextDelta = normalize(c, delta.delta.(*deltaFirstChild).delta)
		if c.ref == nextContext.ref {
			nextDelta = FirstChild(nextDelta)
		}

	case lastChildType:
		nextContext, nextDelta = normalize(c, delta.delta.(*deltaLastChild).delta)
		if c.ref == nextContext.ref {
			nextDelta = LastChild(nextDelta)
		}

	case prevSiblingType:
		nextContext, nextDelta = normalize(c, delta.delta.(*deltaPrevSibling).delta)
		if c.ref == nextContext.ref {
			nextDelta = PrevSibling(nextDelta)
		}

	case nextSiblingType:
		nextContext, nextDelta = normalize(c, delta.delta.(*deltaNextSibling).delta)
		if c.ref == nextContext.ref {
			nextDelta = NextSibling(nextDelta)
		}

	case jumpType:
		return normalize(&normalizationContext{
			ref:      &normalizationRef{},
			keys:     c.baseKeys,
			baseKeys: c.baseKeys,
		}, delta.delta.(*deltaJump).delta)

	case runSyncType:
		f := delta.delta.(*deltaRunSync).handler
		return normalize(c, f())

	case withKeyType:
		d := delta.delta.(*deltaWithKey)

		if c.keys[d.key] {
			discardDelta(d.delta)
			return
		}

		c.keys[d.key] = true
		nextContext, nextDelta = normalize(c, d.delta)
		if c.ref == nextContext.ref {
			nextDelta = WithKey(d.key, nextDelta)
		}

	case clearKeyType:
		key := delta.delta.(*deltaClearKey).key
		delete(c.keys, key)

	}

	return
}
