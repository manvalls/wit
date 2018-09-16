package wit

import (
	"io"
	"io/ioutil"
	"strconv"
)

type jsonRenderer struct {
	delta Delta
}

// NewJSONRenderer returns a new renderer which will render JSON
func NewJSONRenderer(delta Delta) Renderer {
	return &jsonRenderer{delta}
}

func (r *jsonRenderer) Render(w io.Writer) error {
	return writeDeltaJSON(w, r.delta)
}

func writeListOrDelta(w io.Writer, delta Delta) (err error) {
	if delta.typeID == sliceType {
		for i, childDelta := range delta.delta.(*deltaSlice).deltas {
			if i != 0 {
				_, err = w.Write([]byte{','})
				if err != nil {
					return
				}
			}

			err = writeDeltaJSON(w, childDelta)
			if err != nil {
				return
			}
		}

		return
	}

	return writeDeltaJSON(w, delta)
}

func writeDeltaJSON(w io.Writer, delta Delta) (err error) {

	_, err = w.Write([]byte{'['})
	if err != nil {
		return
	}

	switch delta.typeID {

	case sliceType:
		_, err = w.Write(append(sliceTypeString))
		if err != nil {
			return
		}

		for _, childDelta := range delta.delta.(*deltaSlice).deltas {
			_, err = w.Write([]byte{','})
			if err != nil {
				return
			}

			err = writeDeltaJSON(w, childDelta)
			if err != nil {
				return
			}
		}

	case rootType:
		_, err = w.Write(append(sliceTypeString, ','))
		if err != nil {
			return
		}

		err = writeListOrDelta(w, delta.delta.(*deltaRoot).delta)
		if err != nil {
			return
		}

	case selectorType:
		d := delta.delta.(*deltaSelector)
		_, err = w.Write(append(selectorTypeString, ','))
		if err != nil {
			return
		}

		_, err = w.Write(append([]byte(strconv.Quote(d.selector.selectorText)), ','))
		if err != nil {
			return
		}

		err = writeListOrDelta(w, d.delta)
		if err != nil {
			return
		}

	case selectorAllType:
		d := delta.delta.(*deltaSelectorAll)
		_, err = w.Write(append(selectorAllTypeString, ','))
		if err != nil {
			return
		}

		_, err = w.Write(append([]byte(strconv.Quote(d.selector.selectorText)), ','))
		if err != nil {
			return
		}

		err = writeListOrDelta(w, d.delta)
		if err != nil {
			return
		}

	case parentType:
		_, err = w.Write(append(parentTypeString, ','))
		if err != nil {
			return
		}

		err = writeListOrDelta(w, delta.delta.(*deltaParent).delta)
		if err != nil {
			return
		}

	case firstChildType:
		_, err = w.Write(append(firstChildTypeString, ','))
		if err != nil {
			return
		}

		err = writeListOrDelta(w, delta.delta.(*deltaFirstChild).delta)
		if err != nil {
			return
		}

	case lastChildType:
		_, err = w.Write(append(lastChildTypeString, ','))
		if err != nil {
			return
		}

		err = writeListOrDelta(w, delta.delta.(*deltaLastChild).delta)
		if err != nil {
			return
		}

	case prevSiblingType:
		_, err = w.Write(append(prevSiblingTypeString, ','))
		if err != nil {
			return
		}

		err = writeListOrDelta(w, delta.delta.(*deltaPrevSibling).delta)
		if err != nil {
			return
		}

	case nextSiblingType:
		_, err = w.Write(append(nextSiblingTypeString, ','))
		if err != nil {
			return
		}

		err = writeListOrDelta(w, delta.delta.(*deltaNextSibling).delta)
		if err != nil {
			return
		}

	case removeType:
		_, err = w.Write(removeTypeString)
		if err != nil {
			return
		}

	case clearType:
		_, err = w.Write(clearTypeString)
		if err != nil {
			return
		}

	case htmlType:
		_, err = w.Write(append(htmlTypeString, ','))
		if err != nil {
			return
		}

		var result []byte
		result, err = ioutil.ReadAll(delta.delta.(*deltaHTML).factory.HTML())
		if err != nil {
			return
		}

		_, err = w.Write([]byte(strconv.Quote(string(result))))
		if err != nil {
			return
		}

	case replaceType:
		_, err = w.Write(append(replaceTypeString, ','))
		if err != nil {
			return
		}

		var result []byte
		result, err = ioutil.ReadAll(delta.delta.(*deltaReplace).factory.HTML())
		if err != nil {
			return
		}

		_, err = w.Write([]byte(strconv.Quote(string(result))))
		if err != nil {
			return
		}

	case appendType:
		_, err = w.Write(append(appendTypeString, ','))
		if err != nil {
			return
		}

		var result []byte
		result, err = ioutil.ReadAll(delta.delta.(*deltaAppend).factory.HTML())
		if err != nil {
			return
		}

		_, err = w.Write([]byte(strconv.Quote(string(result))))
		if err != nil {
			return
		}

	case prependType:
		_, err = w.Write(append(prependTypeString, ','))
		if err != nil {
			return
		}

		var result []byte
		result, err = ioutil.ReadAll(delta.delta.(*deltaPrepend).factory.HTML())
		if err != nil {
			return
		}

		_, err = w.Write([]byte(strconv.Quote(string(result))))
		if err != nil {
			return
		}

	case insertAfterType:
		_, err = w.Write(append(insertAfterTypeString, ','))
		if err != nil {
			return
		}

		var result []byte
		result, err = ioutil.ReadAll(delta.delta.(*deltaInsertAfter).factory.HTML())
		if err != nil {
			return
		}

		_, err = w.Write([]byte(strconv.Quote(string(result))))
		if err != nil {
			return
		}

	case insertBeforeType:
		_, err = w.Write(append(insertBeforeTypeString, ','))
		if err != nil {
			return
		}

		var result []byte
		result, err = ioutil.ReadAll(delta.delta.(*deltaInsertBefore).factory.HTML())
		if err != nil {
			return
		}

		_, err = w.Write([]byte(strconv.Quote(string(result))))
		if err != nil {
			return
		}

	case addAttrType:
		_, err = w.Write(append(addAttrTypeString, ','))
		if err != nil {
			return
		}

		_, err = w.Write([]byte(strMapToJSON(delta.delta.(*deltaAddAttr).attr)))
		if err != nil {
			return
		}

	case setAttrType:
		_, err = w.Write(append(setAttrTypeString, ','))
		if err != nil {
			return
		}

		_, err = w.Write([]byte(strMapToJSON(delta.delta.(*deltaSetAttr).attr)))
		if err != nil {
			return
		}

	case rmAttrType:
		_, err = w.Write(append(rmAttrTypeString, ','))
		if err != nil {
			return
		}

		_, err = w.Write([]byte(strSliceToCSV(delta.delta.(*deltaRmAttr).attr)))
		if err != nil {
			return
		}

	case addStylesType:
		_, err = w.Write(append(addStylesTypeString, ','))
		if err != nil {
			return
		}

		_, err = w.Write([]byte(strMapToJSON(delta.delta.(*deltaAddStyles).styles)))
		if err != nil {
			return
		}

	case rmStylesType:
		_, err = w.Write(append(rmStylesTypeString, ','))
		if err != nil {
			return
		}

		_, err = w.Write([]byte(strSliceToCSV(delta.delta.(*deltaRmStyles).styles)))
		if err != nil {
			return
		}

	case addClassType:
		_, err = w.Write(append(addClassTypeString, ','))
		if err != nil {
			return
		}

		_, err = w.Write([]byte(strconv.Quote(delta.delta.(*deltaAddClass).class)))
		if err != nil {
			return
		}

	case rmClassType:
		_, err = w.Write(append(rmClassTypeString, ','))
		if err != nil {
			return
		}

		_, err = w.Write([]byte(strconv.Quote(delta.delta.(*deltaRmClass).class)))
		if err != nil {
			return
		}

	default:
		_, err = w.Write([]byte{'0'})
		if err != nil {
			return
		}

	}

	_, err = w.Write([]byte{']'})
	if err != nil {
		return
	}

	return
}
