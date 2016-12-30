package gt

import (
	"fmt"
	"log"
	"testing"
)

func TestLex(t *testing.T) {
	tests := []struct {
		input string
		desc  string
		want  []item
	}{
		{"hello", "Simple text", []item{
			it(itemText, "\n"),
			it(itemText, "hello\n"),
		}},
		{"==header==", "Header", []item{
			ih(itemHeaderStart, "\n==", 2),
			it(itemText, "header"),
			ih(itemHeaderEnd, "==\n", 2),
		}},
		{"{{t}}", "Simple action", []item{
			it(itemText, "\n"),
			it(itemLeftTemplate, "{{"),
			it(itemAction, "t"),
			it(itemRightTemplate, "}}"),
			it(itemText, "\n"),
		}},
		{"\n{{t|1|2}}", "Action, two positional params", []item{
			it(itemText, "\n"),
			it(itemLeftTemplate, "{{"),
			it(itemAction, "t"),
			it(itemParamDelim, "|"),
			it(itemParamText, "1"),
			it(itemParamDelim, "|"),
			it(itemParamText, "2"),
			it(itemRightTemplate, "}}"),
			it(itemText, "\n"),
		}},
		{"\n{{t|1||3}}", "Action, empty param", []item{
			it(itemText, "\n"),
			it(itemLeftTemplate, "{{"),
			it(itemAction, "t"),
			it(itemParamDelim, "|"),
			it(itemParamText, "1"),
			it(itemParamDelim, "|"),
			it(itemParamText, ""),
			it(itemParamDelim, "|"),
			it(itemParamText, "3"),
			it(itemRightTemplate, "}}"),
			it(itemText, "\n"),
		}},
		{"{{t|1|a=2}}", "Action, positional param, named param", []item{
			it(itemText, "\n"),
			it(itemLeftTemplate, "{{"),
			it(itemAction, "t"),
			it(itemParamDelim, "|"),
			it(itemParamText, "1"),
			it(itemParamDelim, "|"),
			it(itemParamName, "a"),
			it(itemParamText, "2"),
			it(itemRightTemplate, "}}"),
			it(itemText, "\n"),
		}},
		{"{{t|<strong>e=mc^2</strong>}}", "Equals in tag", []item{
			it(itemText, "\n"),
			it(itemLeftTemplate, "{{"),
			it(itemAction, "t"),
			it(itemParamDelim, "|"),
			it(itemParamText, "<strong>e=mc^2</strong>"),
			it(itemRightTemplate, "}}"),
			it(itemText, "\n"),
		}},
		{"{{t|[[was|Was] I?}}", "Pipe in link", []item{
			it(itemText, "\n"),
			it(itemLeftTemplate, "{{"),
			it(itemAction, "t"),
			it(itemParamDelim, "|"),
			it(itemParamText, "[[was|Was] I?"),
			it(itemRightTemplate, "}}"),
			it(itemText, "\n"),
		}},
		{"{{t|<डलर>}}", "Non-ASCII tag name", []item{
			it(itemText, "\n"),
			it(itemLeftTemplate, "{{"),
			it(itemAction, "t"),
			it(itemParamDelim, "|"),
			it(itemParamText, "<डलर>"),
			it(itemRightTemplate, "}}"),
			it(itemText, "\n"),
		}},
		{"{{m|la|dictus{{m|la|dictus}}}}", "Nested templates", []item{
			it(itemText, "\n"),
			it(itemLeftTemplate, "{{"),
			it(itemAction, "m"),
			it(itemParamDelim, "|"),
			it(itemParamText, "la"),
			it(itemParamDelim, "|"),
			it(itemParamText, "dictus"),
			it(itemLeftTemplate, "{{"),
			it(itemAction, "m"),
			it(itemParamDelim, "|"),
			it(itemParamText, "la"),
			it(itemParamDelim, "|"),
			it(itemParamText, "dictus"),
			it(itemRightTemplate, "}}"),
			it(itemRightTemplate, "}}"),
			it(itemText, "\n"),
		}},
		{"start{{t\nend", "Unclosed template action", []item{
			it(itemText, "\n"),
			it(itemText, "start"),
			it(itemText, "{{t"),
			it(itemText, "end\n"),
		}},
		{"start{{t\n==Header", "Unclosed template action (header)", []item{
			it(itemText, "\n"),
			it(itemText, "start"),
			it(itemText, "{{t"),
			ih(itemHeaderStart, "\n==", 2),
			it(itemText, "Header\n"),
		}},
		{"start{{t|1\nmiddle\nend", "Unclosed template param", []item{
			it(itemText, "\n"),
			it(itemText, "start"),
			it(itemText, "{{t|1\nmiddle\nend\n"),
		}},
		{"start{{t|1\n==Header", "Unclosed template param (header)", []item{
			it(itemText, "\n"),
			it(itemText, "start"),
			it(itemText, "{{t|1"),
			ih(itemHeaderStart, "\n==", 2),
			it(itemText, "Header\n"),
		}},
		{"start{{t|1\n{{new}}", "Unclosed template param (nested)", []item{
			it(itemText, "\n"),
			it(itemText, "start"),
			it(itemText, "{{t|1\n"),
			it(itemLeftTemplate, "{{"),
			it(itemAction, "new"),
			it(itemRightTemplate, "}}"),
			it(itemText, "\n"),
		}},
		{"{{t|multi\nline}}", "Multi-line template", []item{
			it(itemText, "\n"),
			it(itemLeftTemplate, "{{"),
			it(itemAction, "t"),
			it(itemParamDelim, "|"),
			it(itemParamText, "multi\nline"),
			it(itemRightTemplate, "}}"),
			it(itemText, "\n"),
		}},
		{"{{t\n\t\t|1}}", "Multi-line template w/ leading whitespace", []item{
			it(itemText, "\n"),
			it(itemLeftTemplate, "{{"),
			it(itemAction, "t"),
			it(itemParamDelim, "\n\t\t|"),
			it(itemParamText, "1"),
			it(itemRightTemplate, "}}"),
			it(itemText, "\n"),
		}},
		{"{{gloss|outer{{gloss|inner{{gloss|inner-inner}} and some more {{gloss|ok}} and even more {{gloss|ok}}", "Crazy nesting 1", []item{
			it(itemText, "\n"),
			it(itemText, "{{gloss|outer{{gloss|inner"),
			it(itemLeftTemplate, "{{"),
			it(itemAction, "gloss"),
			it(itemParamDelim, "|"),
			it(itemParamText, "inner-inner"),
			it(itemRightTemplate, "}}"),
			it(itemText, " and some more "),
			it(itemLeftTemplate, "{{"),
			it(itemAction, "gloss"),
			it(itemParamDelim, "|"),
			it(itemParamText, "ok"),
			it(itemRightTemplate, "}}"),
			it(itemText, " and even more "),
			it(itemLeftTemplate, "{{"),
			it(itemAction, "gloss"),
			it(itemParamDelim, "|"),
			it(itemParamText, "ok"),
			it(itemRightTemplate, "}}"),
			it(itemText, "\n"),
		}},
		{"{{gloss|outer{{gloss|inner{{gloss|inner-inner}} and some more {{gloss|ok and even more {{gloss|ok}}", "Crazy nesting 2", []item{
			it(itemText, "\n"),
			it(itemText, "{{gloss|outer{{gloss|inner"),
			it(itemLeftTemplate, "{{"),
			it(itemAction, "gloss"),
			it(itemParamDelim, "|"),
			it(itemParamText, "inner-inner"),
			it(itemRightTemplate, "}}"),
			it(itemText, " and some more {{gloss|ok and even more "),
			it(itemLeftTemplate, "{{"),
			it(itemAction, "gloss"),
			it(itemParamDelim, "|"),
			it(itemParamText, "ok"),
			it(itemRightTemplate, "}}"),
			it(itemText, "\n"),
		}},
	}
	for _, tt := range tests {
		l := NewLexer(tt.input)
		got := items(l)
		itemsEqual(got, tt.want)
		if eq, reason := itemsEqual(got, tt.want); !eq {
			t.Errorf("%s: %s. NewLexer(%q) = %+v, want %+v.", tt.desc, reason, tt.input, got, tt.want)
		}
	}
}

func it(typ itemType, val string) item {
	return item{typ: typ, val: val}
}

func ih(typ itemType, val string, depth int) item {
	return item{typ: typ, val: val, depth: depth}
}

func items(l *lexer) (is []item) {
	for {
		i := l.NextItem()
		if i.typ == itemEOF || i.typ == itemError {
			if i.typ == itemError {
				log.Fatalf("Error getting next item: %s", i.val)
			}
			break
		}
		is = append(is, i)
	}
	return is
}

func itemsEqual(is1, is2 []item) (eq bool, reason string) {
	n1 := len(is1)
	n2 := len(is2)
	if n1 != n2 {
		return false, fmt.Sprintf("Length = %d, want %d", n1, n2)
	}
	for i := range is1 {
		i1 := is1[i]
		i2 := is2[i]
		j := i + 1
		if i1.typ != i2.typ {
			return false, fmt.Sprintf("Item #%d typ = %q, want %q", j, i1.typ, i2.typ)
		} else if i1.val != i2.val {
			return false, fmt.Sprintf("Item #%d val = %q, want %q", j, i1.val, i2.val)
		} else if i1.typ == itemHeaderStart || i1.typ == itemHeaderEnd {
			if i1.depth != i2.depth {
				return false, fmt.Sprintf("Header item #%d depth = %d, want %d", j, i1.depth, i2.depth)
			}
		}
	}
	return true, ""
}
