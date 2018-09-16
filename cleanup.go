package wit

// Discard cleans up the resources associated with the given delta
func Discard(delta Delta) {

	switch delta.typeID {

	case cleanupType:
		delta.delta.(*deltaCleanup).handler()

	case sliceType:
		deltas := delta.delta.(*deltaSlice).deltas
		for _, childDelta := range deltas {
			Discard(childDelta)
		}

	case channelType:
		d := delta.delta.(*deltaChannel)

		d.cancel()
		for childDelta := range d.channel {
			Discard(childDelta)
		}

	case rootType:
		Discard(delta.delta.(*deltaRoot).delta)

	case selectorType:
		Discard(delta.delta.(*deltaSelector).delta)

	case selectorAllType:
		Discard(delta.delta.(*deltaSelectorAll).delta)

	case firstChildType:
		Discard(delta.delta.(*deltaFirstChild).delta)

	case lastChildType:
		Discard(delta.delta.(*deltaLastChild).delta)

	case prevSiblingType:
		Discard(delta.delta.(*deltaPrevSibling).delta)

	case nextSiblingType:
		Discard(delta.delta.(*deltaNextSibling).delta)

	}

}
