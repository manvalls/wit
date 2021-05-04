package wit

import (
	"bytes"
	"testing"
)

var delta = List{[]Delta{
	First{Body, List{[]Delta{
		HTML{HTMLFromString("<div class=\"one\"></div>")},
		Append{HTMLFromString("<div class=\"two\"></div>")},
		Prepend{HTMLFromString("<h1 foo=bar></h1>")},
		First{S(".two"), List{[]Delta{
			InsertAfter{HTMLFromString("<div class=\"three\" style=\"background: black;border: 1px solid\"></div>")},
			NextSibling{SetAttr{map[string]string{"foo": "bar"}}},
			PrevSibling{ReplaceAttr{map[string]string{"number": "one"}}},
			Parent{
				AddClasses{"foo bar"},
			},
			Root{
				First{Body, List{[]Delta{
					RmClasses{"baz foo"},
					FirstChild{SetStyles{map[string]string{"color": "black"}}},
					LastChild{RmStyles{[]string{"background"}}},
				}}},
			}}},
		},
		All{S("div"), InsertBefore{HTMLFromString("<hr>")}},
		First{S("h1"), List{[]Delta{
			ReplaceAttr{map[string]string{"bar": "foo"}},
			SetAttr{map[string]string{"bar2": "foo2"}},
			RmAttr{[]string{"bar2"}},
			SetAttr{map[string]string{"bar3": "foo3"}},
		}}},
	}}},
}}

var expectedJSON = `[1,[3,"body",[12,"<div class=\"one\"></div>"],[14,"<div class=\"two\"></div>"],[15,"<h1 foo=bar></h1>"],[3,".two",[16,"<div class=\"three\" style=\"background: black;border: 1px solid\"></div>"],[9,[18,{"foo":"bar"}]],[8,[19,{"number":"one"}]],[5,[23,"foo bar"]],[2,[3,"body",[24,"baz foo"],[6,[21,{"color":"black"}]],[7,[22,"background"]]]]],[4,"div",[17,"<hr>"]],[3,"h1",[19,{"bar":"foo"}],[18,{"bar2":"foo2"}],[20,"bar2"],[18,{"bar3":"foo3"}]]]]`
var baseDocument = NewDocument()
var expectedHTML = `<!DOCTYPE html><html><head></head><body class="bar"><h1 bar="foo" bar3="foo3"></h1><hr/><div number="one"></div><hr/><div class="two"></div><hr/><div class="three" style="border: 1px solid;" foo="bar"></div></body></html>`

func TestJSONMarshal(t *testing.T) {
	result, _ := delta.MarshalJSON()

	if string(result) != expectedJSON {
		t.Error("Expected " + expectedJSON + ", got" + string(result))
	}
}

func TestJSONUnmarshal(t *testing.T) {
	var result List
	(&result).UnmarshalJSON([]byte(expectedJSON))
	marshaledDelta, _ := delta.MarshalJSON()
	marshaledResult, _ := result.MarshalJSON()

	if string(marshaledDelta) != string(marshaledResult) {
		t.Error("Expected ", string(marshaledDelta), ", got", string(marshaledResult))
	}
}

func TestApply(t *testing.T) {
	var b bytes.Buffer

	First{
		Body,
		AddClasses{"foo baz"},
	}.Apply(baseDocument)

	delta.Apply(baseDocument)
	baseDocument.Render(&b)

	if string(b.Bytes()) != expectedHTML {
		t.Error("Expected ", expectedHTML, ", got", string(b.Bytes()))
	}
}
