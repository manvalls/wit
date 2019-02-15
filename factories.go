package wit

import (
	"io"
	"sync"

	"golang.org/x/net/html"
)

func extractDeltas(commands []Command) []Delta {
	filteredDeltas := make([]Delta, 0, len(commands))
	for _, command := range commands {
		if !IsNil(command) {
			delta := command.Delta()

			switch delta.typeID {
			case sliceType:
				childDeltas := delta.info.deltas
				for _, childDelta := range childDeltas {
					filteredDeltas = append(filteredDeltas, childDelta)
				}
			default:
				filteredDeltas = append(filteredDeltas, delta)
			}
		}
	}

	return filteredDeltas
}

// List groups a list of commands together
func List(commands ...Command) *ListCommand {
	return &ListCommand{commands: commands}
}

// ListCommand holds a list of commands
type ListCommand struct {
	mutex    sync.Mutex
	commands []Command
}

// Delta computes the delta based on the list of stored commands
func (l *ListCommand) Delta() Delta {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	filteredDeltas := extractDeltas(l.commands)

	switch len(filteredDeltas) {
	case 0:
		return Delta{}
	case 1:
		return filteredDeltas[0]
	default:
		return Delta{sliceType, &deltaInfo{deltas: filteredDeltas}}
	}
}

// Add adds a command to this list
func (l *ListCommand) Add(c Command) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.commands = append(l.commands, c)
}

// Root applies given commands to the root of the document
func Root(commands ...Command) Command {
	deltas := extractDeltas(commands)
	if len(deltas) == 0 {
		return Nil
	}

	return Delta{rootType, &deltaInfo{deltas: deltas}}
}

// Nil represents an effectless command
var Nil Command = Delta{}

// Remove removes from the document matching elements
var Remove Command = Delta{removeType, nil}

// Clear empties matching elements
var Clear Command = Delta{clearType, nil}

// Factory builds HTML documents on demand
type Factory interface {
	HTML() io.Reader
	Nodes(context *html.Node) []*html.Node
}

// HTML sets the inner HTML of the matching elements
func HTML(factory Factory) Command {
	return Delta{htmlType, &deltaInfo{factory: factory}}
}

// Parent applies provided commands to the parent of matching elements
func Parent(commands ...Command) Command {
	deltas := extractDeltas(commands)
	if len(deltas) == 0 {
		return Nil
	}

	return Delta{parentType, &deltaInfo{deltas: deltas}}
}

// FirstChild applies provided commands to the first child of matching elements
func FirstChild(commands ...Command) Command {
	deltas := extractDeltas(commands)
	if len(deltas) == 0 {
		return Nil
	}

	return Delta{firstChildType, &deltaInfo{deltas: deltas}}
}

// LastChild applies provided commands to the last child of matching elements
func LastChild(commands ...Command) Command {
	deltas := extractDeltas(commands)
	if len(deltas) == 0 {
		return Nil
	}

	return Delta{lastChildType, &deltaInfo{deltas: deltas}}
}

// PrevSibling applies provided commands to the previous sibling of matching elements
func PrevSibling(commands ...Command) Command {
	deltas := extractDeltas(commands)
	if len(deltas) == 0 {
		return Nil
	}

	return Delta{prevSiblingType, &deltaInfo{deltas: deltas}}
}

// NextSibling applies provided commands to the previous sibling of matching elements
func NextSibling(commands ...Command) Command {
	deltas := extractDeltas(commands)
	if len(deltas) == 0 {
		return Nil
	}

	return Delta{nextSiblingType, &deltaInfo{deltas: deltas}}
}

// Replace replaces matching elements with the provided HTML
func Replace(factory Factory) Command {
	return Delta{replaceType, &deltaInfo{factory: factory}}
}

// Append adds the provided HTML at the end of matching elements
func Append(factory Factory) Command {
	return Delta{appendType, &deltaInfo{factory: factory}}
}

// Prepend adds the provided HTML at the beginning of matching elements
func Prepend(factory Factory) Command {
	return Delta{prependType, &deltaInfo{factory: factory}}
}

// InsertAfter inserts the provided HTML after matching elements
func InsertAfter(factory Factory) Command {
	return Delta{insertAfterType, &deltaInfo{factory: factory}}
}

// InsertBefore inserts the provided HTML before matching elements
func InsertBefore(factory Factory) Command {
	return Delta{insertBeforeType, &deltaInfo{factory: factory}}
}

// AddAttr adds the provided attributes to the matching elements
func AddAttr(attr map[string]string) Command {
	if len(attr) == 0 {
		return Nil
	}

	return Delta{addAttrType, &deltaInfo{strMap: attr}}
}

// SetAttr sets the attributes of the matching elements
func SetAttr(attr map[string]string) Command {
	return Delta{setAttrType, &deltaInfo{strMap: attr}}
}

// RmAttr removes the provided attributes from the matching elements
func RmAttr(attrs ...string) Command {
	if len(attrs) == 0 {
		return Nil
	}

	return Delta{rmAttrType, &deltaInfo{strList: attrs}}
}

// AddStyles adds the provided styles to the matching elements
func AddStyles(styles map[string]string) Command {
	if len(styles) == 0 {
		return Nil
	}

	return Delta{addStylesType, &deltaInfo{strMap: styles}}
}

// RmStyles removes the provided styles from the matching elements
func RmStyles(styles ...string) Command {
	if len(styles) == 0 {
		return Nil
	}

	return Delta{rmStylesType, &deltaInfo{strList: styles}}
}

// AddClass adds the provided class to the matching elements
func AddClass(class string) Command {
	return Delta{addClassType, &deltaInfo{class: class}}
}

// RmClass removes the provided class from the matching elements
func RmClass(class string) Command {
	return Delta{rmClassType, &deltaInfo{class: class}}
}
