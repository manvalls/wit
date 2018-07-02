package wit

// Normalize resolves the provided delta to its normalized representation
func Normalize(delta Delta) Delta {
	return delta
}

type normalizationContext struct {
	ref      *Delta
	deferred []*deltaWithRef
}

type deltaWithRef struct {
	delta Delta
	ref   *Delta
}

func normalize(c *normalizationContext, delta Delta) (nextContext *normalizationContext, nextDelta Delta) {
	nextContext = c
	nextDelta = delta

	switch delta.typeID {

	case sliceType:
		var nd Delta
		deltas := delta.delta.(*deltaSlice).deltas
		nextDeltas := make([]Delta, len(deltas))

		for i, childDelta := range deltas {
			if c.ref != nextContext.ref {
				discardDelta(childDelta)
			} else {
				nextContext, nd = normalize(nextContext, childDelta)
				nextDeltas[i] = nd
			}
		}

		if c.ref == nextContext.ref {
			nextDelta = List(nextDeltas...)
		} else {
			nextDelta = nd
		}

	}

	return
}
