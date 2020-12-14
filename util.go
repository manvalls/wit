package wit

import (
	"strconv"

	"golang.org/x/net/html"
)

func clone(nodes []*html.Node) []*html.Node {
	clones := make([]*html.Node, len(nodes))
	cache := map[*html.Node]*html.Node{}

	for i, node := range nodes {
		clones[i] = cloneNode(node, cache)
	}

	return clones
}

func cloneNode(node *html.Node, cache map[*html.Node]*html.Node) *html.Node {
	if node == nil {
		return nil
	}

	if val, ok := cache[node]; ok {
		return val
	}

	newNode := &html.Node{}
	cache[node] = newNode

	newNode.Parent = cloneNode(node.Parent, cache)
	newNode.FirstChild = cloneNode(node.FirstChild, cache)
	newNode.LastChild = cloneNode(node.LastChild, cache)
	newNode.PrevSibling = cloneNode(node.PrevSibling, cache)
	newNode.NextSibling = cloneNode(node.NextSibling, cache)

	newNode.Type = node.Type
	newNode.DataAtom = node.DataAtom
	newNode.Data = node.Data
	newNode.Namespace = node.Namespace

	newNode.Attr = make([]html.Attribute, len(node.Attr))
	copy(newNode.Attr, node.Attr)

	return newNode
}

func parseStyle(style string) map[string]string {
	styleMap := map[string]string{}

	key := ""
	value := ""
	lookingForKey := true
	lookingForValue := false
	fillingKey := false
	fillingValue := false
	inSingleQuoteString := false
	inDoubleQuoteString := false
	escapedChar := false

	for _, r := range style {
	start:
		if fillingValue {
			if escapedChar {
				escapedChar = false
				value += string(r)
				continue
			}

			switch r {
			case '\\':
				escapedChar = true
				value += string(r)
				continue
			case '\'':
				if inSingleQuoteString {
					inSingleQuoteString = false
				} else {
					inSingleQuoteString = true
				}

				value += string(r)
				continue
			case '"':
				if inDoubleQuoteString {
					inDoubleQuoteString = false
				} else {
					inDoubleQuoteString = true
				}

				value += string(r)
				continue
			default:
				if inSingleQuoteString || inDoubleQuoteString {
					value += string(r)
					continue
				}
			}

			switch r {
			case ';':
				if key != "" {
					styleMap[key] = value
				}

				key = ""
				value = ""

				lookingForKey = true
				lookingForValue = false
				fillingKey = false
				fillingValue = false
			default:
				value += string(r)
			}

		} else {
			switch r {
			case ' ', '\t', '\r', '\n', '\f':
				fillingKey = false
			case ':':
				value = ""
				lookingForKey = false
				lookingForValue = true
				fillingKey = false
				fillingValue = false
			default:
				if lookingForKey || fillingKey {
					fillingKey = true
					lookingForKey = false
					key += string(r)
				} else if lookingForValue || fillingValue {
					fillingValue = true
					lookingForValue = false
					goto start
				}
			}
		}

	}

	if fillingValue && key != "" {
		styleMap[key] = value
	}

	return styleMap
}

func buildStyle(style map[string]string) string {
	attr := ""
	for key, value := range style {
		attr += key + ": " + value + ";"
	}

	return attr
}

func parseClass(class string) map[string]bool {
	currentClass := ""
	classes := map[string]bool{}

	flush := func() {
		if currentClass != "" {
			classes[currentClass] = true
			currentClass = ""
		}
	}

	for _, r := range class {
		switch r {
		case ' ', '\t', '\r', '\n', '\f':
			flush()
		default:
			currentClass += string(r)
		}
	}

	flush()
	return classes
}

func buildClass(classes map[string]bool) string {
	class := ""
	i := 0

	for key, value := range classes {
		if i != 0 {
			class += " "
		}

		if value {
			class += key
		}

		i++
	}

	return class
}

func strMapToJSON(args map[string]string) string {
	result := "{"
	i := 0

	for key, value := range args {
		if i != 0 {
			result += ","
		}

		result += strconv.Quote(key) + ":" + strconv.Quote(value)
		i++
	}

	result += "}"
	return result
}

func deltasToCSV(deltas []Delta) string {
	result := ""

	for _, delta := range deltas {
		deltaJSON, err := delta.MarshalJSON()
		if err != nil {
			return ""
		}

		result += "," + string(deltaJSON)
	}

	return result
}

func deltaToCSV(delta Delta) string {
	if list, ok := delta.(List); ok {
		return deltasToCSV(list.Deltas)
	}

	deltaJSON, err := delta.MarshalJSON()
	if err != nil {
		return ""
	}

	return "," + string(deltaJSON)
}

func strSliceToQuotedCSV(arr []string) string {
	result := ""

	for _, str := range arr {
		result += "," + strconv.Quote(str)
	}

	return result
}
