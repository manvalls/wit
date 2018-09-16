package wit

// Normalize resolves the provided delta to its normalized representation
func Normalize(delta Delta) (normalizedDelta Delta, err error) {
	cleanupHandlers := []func(){}
	normalizedDelta, err = normalize(delta, &cleanupHandlers)
	if err != nil {
		for _, handler := range cleanupHandlers {
			handler()
		}
	}

	return
}

func normalize(delta Delta, cleanupHandlers *[]func()) (normalizedDelta Delta, err error) {
	normalizedDelta = delta

	switch delta.typeID {

	case sliceType:
		deltas := delta.delta.(*deltaSlice).deltas
		nextDeltas := make([]Delta, 0, len(deltas))

		for _, childDelta := range deltas {
			if err != nil {
				Discard(childDelta)
			} else {
				var nd Delta
				nd, err = normalize(childDelta, cleanupHandlers)
				if nd.typeID != 0 && err == nil {
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
				Discard(childDelta)
			} else {
				nd, err = normalize(childDelta, cleanupHandlers)
				if err != nil {
					cancel()
				} else if nd.typeID != 0 {
					nextDeltas = append(nextDeltas, nd)
				}
			}
		}

		cancel()

		if err != nil {
			normalizedDelta = Nil
			return
		}

		normalizedDelta = List(nextDeltas...)

	case rootType:
		normalizedDelta, err = normalize(delta.delta.(*deltaRoot).delta, cleanupHandlers)
		if err == nil {
			normalizedDelta = Root(normalizedDelta)
		}

	case selectorType:
		d := delta.delta.(*deltaSelector)
		normalizedDelta, err = normalize(d.delta, cleanupHandlers)
		if err == nil {
			normalizedDelta = d.selector.One(normalizedDelta)
		}

	case selectorAllType:
		d := delta.delta.(*deltaSelectorAll)
		normalizedDelta, err = normalize(d.delta, cleanupHandlers)
		if err == nil {
			normalizedDelta = d.selector.All(normalizedDelta)
		}

	case parentType:
		normalizedDelta, err = normalize(delta.delta.(*deltaParent).delta, cleanupHandlers)
		if err == nil {
			normalizedDelta = Parent(normalizedDelta)
		}

	case firstChildType:
		normalizedDelta, err = normalize(delta.delta.(*deltaFirstChild).delta, cleanupHandlers)
		if err == nil {
			normalizedDelta = FirstChild(normalizedDelta)
		}

	case lastChildType:
		normalizedDelta, err = normalize(delta.delta.(*deltaLastChild).delta, cleanupHandlers)
		if err == nil {
			normalizedDelta = LastChild(normalizedDelta)
		}

	case prevSiblingType:
		normalizedDelta, err = normalize(delta.delta.(*deltaPrevSibling).delta, cleanupHandlers)
		if err == nil {
			normalizedDelta = PrevSibling(normalizedDelta)
		}

	case nextSiblingType:
		normalizedDelta, err = normalize(delta.delta.(*deltaNextSibling).delta, cleanupHandlers)
		if err == nil {
			normalizedDelta = NextSibling(normalizedDelta)
		}

	case errorType:
		err = delta.delta.(*deltaError).err
		normalizedDelta = Nil

	case runSyncType:
		f := delta.delta.(*deltaRunSync).handler
		return normalize(f(), cleanupHandlers)

	case cleanupType:
		*cleanupHandlers = append(*cleanupHandlers, delta.delta.(*deltaCleanup).handler)

	}

	return
}
