package wit

func discardDelta(delta Delta) {

	switch delta.typeID {

	case sliceType:
		deltas := delta.delta.(*deltaSlice).deltas
		for _, childDelta := range deltas {
			discardDelta(childDelta)
		}

	case channelType:
		d := delta.delta.(*deltaChannel)

		d.cancel()
		for childDelta := range d.channel {
			discardDelta(childDelta)
		}

	case rootType:
		discardDelta(delta.delta.(*deltaRoot).delta)

	case selectorType:
		discardDelta(delta.delta.(*deltaSelector).delta)

	case selectorAllType:
		discardDelta(delta.delta.(*deltaSelectorAll).delta)

	case firstChildType:
		discardDelta(delta.delta.(*deltaFirstChild).delta)

	case lastChildType:
		discardDelta(delta.delta.(*deltaLastChild).delta)

	case prevSiblingType:
		discardDelta(delta.delta.(*deltaPrevSibling).delta)

	case nextSiblingType:
		discardDelta(delta.delta.(*deltaNextSibling).delta)

	case jumpType:
		discardDelta(delta.delta.(*deltaJump).delta)

	case withKeyType:
		discardDelta(delta.delta.(*deltaWithKey).delta)

	case deferType:
		discardDelta(delta.delta.(*deltaDefer).delta)

	}

}
