package wit

// Normalize resolves the provided delta to its normalized representation
func Normalize(delta Delta) (normalizedDelta Delta, err error) {
	normalizedDelta = delta

	switch delta.typeID {

	case sliceType:
		deltas := delta.delta.(*deltaSlice).deltas
		nextDeltas := make([]Delta, 0, len(deltas))

		for _, childDelta := range deltas {
			if err != nil {
				discardDelta(childDelta)
			} else {
				var nd Delta
				nd, err = Normalize(childDelta)
				if nd.typeID != 0 && err != nil {
					nextDeltas = append(nextDeltas, nd)
				}
			}
		}

		if err != nil {
			normalizedDelta = Nil
			return
		}

		normalizedDelta = List(nextDeltas...)

	case channelType:
		var nd Delta
		d := delta.delta.(*deltaChannel)
		nextDeltas := []Delta{}

		channel := d.channel
		cancel := d.cancel

		for childDelta := range channel {
			if err != nil {
				discardDelta(childDelta)
			} else {
				nd, err = Normalize(childDelta)
				if err != nil {
					cancel()
				} else if nd.typeID != 0 {
					nextDeltas = append(nextDeltas, nd)
				}
			}
		}

		if err != nil {
			normalizedDelta = Nil
			return
		}

		normalizedDelta = List(nextDeltas...)

	case rootType:
		normalizedDelta, err = Normalize(delta.delta.(*deltaRoot).delta)
		if err == nil {
			normalizedDelta = Root(normalizedDelta)
		}

	case selectorType:
		d := delta.delta.(*deltaSelector)
		normalizedDelta, err = Normalize(d.delta)
		if err == nil {
			normalizedDelta = d.selector.One(normalizedDelta)
		}

	case selectorAllType:
		d := delta.delta.(*deltaSelectorAll)
		normalizedDelta, err = Normalize(d.delta)
		if err == nil {
			normalizedDelta = d.selector.All(normalizedDelta)
		}

	case parentType:
		normalizedDelta, err = Normalize(delta.delta.(*deltaParent).delta)
		if err == nil {
			normalizedDelta = Parent(normalizedDelta)
		}

	case firstChildType:
		normalizedDelta, err = Normalize(delta.delta.(*deltaFirstChild).delta)
		if err == nil {
			normalizedDelta = FirstChild(normalizedDelta)
		}

	case lastChildType:
		normalizedDelta, err = Normalize(delta.delta.(*deltaLastChild).delta)
		if err == nil {
			normalizedDelta = LastChild(normalizedDelta)
		}

	case prevSiblingType:
		normalizedDelta, err = Normalize(delta.delta.(*deltaPrevSibling).delta)
		if err == nil {
			normalizedDelta = PrevSibling(normalizedDelta)
		}

	case nextSiblingType:
		normalizedDelta, err = Normalize(delta.delta.(*deltaNextSibling).delta)
		if err == nil {
			normalizedDelta = NextSibling(normalizedDelta)
		}

	case errorType:
		err = delta.delta.(*deltaError).err
		normalizedDelta = Nil

	case runSyncType:
		f := delta.delta.(*deltaRunSync).handler
		return Normalize(f())

	}

	return
}
