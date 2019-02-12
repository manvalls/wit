package wit

import (
	"io"
	"io/ioutil"
	"strconv"
)

type jsonRenderer struct {
	command Command
}

// NewJSONRenderer returns a new renderer which will render JSON
func NewJSONRenderer(command Command) Renderer {
	return &jsonRenderer{command}
}

func (r *jsonRenderer) Render(w io.Writer) error {
	if IsNil(r.command) {
		return writeDeltaJSON(w, Nil.Delta())
	}

	return writeDeltaJSON(w, r.command.Delta())
}

func writeList(w io.Writer, deltas []Delta) (err error) {
	for _, childDelta := range deltas {
		_, err = w.Write([]byte{','})
		if err != nil {
			return
		}

		err = writeDeltaJSON(w, childDelta)
		if err != nil {
			return
		}
	}

	return
}

func writeDeltaJSON(w io.Writer, delta Delta) (err error) {

	_, err = w.Write([]byte{'['})
	if err != nil {
		return
	}

	switch delta.typeID {

	case sliceType:
		_, err = w.Write(sliceTypeString)
		if err != nil {
			return
		}

		err = writeList(w, delta.info.deltas)
		if err != nil {
			return
		}

	case rootType:
		_, err = w.Write(rootTypeString)
		if err != nil {
			return
		}

		err = writeList(w, delta.info.deltas)
		if err != nil {
			return
		}

	case selectorType:
		_, err = w.Write(append(selectorTypeString, ','))
		if err != nil {
			return
		}

		_, err = w.Write([]byte(strconv.Quote(delta.info.selector.selectorText)))
		if err != nil {
			return
		}

		err = writeList(w, delta.info.deltas)
		if err != nil {
			return
		}

	case selectorAllType:
		_, err = w.Write(append(selectorAllTypeString, ','))
		if err != nil {
			return
		}

		_, err = w.Write([]byte(strconv.Quote(delta.info.selector.selectorText)))
		if err != nil {
			return
		}

		err = writeList(w, delta.info.deltas)
		if err != nil {
			return
		}

	case parentType:
		_, err = w.Write(parentTypeString)
		if err != nil {
			return
		}

		err = writeList(w, delta.info.deltas)
		if err != nil {
			return
		}

	case firstChildType:
		_, err = w.Write(firstChildTypeString)
		if err != nil {
			return
		}

		err = writeList(w, delta.info.deltas)
		if err != nil {
			return
		}

	case lastChildType:
		_, err = w.Write(lastChildTypeString)
		if err != nil {
			return
		}

		err = writeList(w, delta.info.deltas)
		if err != nil {
			return
		}

	case prevSiblingType:
		_, err = w.Write(prevSiblingTypeString)
		if err != nil {
			return
		}

		err = writeList(w, delta.info.deltas)
		if err != nil {
			return
		}

	case nextSiblingType:
		_, err = w.Write(nextSiblingTypeString)
		if err != nil {
			return
		}

		err = writeList(w, delta.info.deltas)
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
		result, err = ioutil.ReadAll(delta.info.factory.HTML())
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
		result, err = ioutil.ReadAll(delta.info.factory.HTML())
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
		result, err = ioutil.ReadAll(delta.info.factory.HTML())
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
		result, err = ioutil.ReadAll(delta.info.factory.HTML())
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
		result, err = ioutil.ReadAll(delta.info.factory.HTML())
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
		result, err = ioutil.ReadAll(delta.info.factory.HTML())
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

		_, err = w.Write([]byte(strMapToJSON(delta.info.strMap)))
		if err != nil {
			return
		}

	case setAttrType:
		_, err = w.Write(append(setAttrTypeString, ','))
		if err != nil {
			return
		}

		_, err = w.Write([]byte(strMapToJSON(delta.info.strMap)))
		if err != nil {
			return
		}

	case rmAttrType:
		_, err = w.Write(append(rmAttrTypeString, ','))
		if err != nil {
			return
		}

		_, err = w.Write([]byte(strSliceToCSV(delta.info.strList)))
		if err != nil {
			return
		}

	case addStylesType:
		_, err = w.Write(append(addStylesTypeString, ','))
		if err != nil {
			return
		}

		_, err = w.Write([]byte(strMapToJSON(delta.info.strMap)))
		if err != nil {
			return
		}

	case rmStylesType:
		_, err = w.Write(append(rmStylesTypeString, ','))
		if err != nil {
			return
		}

		_, err = w.Write([]byte(strSliceToCSV(delta.info.strList)))
		if err != nil {
			return
		}

	case addClassType:
		_, err = w.Write(append(addClassTypeString, ','))
		if err != nil {
			return
		}

		_, err = w.Write([]byte(strconv.Quote(delta.info.class)))
		if err != nil {
			return
		}

	case rmClassType:
		_, err = w.Write(append(rmClassTypeString, ','))
		if err != nil {
			return
		}

		_, err = w.Write([]byte(strconv.Quote(delta.info.class)))
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
