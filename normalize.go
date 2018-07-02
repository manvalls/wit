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

	}

	return
}
