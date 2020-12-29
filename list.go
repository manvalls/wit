package wit

import (
	"encoding/json"

	"golang.org/x/net/html"
)

// List holds a list of deltas
type List struct {
	Deltas []Delta
}

func unmarshalJSONList(input []interface{}, offset int) List {
	deltas := make([]Delta, 0, len(input)-offset)

	for i := offset; i < len(input); i++ {
		if subjson, ok := input[i].([]interface{}); ok {
			deltas = append(deltas, unmarshalJSONDelta(subjson))
		}
	}

	return List{deltas}
}

func unmarshalDeltaParameter(input []interface{}, offset int) Delta {
	list := unmarshalJSONList(input, offset)

	if len(list.Deltas) == 1 {
		return list.Deltas[0]
	}

	return list
}

func unmarshalStrMap(input interface{}) map[string]string {
	strMap := map[string]string{}
	if jsonMap, ok := input.(map[string]interface{}); ok {
		for key, value := range jsonMap {
			if str, ok := value.(string); ok {
				strMap[key] = str
			}
		}
	}

	return strMap
}

func unmarshalStrArray(input []interface{}, offset int) []string {
	strSlice := []string{}
	for i := offset; i < len(input); i++ {
		if str, ok := input[i].(string); ok {
			strSlice = append(strSlice, str)
		}
	}

	return strSlice
}

func unmarshalJSONDelta(input []interface{}) Delta {

	if len(input) == 0 {
		return nil
	}

	if label, ok := input[0].(float64); ok {
		switch label {
		case listLabel:
			return unmarshalJSONList(input, 1)

		case rootLabel:
			return Root{unmarshalDeltaParameter(input, 1)}

		case selectorLabel:
			if selector, ok := input[1].(string); ok {
				return First{S(selector), unmarshalDeltaParameter(input, 2)}
			}

		case selectorAllLabel:
			if selector, ok := input[1].(string); ok {
				return All{S(selector), unmarshalDeltaParameter(input, 2)}
			}

		case parentLabel:
			return Parent{unmarshalDeltaParameter(input, 1)}

		case firstChildLabel:
			return FirstChild{unmarshalDeltaParameter(input, 1)}

		case lastChildLabel:
			return LastChild{unmarshalDeltaParameter(input, 1)}

		case prevSiblingLabel:
			return PrevSibling{unmarshalDeltaParameter(input, 1)}

		case nextSiblingLabel:
			return NextSibling{unmarshalDeltaParameter(input, 1)}

		case removeLabel:
			return Remove{}

		case clearLabel:
			return Clear{}

		case htmlLabel:
			if html, ok := input[1].(string); ok {
				return HTML{HTMLFromString(html)}
			}

		case replaceLabel:
			if html, ok := input[1].(string); ok {
				return Replace{HTMLFromString(html)}
			}

		case appendLabel:
			if html, ok := input[1].(string); ok {
				return Append{HTMLFromString(html)}
			}

		case prependLabel:
			if html, ok := input[1].(string); ok {
				return Prepend{HTMLFromString(html)}
			}

		case insertAfterLabel:
			if html, ok := input[1].(string); ok {
				return InsertAfter{HTMLFromString(html)}
			}

		case insertBeforeLabel:
			if html, ok := input[1].(string); ok {
				return InsertBefore{HTMLFromString(html)}
			}

		case setAttrLabel:
			return SetAttr{unmarshalStrMap(input[1])}

		case replaceAttrLabel:
			return ReplaceAttr{unmarshalStrMap(input[1])}

		case rmAttrLabel:
			return RmAttr{unmarshalStrArray(input, 1)}

		case setStylesLabel:
			return SetStyles{unmarshalStrMap(input[1])}

		case rmStylesLabel:
			return RmStyles{unmarshalStrArray(input, 1)}

		case addClassesLabel:
			if classes, ok := input[1].(string); ok {
				return AddClasses{classes}
			}

		case rmClassesLabel:
			if classes, ok := input[1].(string); ok {
				return RmClasses{classes}
			}

		}
	}

	return nil
}

// UnmarshalJSON sets *l to the unmarshalled list of deltas
func (l *List) UnmarshalJSON(payload []byte) error {
	var input []interface{}

	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}

	switch delta := unmarshalJSONDelta(input).(type) {
	case List:
		l.Deltas = delta.Deltas
	case nil:
	default:
		l.Deltas = []Delta{delta}
	}

	return nil
}

// Apply applies the delta to the provided elements
func (l List) Apply(root *html.Node, nodes []*html.Node) {
	for _, delta := range l.Deltas {
		delta.Apply(root, nodes)
	}
}

// MarshalJSON marshals the delta to JSON format
func (l List) MarshalJSON() ([]byte, error) {
	return []byte("[" + listLabelJSON + deltasToCSV(l.Deltas) + "]"), nil
}
