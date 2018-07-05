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
		ref:        &normalizationRef{},
		cleanDelta: &CleanDelta{},
	}, delta)

	c.cleanDelta.Delta = nextDelta
	return *c.cleanDelta
}

// CleanDelta holds the result of a Clean operation
type CleanDelta struct {
	Delta   Delta
	Status  int
	Headers http.Header
	Answer  io.ReadCloser
}

type normalizationRef struct{}

type normalizationContext struct {
	ref        *normalizationRef
	cleanDelta *CleanDelta
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
		if c.cleanDelta != nil {
			return normalize(&normalizationContext{
				ref:        &normalizationRef{},
				cleanDelta: &CleanDelta{},
			}, delta.delta.(*deltaJump).delta)
		}

		return normalize(&normalizationContext{
			ref: &normalizationRef{},
		}, delta.delta.(*deltaJump).delta)

	case runSyncType:
		f := delta.delta.(*deltaRunSync).handler
		return normalize(c, f())

	case statusType:
		if c.cleanDelta != nil {
			nextDelta = Nil
			c.cleanDelta.Status = delta.delta.(*deltaStatus).code
		}

	case addHeadersType:
		if c.cleanDelta != nil {
			nextDelta = Nil
			headers := delta.delta.(*deltaAddHeaders).headers
			for key, value := range headers {
				for _, h := range value {
					c.cleanDelta.Headers.Add(key, h)
				}
			}
		}

	case setHeadersType:
		if c.cleanDelta != nil {
			nextDelta = Nil
			headers := delta.delta.(*deltaSetHeaders).headers
			for key, value := range headers {
				c.cleanDelta.Headers[key] = value
			}
		}

	case rmHeadersType:
		if c.cleanDelta != nil {
			nextDelta = Nil
			headers := delta.delta.(*deltaRmHeaders).headers
			for _, header := range headers {
				c.cleanDelta.Headers.Del(header)
			}
		}

	case answerType:
		if c.cleanDelta != nil {
			nextDelta = Nil
			c.cleanDelta.Answer = delta.delta.(*deltaAnswer).reader
			c.ref = &normalizationRef{}
		}

	}

	return
}
