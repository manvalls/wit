package wit

import (
	"io"
	"strconv"
)

// Renderer renders relevant content to a writer
type Renderer interface {
	Render(io.Writer) error
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
	}

	return class
}

func strSliceToJSON(arr []string) string {
	result := "["

	for i, str := range arr {
		if i != 0 {
			result += ","
		}

		result += strconv.Quote(str)
	}

	result += "]"
	return result
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
