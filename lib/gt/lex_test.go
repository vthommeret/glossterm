package gt

import (
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestLex(t *testing.T) {
	tests := []struct {
		input string
		desc  string
		want  []item
		debug bool
	}{
		{"hello", "Simple text", []item{
			it(itemText, "hello"),
		}, false},
		{"==header==", "Header", []item{
			ih(itemHeaderStart, "==", 2),
			it(itemText, "header"),
			ih(itemHeaderEnd, "==", 2),
		}, false},
		{"==header == ignore inner == delimiters==", "Header ignore inner delimiters", []item{
			ih(itemHeaderStart, "==", 2),
			it(itemText, "header == ignore inner == delimiters"),
			ih(itemHeaderEnd, "==", 2),
		}, false},
		{"==header1==\n\nSome text\n\n==header2==\n\nSome more text", "Headers and text", []item{
			ih(itemHeaderStart, "==", 2),
			it(itemText, "header1"),
			ih(itemHeaderEnd, "==", 2),
			it(itemText, "\n\nSome text\n\n"),
			ih(itemHeaderStart, "==", 2),
			it(itemText, "header2"),
			ih(itemHeaderEnd, "==", 2),
			it(itemText, "\n\nSome more text"),
		}, false},
		{"==Header {{t}}==", "Header with action", []item{
			ih(itemHeaderStart, "==", 2),
			it(itemText, "Header "),
			it(itemLeftTemplate, "{{"),
			it(itemAction, "t"),
			it(itemRightTemplate, "}}"),
			ih(itemHeaderEnd, "==", 2),
		}, false},
		{"{{t}}", "Simple action", []item{
			it(itemLeftTemplate, "{{"),
			it(itemAction, "t"),
			it(itemRightTemplate, "}}"),
		}, false},
		{"{{gloss|hello}}", "Action, one positional param", []item{
			it(itemLeftTemplate, "{{"),
			it(itemAction, "gloss"),
			it(itemParamDelim, "|"),
			it(itemText, "hello"),
			it(itemRightTemplate, "}}"),
		}, false},
		{"{{t|1|2}}", "Action, two positional params", []item{
			it(itemLeftTemplate, "{{"),
			it(itemAction, "t"),
			it(itemParamDelim, "|"),
			it(itemText, "1"),
			it(itemParamDelim, "|"),
			it(itemText, "2"),
			it(itemRightTemplate, "}}"),
		}, false},
		{"{{t|1||3}}", "Action, empty param", []item{
			it(itemLeftTemplate, "{{"),
			it(itemAction, "t"),
			it(itemParamDelim, "|"),
			it(itemText, "1"),
			it(itemParamDelim, "|"),
			it(itemParamDelim, "|"),
			it(itemText, "3"),
			it(itemRightTemplate, "}}"),
		}, false},
		{"{{t|1|a=2}}", "Action, positional param, named param", []item{
			it(itemLeftTemplate, "{{"),
			it(itemAction, "t"),
			it(itemParamDelim, "|"),
			it(itemText, "1"),
			it(itemParamDelim, "|"),
			it(itemParamName, "a"),
			it(itemText, "2"),
			it(itemRightTemplate, "}}"),
		}, false},

		{"A [[simple]] link", "A simple link", []item{
			it(itemText, "A "),
			it(itemLeftLink, "[["),
			it(itemLink, "simple"),
			it(itemRightLink, "]]"),
			it(itemText, " link"),
		}, false},
		/*
			{"An [[unclosed link", "An unclosed link", []item{
				it(itemText, "An [[unclosed link"),
			}, true},
		*/
		{"A [[simple|named]] link", "A simple, named link", []item{
			it(itemText, "A "),
			it(itemLeftLink, "[["),
			it(itemLink, "simple"),
			it(itemLinkDelim, "|"),
			it(itemText, "named"),
			it(itemRightLink, "]]"),
			it(itemText, " link"),
		}, false},
		/*
			{"An [[unclosed|named] link", "An unclosed, named link", []item{
				it(itemText, "An "),
				it(itemText, "[[unclosed|named] link"),
			}, false},
		*/
		{"An [[embedded|named{{template}}]] link", "Embedded template in named link", []item{
			it(itemText, "An "),
			it(itemLeftLink, "[["),
			it(itemLink, "embedded"),
			it(itemLinkDelim, "|"),
			it(itemText, "named"),
			it(itemLeftTemplate, "{{"),
			it(itemAction, "template"),
			it(itemRightTemplate, "}}"),
			it(itemRightLink, "]]"),
			it(itemText, " link"),
		}, false},
		{"A [[simple|'''strong''']] link", "A strong link", []item{
			it(itemText, "A "),
			it(itemLeftLink, "[["),
			it(itemLink, "simple"),
			it(itemLinkDelim, "|"),
			it(itemStrong, "'''"),
			it(itemText, "strong"),
			it(itemStrong, "'''"),
			it(itemRightLink, "]]"),
			it(itemText, " link"),
		}, false},
		/*
			{"A [[multi\nline]] link", "A multiline link", []item{
				// TODO: Try to merge next two tokens?
				it(itemText, "A "),
				it(itemText, "[[multi\n"),
				it(itemText, "line]] link"),
			}, false},
			{"A [[simple [[nested]]]] link", "A simple nested link", []item{
				it(itemText, "A "),
				it(itemText, "[[simple "),
				it(itemLeftLink, "[["),
				it(itemLink, "nested"),
				it(itemRightLink, "]]"),
				it(itemText, "]] link"),
			}, false},
			{"An [[unclosed link", "An unclosed link", []item{
				it(itemText, "An "),
				it(itemText, "[[unclosed link"),
			}, false},
			{"A [[partially closed] link", "A partially closed link", []item{
				it(itemText, "A "),
				it(itemText, "[[partially closed] link"),
			}, false},
			{"An [[embedded {{template}}]]", "An embedded template", []item{
				it(itemText, "An "),
				it(itemText, "[[embedded "),
				it(itemLeftTemplate, "{{"),
				it(itemAction, "template"),
				it(itemRightTemplate, "}}"),
				it(itemText, "]]"),
			}, false},
		*/

		{"{{gloss|[[was|Was]] I?}}", "Pipe in link in template", []item{
			it(itemLeftTemplate, "{{"),
			it(itemAction, "gloss"),
			it(itemParamDelim, "|"),
			it(itemLeftLink, "[["),
			it(itemLink, "was"),
			it(itemLinkDelim, "|"),
			it(itemText, "Was"),
			it(itemRightLink, "]]"),
			it(itemText, " I?"),
			it(itemRightTemplate, "}}"),
		}, false},
		{"{{t|<डलर>}}", "Non-ASCII tag name", []item{
			it(itemLeftTemplate, "{{"),
			it(itemAction, "t"),
			it(itemParamDelim, "|"),
			it(itemText, "<डलर>"),
			it(itemRightTemplate, "}}"),
		}, false},
		{"{{m|la|dictus{{m|la|dictus}}}}", "Nested templates", []item{
			it(itemLeftTemplate, "{{"),
			it(itemAction, "m"),
			it(itemParamDelim, "|"),
			it(itemText, "la"),
			it(itemParamDelim, "|"),
			it(itemText, "dictus"),
			it(itemLeftTemplate, "{{"),
			it(itemAction, "m"),
			it(itemParamDelim, "|"),
			it(itemText, "la"),
			it(itemParamDelim, "|"),
			it(itemText, "dictus"),
			it(itemRightTemplate, "}}"),
			it(itemRightTemplate, "}}"),
		}, false},
		/*
			{"start{{t\nend", "Unclosed template action", []item{
				it(itemText, "start"),
				it(itemText, "{{t"),
				it(itemText, "end"),
			}, false},
			{"start{{t\n==Header", "Unclosed template action (header)", []item{
				it(itemText, "start"),
				it(itemText, "{{t"),
				it(itemText, "\n"),
				ih(itemHeaderStart, "==", 2),
				it(itemText, "Header"),
			}, false},
			{"start{{t|1\nmiddle\nend", "Unclosed template param", []item{
				it(itemText, "start"),
				it(itemText, "{{t|1\nmiddle\nend"),
			}, false},
			{"start{{t|1\n==Header", "Unclosed template param (header)", []item{
				it(itemText, "start"),
				it(itemText, "{{t|1"),
				it(itemText, "\n"),
				ih(itemHeaderStart, "==", 2),
				it(itemText, "Header"),
			}, false},
			{"start{{t|1\n{{new}}", "Unclosed template param (nested)", []item{
				it(itemText, "start"),
				it(itemText, "{{t|1\n"),
				it(itemLeftTemplate, "{{"),
				it(itemAction, "new"),
				it(itemRightTemplate, "}}"),
			}, false},
			{"{{t|multi\nline}}", "Multi-line template", []item{
				it(itemLeftTemplate, "{{"),
				it(itemAction, "t"),
				it(itemParamDelim, "|"),
				it(itemParamText, "multi\nline"),
				it(itemRightTemplate, "}}"),
			}, false},
			{"{{t\n\t\t|1}}", "Multi-line template w/ leading whitespace", []item{
				it(itemLeftTemplate, "{{"),
				it(itemAction, "t"),
				it(itemParamDelim, "\n\t\t|"),
				it(itemParamText, "1"),
				it(itemRightTemplate, "}}"),
			}, false},
		*/
		// TODO: Re-handle item prefixes
		/*
			{"====Descendants====\n* English: [[lettuce]]", "List items", []item{
				ih(itemHeaderStart, "====", 4),
				it(itemText, "Descendants"),
				ih(itemHeaderEnd, "====", 4),
				it(itemText, "\n"),
				it(itemUnorderedListItemStart, "*"),
				it(itemListItemPrefix, "English"),
				it(itemListItemEnd, ": "),
				it(itemLeftLink, "[["),
				it(itemLink, "lettuce"),
				it(itemRightLink, "]]"),
			}, false},
		*/
		/*
			{"Text with * asterisk ignored.", "Ignored asterisk", []item{
				it(itemText, "Text with * asterisk ignored."),
			}, false},
		*/

		// Markup tests
		{"A '''''strong emphasized''''' statement", "Strong emphasized text", []item{
			it(itemText, "A "),
			it(itemStrongEmphasized, "'''''"),
			it(itemText, "strong emphasized"),
			it(itemStrongEmphasized, "'''''"),
			it(itemText, " statement"),
		}, false},
		{"A '''strong''' statement", "Strong text", []item{
			it(itemText, "A "),
			it(itemStrong, "'''"),
			it(itemText, "strong"),
			it(itemStrong, "'''"),
			it(itemText, " statement"),
		}, false},
		{"An ''emphasized'' statement", "Emphasized text", []item{
			it(itemText, "An "),
			it(itemEmphasized, "''"),
			it(itemText, "emphasized"),
			it(itemEmphasized, "''"),
			it(itemText, " statement"),
		}, false},
	}

	for _, tt := range tests {
		l := NewLexer2(tt.input, true)
		l.debug = tt.debug

		want := l.String(tt.want)
		got := l.String(l.Items())

		if want != got {
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(want, got, false)
			diffText := dmp.DiffPrettyText(diffs)
			t.Errorf("%s: NewLexer(%q) diff: %s", tt.desc, tt.input, diffText)
		}
	}
}

func it(typ itemType, val string) item {
	return item{typ: typ, val: val}
}

func ih(typ itemType, val string, depth int) item {
	return item{typ: typ, val: val, depth: depth}
}
