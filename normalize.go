package wit

// Normalize resolves the provided delta to its normalized representation,
// following jumps locally
func Normalize(delta Delta) Delta {
	_, nextDelta := normalize(&normalizationContext{}, delta)
	return nextDelta
}

type normalizationContext struct{}

func normalize(c *normalizationContext, delta Delta) (nextContext *normalizationContext, nextDelta Delta) {
	nextContext = c
	nextDelta = delta

	switch delta.typeID {

	case sliceType:
		var nd Delta
		deltas := delta.delta.(*deltaSlice).deltas
		nextDeltas := make([]Delta, 0, len(deltas))

		for _, childDelta := range deltas {
			if c != nextContext {
				discardDelta(childDelta)
			} else {
				nextContext, nd = normalize(nextContext, childDelta)
				if nd.typeID != 0 && c == nextContext {
					nextDeltas = append(nextDeltas, nd)
				}
			}
		}

		if c == nextContext {
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
			if c != nextContext {
				discardDelta(childDelta)
			} else {
				nextContext, nd = normalize(nextContext, childDelta)
				if c != nextContext {
					cancel()
				} else if nd.typeID != 0 {
					nextDeltas = append(nextDeltas, nd)
				}
			}
		}

		if c == nextContext {
			nextDelta = List(nextDeltas...)
		} else {
			nextDelta = nd
		}

	case rootType:
		nextContext, nextDelta = normalize(c, delta.delta.(*deltaRoot).delta)
		if c == nextContext {
			nextDelta = Root(nextDelta)
		}

	case selectorType:
		d := delta.delta.(*deltaSelector)
		nextContext, nextDelta = normalize(c, d.delta)
		if c == nextContext {
			nextDelta = d.selector.One(nextDelta)
		}

	case selectorAllType:
		d := delta.delta.(*deltaSelectorAll)
		nextContext, nextDelta = normalize(c, d.delta)
		if c == nextContext {
			nextDelta = d.selector.All(nextDelta)
		}

	case parentType:
		nextContext, nextDelta = normalize(c, delta.delta.(*deltaParent).delta)
		if c == nextContext {
			nextDelta = Parent(nextDelta)
		}

	case firstChildType:
		nextContext, nextDelta = normalize(c, delta.delta.(*deltaFirstChild).delta)
		if c == nextContext {
			nextDelta = FirstChild(nextDelta)
		}

	case lastChildType:
		nextContext, nextDelta = normalize(c, delta.delta.(*deltaLastChild).delta)
		if c == nextContext {
			nextDelta = LastChild(nextDelta)
		}

	case prevSiblingType:
		nextContext, nextDelta = normalize(c, delta.delta.(*deltaPrevSibling).delta)
		if c == nextContext {
			nextDelta = PrevSibling(nextDelta)
		}

	case nextSiblingType:
		nextContext, nextDelta = normalize(c, delta.delta.(*deltaNextSibling).delta)
		if c == nextContext {
			nextDelta = NextSibling(nextDelta)
		}

	case jumpType:
		return normalize(&normalizationContext{}, delta.delta.(*deltaJump).delta)

	case runSyncType:
		f := delta.delta.(*deltaRunSync).handler
		return normalize(c, f())

	}

	return
}
