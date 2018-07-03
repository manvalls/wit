package wit

import (
	"io"
	"net/http"
)

// Normalize resolves the provided delta to its normalized representation. Jumps
// are resolved locally.
func Normalize(delta Delta) Delta {
	return delta
}

// Clean resolves the provided delta to its normalized representation,
// removing header and status information and returning it. Jumps are
// resolved locally.
func Clean(delta Delta) CleanDelta {
	c, nextDelta := normalize(&normalizationContext{
		ref: &normalizationRef{},
	}, delta)

	for len(c.deferred) > 0 {
		var nd Delta
		deferred := c.deferred
		c.deferred = nil

		for _, def := range deferred {
			if def.ref != c.ref {
				discardDelta(def.delta)
			} else {
				d := delta.delta.(*deltaDefer)
				c, nd = normalize(c, d.delta)
				if def.ref == c.ref {
					d.delta = nd
				} else {
					nextDelta = nd
				}
			}
		}
	}

	return CleanDelta{Delta: nextDelta}
}

// CleanDelta holds the result of a Clean operation
type CleanDelta struct {
	Delta        Delta
	Status       int
	HeadersToRm  []string
	HeadersToSet http.Header
	HeadersToAdd http.Header
	Answer       io.ReadCloser
}

type normalizationRef struct{}

type normalizationContext struct {
	ref      *normalizationRef
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
			ref: &normalizationRef{},
		}, delta.delta.(*deltaJump).delta)

	case runSyncType:
		f := delta.delta.(*deltaRunSync).handler
		return normalize(c, f())

	case deferType:
		c.deferred = append(c.deferred, &deltaWithRef{
			ref:   c.ref,
			delta: Defer(delta.delta.(*deltaDefer).delta),
		})

	}

	return
}
