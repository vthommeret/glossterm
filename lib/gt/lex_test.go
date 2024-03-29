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
		{"A [[simple|named]] link", "A simple, named link", []item{
			it(itemText, "A "),
			it(itemLeftLink, "[["),
			it(itemLink, "simple"),
			it(itemLinkDelim, "|"),
			it(itemText, "named"),
			it(itemRightLink, "]]"),
			it(itemText, " link"),
		}, false},
		{"An [[unclosed|named] link", "An unclosed, named link", []item{
			it(itemText, "An "),
			it(itemLeftLink, "[["),
			it(itemLink, "unclosed"),
			it(itemLinkDelim, "|"),
			it(itemText, "named] link"),
		}, false},
		{"A [[simple|'''strong''']] link", "Markup in link", []item{
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
		{"An [[amazing|<em>emphasized</em>]] link", "HTML in link", []item{
			it(itemText, "An "),
			it(itemLeftLink, "[["),
			it(itemLink, "amazing"),
			it(itemLinkDelim, "|"),
			it(itemOpenTagLeft, "<"),
			it(itemTagName, "em"),
			it(itemTagRight, ">"),
			it(itemText, "emphasized"),
			it(itemCloseTagLeft, "</"),
			it(itemTagName, "em"),
			it(itemTagRight, ">"),
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
		*/
		/*
			TODO: Nested links not allowed
			{"A [[simple [[nested]]]] link", "A simple nested link", []item{
				it(itemText, "A "),
				it(itemText, "[[simple "),
				it(itemLeftLink, "[["),
				it(itemLink, "nested"),
				it(itemRightLink, "]]"),
				it(itemText, "]] link"),
			}, false},
		*/
		{"An [[unclosed link", "An unclosed link", []item{
			it(itemText, "An "),
			it(itemLeftLink, "[["),
			it(itemText, "unclosed link"),
		}, false},
		{"A [[partially closed] link", "A partially closed link", []item{
			it(itemText, "A "),
			it(itemLeftLink, "[["),
			it(itemText, "partially closed] link"),
		}, false},
		/*
			// TODO: Template only allowed in text section
				{"An [[embedded {{template}}]]", "Template in link (not allowed)", []item{
					it(itemText, "An "),
					it(itemText, "[[embedded "),
					it(itemLeftTemplate, "{{"),
					it(itemAction, "template"),
					it(itemRightTemplate, "}}"),
					it(itemText, "]]"),
				}, false},
		*/
		{"An [[link|embedded {{template}}]]", "Template in link text (allowed)", []item{
			it(itemText, "An "),
			it(itemLeftLink, "[["),
			it(itemLink, "link"),
			it(itemLinkDelim, "|"),
			it(itemText, "embedded "),
			it(itemLeftTemplate, "{{"),
			it(itemAction, "template"),
			it(itemRightTemplate, "}}"),
			it(itemRightLink, "]]"),
		}, false},

		{"{{t|<डलर>}}", "Non-ASCII tag name", []item{
			it(itemLeftTemplate, "{{"),
			it(itemAction, "t"),
			it(itemParamDelim, "|"),
			it(itemOpenTagLeft, "<"),
			it(itemTagName, "डलर"),
			it(itemTagRight, ">"),
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

		// HTML tests
		{"An <em>emphasized</em> test", "HTML emphasized text", []item{
			it(itemText, "An "),
			it(itemOpenTagLeft, "<"),
			it(itemTagName, "em"),
			it(itemTagRight, ">"),
			it(itemText, "emphasized"),
			it(itemCloseTagLeft, "</"),
			it(itemTagName, "em"),
			it(itemTagRight, ">"),
			it(itemText, " test"),
		}, false},
		{"An <span style=\"color: red\">attribute</span> test", "HTML attribute text", []item{
			it(itemText, "An "),
			it(itemOpenTagLeft, "<"),
			it(itemTagName, "span"),
			it(itemTagAttrName, "style"),
			it(itemTagAttrValue, "color: red"),
			it(itemTagRight, ">"),
			it(itemText, "attribute"),
			it(itemCloseTagLeft, "</"),
			it(itemTagName, "span"),
			it(itemTagRight, ">"),
			it(itemText, " test"),
		}, false},
		{"Multiple <span attr1=\"1\" attr2=\"2\">attribute</span> test", "HTML multiple attribute text", []item{
			it(itemText, "Multiple "),
			it(itemOpenTagLeft, "<"),
			it(itemTagName, "span"),
			it(itemTagAttrName, "attr1"),
			it(itemTagAttrValue, "1"),
			it(itemTagAttrName, "attr2"),
			it(itemTagAttrValue, "2"),
			it(itemTagRight, ">"),
			it(itemText, "attribute"),
			it(itemCloseTagLeft, "</"),
			it(itemTagName, "span"),
			it(itemTagRight, ">"),
			it(itemText, " test"),
		}, false},
		// TODO: Uncomment when emitTrim usage is eliminated
		/*
			{"A <span   style  =   \"color: red\"   >whitespace</span> test", "HTML whitespace text", []item{
				it(itemText, "A "),
				it(itemOpenTagLeft, "<"),
				it(itemTagName, "span"),
				it(itemTagAttrName, "style"),
				it(itemTagAttrValue, "color: red"),
				it(itemTagRight, ">"),
				it(itemText, "whitespace"),
				it(itemCloseTagLeft, "</"),
				it(itemTagName, "span"),
				it(itemTagRight, ">"),
				it(itemText, " test"),
			}, false},
		*/
		{"This is an <!--HTML--> comment", "HTML comment", []item{
			it(itemText, "This is an "),
			it(itemTagCommentLeft, "<!--"),
			it(itemTagComment, "HTML"),
			it(itemTagCommentRight, "-->"),
			it(itemText, " comment"),
		}, false},
		{"<br/>", "HTML void tag", []item{
			it(itemOpenTagLeft, "<"),
			it(itemTagName, "br"),
			it(itemCloseTagRight, "/>"),
		}, false},
		{"<input readonly />", "HTML void tag valueless attribute", []item{
			it(itemOpenTagLeft, "<"),
			it(itemTagName, "input"),
			it(itemTagAttrName, "readonly"),
			it(itemCloseTagRight, "/>"),
		}, false},
		{"<ref name=\"SOED\"/>", "HTML void tag attribute", []item{
			it(itemOpenTagLeft, "<"),
			it(itemTagName, "ref"),
			it(itemTagAttrName, "name"),
			it(itemTagAttrValue, "SOED"),
			it(itemCloseTagRight, "/>"),
		}, false},
		{"<ref name=SOED/>", "HTML void tag unquoted attribute", []item{
			it(itemOpenTagLeft, "<"),
			it(itemTagName, "ref"),
			it(itemTagAttrName, "name"),
			it(itemTagAttrValue, "SOED"),
			it(itemCloseTagRight, "/>"),
		}, false},
		{"<span valueless>Blah</span>", "HTML valueless attribute", []item{
			it(itemOpenTagLeft, "<"),
			it(itemTagName, "span"),
			it(itemTagAttrName, "valueless"),
			it(itemTagRight, ">"),
			it(itemText, "Blah"),
			it(itemCloseTagLeft, "</"),
			it(itemTagName, "span"),
			it(itemTagRight, ">"),
		}, false},
		{"<ref name=OCD>Lindberg</ref>", "HTML unquoted attribute values", []item{
			it(itemOpenTagLeft, "<"),
			it(itemTagName, "ref"),
			it(itemTagAttrName, "name"),
			it(itemTagAttrValue, "OCD"),
			it(itemTagRight, ">"),
			it(itemText, "Lindberg"),
			it(itemCloseTagLeft, "</"),
			it(itemTagName, "ref"),
			it(itemTagRight, ">"),
		}, false},
		// TODO: Need to support balancing quotes / ghost quotes for below cases
		/*
			{"<ref name=   OCD  >Lindberg</ref>", "HTML unquoted attribute values whitespace", []item{
				it(itemOpenTagLeft, "<"),
				it(itemTagName, "ref"),
				it(itemTagAttrName, "name"),
				it(itemTagAttrValue, "OCD"),
				it(itemTagRight, ">"),
				it(itemText, "Lindberg"),
				it(itemCloseTagLeft, "</"),
				it(itemTagName, "ref"),
				it(itemTagRight, ">"),
			}, true},
			{"<ref name=   OCD  class=\"foo\">Lindberg</ref>", "HTML unquoted multiple attribute values whitespace", []item{
				it(itemOpenTagLeft, "<"),
				it(itemTagName, "ref"),
				it(itemTagAttrName, "name"),
				it(itemTagAttrValue, "OCD"),
				it(itemTagAttrName, "class"),
				it(itemTagAttrValue, "foo"),
				it(itemTagRight, ">"),
				it(itemText, "Lindberg"),
				it(itemCloseTagLeft, "</"),
				it(itemTagName, "ref"),
				it(itemTagRight, ">"),
			}, true},
		*/

		// Nesting tests
		{"{{gloss|hello <strong>world</strong>}}", "HTML tag in template", []item{
			it(itemLeftTemplate, "{{"),
			it(itemAction, "gloss"),
			it(itemParamDelim, "|"),
			it(itemText, "hello "),
			it(itemOpenTagLeft, "<"),
			it(itemTagName, "strong"),
			it(itemTagRight, ">"),
			it(itemText, "world"),
			it(itemCloseTagLeft, "</"),
			it(itemTagName, "strong"),
			it(itemTagRight, ">"),
			it(itemRightTemplate, "}}"),
		}, false},
		{"<strong>hello {{m|en|world}}</strong>", "Template in HTML tag", []item{
			it(itemOpenTagLeft, "<"),
			it(itemTagName, "strong"),
			it(itemTagRight, ">"),
			it(itemText, "hello "),
			it(itemLeftTemplate, "{{"),
			it(itemAction, "m"),
			it(itemParamDelim, "|"),
			it(itemText, "en"),
			it(itemParamDelim, "|"),
			it(itemText, "world"),
			it(itemRightTemplate, "}}"),
			it(itemCloseTagLeft, "</"),
			it(itemTagName, "strong"),
			it(itemTagRight, ">"),
		}, false},
		{"<strong>hello [[world|world!]]</strong>", "Link in HTML tag", []item{
			it(itemOpenTagLeft, "<"),
			it(itemTagName, "strong"),
			it(itemTagRight, ">"),
			it(itemText, "hello "),
			it(itemLeftLink, "[["),
			it(itemLink, "world"),
			it(itemLinkDelim, "|"),
			it(itemText, "world!"),
			it(itemRightLink, "]]"),
			it(itemCloseTagLeft, "</"),
			it(itemTagName, "strong"),
			it(itemTagRight, ">"),
		}, false},
		{"Hello '''[[world|world!]]'''", "Link in markup", []item{
			it(itemText, "Hello "),
			it(itemStrong, "'''"),
			it(itemLeftLink, "[["),
			it(itemLink, "world"),
			it(itemLinkDelim, "|"),
			it(itemText, "world!"),
			it(itemRightLink, "]]"),
			it(itemStrong, "'''"),
		}, false},
		{"{{m|en|world|gloss=<span style=\"color: red\">hello</span>}}", "HTML tag in named parameter", []item{
			it(itemLeftTemplate, "{{"),
			it(itemAction, "m"),
			it(itemParamDelim, "|"),
			it(itemText, "en"),
			it(itemParamDelim, "|"),
			it(itemText, "world"),
			it(itemParamDelim, "|"),
			it(itemParamName, "gloss"),
			it(itemOpenTagLeft, "<"),
			it(itemTagName, "span"),
			it(itemTagAttrName, "style"),
			it(itemTagAttrValue, "color: red"),
			it(itemTagRight, ">"),
			it(itemText, "hello"),
			it(itemCloseTagLeft, "</"),
			it(itemTagName, "span"),
			it(itemTagRight, ">"),
			it(itemRightTemplate, "}}"),
		}, false},
		{"{{gloss|hello ''world''}}", "Markup in template", []item{
			it(itemLeftTemplate, "{{"),
			it(itemAction, "gloss"),
			it(itemParamDelim, "|"),
			it(itemText, "hello "),
			it(itemEmphasized, "''"),
			it(itemText, "world"),
			it(itemEmphasized, "''"),
			it(itemRightTemplate, "}}"),
		}, false},
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
		{"This is an <!--{{gloss|HTML}}--> comment", "Template in HTML comment", []item{
			it(itemText, "This is an "),
			it(itemTagCommentLeft, "<!--"),
			it(itemTagComment, "{{gloss|HTML}}"),
			it(itemTagCommentRight, "-->"),
			it(itemText, " comment"),
		}, false},
		{"This is an <!--[[HTML]]--> comment", "Link in HTML comment", []item{
			it(itemText, "This is an "),
			it(itemTagCommentLeft, "<!--"),
			it(itemTagComment, "[[HTML]]"),
			it(itemTagCommentRight, "-->"),
			it(itemText, " comment"),
		}, false},
		{"{{gloss|This is a <!-- comment -->}}", "Comment in template", []item{
			it(itemLeftTemplate, "{{"),
			it(itemAction, "gloss"),
			it(itemParamDelim, "|"),
			it(itemText, "This is a "),
			it(itemTagCommentLeft, "<!--"),
			it(itemTagComment, " comment "),
			it(itemTagCommentRight, "-->"),
			it(itemRightTemplate, "}}"),
		}, false},

		// Balance tests - Templates and links

		{"[[amazing|templated {{gloss|a note}}]]", "Balanced template in balanced link", []item{
			it(itemLeftLink, "[["),
			it(itemLink, "amazing"),
			it(itemLinkDelim, "|"),
			it(itemText, "templated "),
			it(itemLeftTemplate, "{{"),
			it(itemAction, "gloss"),
			it(itemParamDelim, "|"),
			it(itemText, "a note"),
			it(itemRightTemplate, "}}"),
			it(itemRightLink, "]]"),
		}, false},
		{"[[amazing|templated {{gloss|a note}}", "Balanced template in unbalanced link", []item{
			it(itemLeftLink, "[["),
			it(itemLink, "amazing"),
			it(itemLinkDelim, "|"),
			it(itemText, "templated "),
			it(itemLeftTemplate, "{{"),
			it(itemAction, "gloss"),
			it(itemParamDelim, "|"),
			it(itemText, "a note"),
			it(itemRightTemplate, "}}"),
		}, false},
		/*
			{"[[amazing|templated {{gloss|a note]]", "Unbalanced template in balanced link", []item{
				it(itemLeftLink, "[["),
				it(itemLink, "amazing"),
				it(itemLinkDelim, "|"),
				it(itemText, "templated {{gloss|a note"),
				it(itemRightLink, "]]"),
			}, false},
			{"[[amazing|templated {{gloss|a note", "Unbalanced template in unbalanced link", []item{
				it(itemText, "[[amazing|templated {{gloss|a note"),
			}, false},

			// Balance tests - Markup and links

			{"[[amazing|markup '''text''']]", "Balanced markup in balanced link", []item{
				it(itemLeftLink, "[["),
				it(itemLink, "amazing"),
				it(itemLinkDelim, "|"),
				it(itemText, "markup "),
				it(itemStrong, "'''"),
				it(itemText, "text"),
				it(itemStrong, "'''"),
				it(itemRightLink, "]]"),
			}, false},
			{"[[amazing|markup '''text'''", "Balanced markup in unbalanced link", []item{
				it(itemText, "[[amazing|markup "),
				it(itemStrong, "'''"),
				it(itemText, "text"),
				it(itemStrong, "'''"),
			}, false},
			{"[[amazing|markup '''text]]", "Unbalanced markup in balanced link", []item{
				it(itemLeftLink, "[["),
				it(itemLink, "amazing"),
				it(itemLinkDelim, "|"),
				it(itemText, "markup '''text"),
				it(itemRightLink, "]]"),
			}, false},
			{"[[amazing|markup '''text", "Unbalanced markup in unbalanced link", []item{
				it(itemText, "[[amazing|markup '''text"),
			}, false},

			// Balance tests - HTML tags and links

			{"[[amazing|html <span>text</span>]]", "Balanced HTML in balanced link", []item{
				it(itemLeftLink, "[["),
				it(itemLink, "amazing"),
				it(itemLinkDelim, "|"),
				it(itemText, "html "),
				it(itemOpenTagLeft, "<"),
				it(itemTagName, "span"),
				it(itemTagRight, ">"),
				it(itemText, "text"),
				it(itemCloseTagLeft, "</"),
				it(itemTagName, "span"),
				it(itemTagRight, ">"),
				it(itemRightLink, "]]"),
			}, false},
			{"[[amazing|html <span>text</span>", "Balanced HTML in unbalanced link", []item{
				it(itemText, "[[amazing|html "),
				it(itemOpenTagLeft, "<"),
				it(itemTagName, "span"),
				it(itemTagRight, ">"),
				it(itemText, "text"),
				it(itemCloseTagLeft, "</"),
				it(itemTagName, "span"),
				it(itemTagRight, ">"),
			}, false},
			{"[[amazing|html <span>text]]", "Unbalanced HTML in balanced link", []item{
				it(itemLeftLink, "[["),
				it(itemLink, "amazing"),
				it(itemLinkDelim, "|"),
				it(itemText, "html "),
				it(itemOpenTagLeft, "<"),
				it(itemTagName, "span"),
				it(itemTagRight, ">"),
				it(itemText, "text"),
				it(itemRightLink, "]]"),
			}, false},
			{"[[amazing|html <span>text", "Unbalanced HTML in unbalanced link", []item{
				it(itemText, "[[amazing|html <span>text"),
			}, false},
				{"[[amazing|html <span>text</em>]]", "Unbalanced HTML (incorrect tag) in balanced link", []item{}, false},

			// Balance tests - HTML opening delimiters and links

			{"[[amazing|html <span text]]", "Unbalanced HTML opening delimiters in balanced link", []item{
				it(itemLeftLink, "[["),
				it(itemLink, "amazing"),
				it(itemLinkDelim, "|"),
				it(itemText, "html <span text"),
				it(itemRightLink, "]]"),
			}, false},
			{"[[amazing|html <span text", "Unbalanced HTML opening delimiters in unbalanced link", []item{
				it(itemText, "[[amazing|html <span text"),
			}, false},

			// Balance tests - HTML closing delimiters and links

			{"[[amazing|html <span>text</span]]", "Unbalanced HTML closing delimiters in balanced link", []item{
				it(itemLeftLink, "[["),
				it(itemLink, "amazing"),
				it(itemLinkDelim, "|"),
				it(itemText, "html <span text"),
				it(itemRightLink, "]]"),
			}, false},
			{"[[amazing|html <span>text</span", "Unbalanced HTML closing delimiters in unbalanced link", []item{
				it(itemText, "[[amazing|html "),
				it(itemOpenTagLeft, "<"),
				it(itemTagName, "span"),
				it(itemTagRight, ">"),
				it(itemText, "text</span"),
			}, false},

			// Balance tests - HTML quoted attributes and links

			{"[[amazing|html <span class=\"test\">text</span>]]", "Balanced HTML quoted attributes in balanced link", []item{
				it(itemLeftLink, "[["),
				it(itemLink, "amazing"),
				it(itemLinkDelim, "|"),
				it(itemText, "html "),
				it(itemOpenTagLeft, "<"),
				it(itemTagName, "span"),
				it(itemTagAttrName, "class"),
				it(itemTagAttrValue, "test"),
				it(itemTagRight, ">"),
				it(itemText, "text"),
				it(itemCloseTagLeft, "</"),
				it(itemTagName, "span"),
				it(itemTagRight, ">"),
				it(itemRightLink, "]]"),
			}, false},
			{"[[amazing|html <span class=\"test\">text</span>", "Balanced HTML quoted attributes in unbalanced link", []item{
				it(itemText, "[[amazing|html "),
				it(itemOpenTagLeft, "<"),
				it(itemTagName, "span"),
				it(itemTagAttrName, "class"),
				it(itemTagAttrValue, "test"),
				it(itemTagRight, ">"),
				it(itemText, "text"),
				it(itemCloseTagLeft, "</"),
				it(itemTagName, "span"),
				it(itemTagRight, ">"),
			}, false},
			{"[[amazing|html <span class=\"test>text</span>]]", "Unbalanced HTML quoted attributes in balanced link", []item{
				it(itemLeftLink, "[["),
				it(itemLink, "amazing"),
				it(itemLinkDelim, "|"),
				it(itemText, "html "),
				it(itemOpenTagLeft, "<"),
				it(itemTagName, "span"),
				it(itemTagAttrName, "class"),
				it(itemTagAttrValue, "test"),
				it(itemTagRight, ">"),
				it(itemText, "text"),
				it(itemCloseTagLeft, "</"),
				it(itemTagName, "span"),
				it(itemTagRight, ">"),
				it(itemRightLink, "]]"),
			}, false},
			{"[[amazing|html <span class=\"test>text</span>", "Unbalanced HTML quoted attributes in unbalanced link", []item{
				it(itemText, "[[amazing|html "),
				it(itemOpenTagLeft, "<"),
				it(itemTagName, "span"),
				it(itemTagAttrName, "class"),
				it(itemTagAttrValue, "test"),
				it(itemTagRight, ">"),
				it(itemText, "text"),
				it(itemCloseTagLeft, "</"),
				it(itemTagName, "span"),
				it(itemTagRight, ">"),
			}, false},
			{"[[amazing|html <span class=\"test style=\"color: red\">text</span>", "Unbalanced HTML quoted multiple attributes in unbalanced link 1", []item{
				it(itemText, "[[amazing|html "),
				it(itemOpenTagLeft, "<"),
				it(itemTagName, "span"),
				it(itemTagAttrName, "class"),
				it(itemTagAttrValue, "test style="),
				it(itemTagRight, ">"),
				it(itemText, "text"),
				it(itemCloseTagLeft, "</"),
				it(itemTagName, "span"),
				it(itemTagRight, ">"),
			}, false},
			{"[[amazing|html <span class=\"test\" style=\"color: red>text</span>", "Unbalanced HTML quoted multiple attributes in unbalanced link 2", []item{
				it(itemText, "[[amazing|html "),
				it(itemOpenTagLeft, "<"),
				it(itemTagName, "span"),
				it(itemTagAttrName, "class"),
				it(itemTagAttrValue, "test"),
				it(itemTagAttrName, "style"),
				it(itemTagAttrValue, "color: red"),
				it(itemTagRight, ">"),
				it(itemText, "text"),
				it(itemCloseTagLeft, "</"),
				it(itemTagName, "span"),
				it(itemTagRight, ">"),
			}, false},

			// Balance tests - HTML unquoted attributes and links

			{"[[amazing|html <span class=\"test\" unquoted=blah>text</span>]]", "Balanced HTML quoted and unquoted attributes in balanced link", []item{
				it(itemLeftLink, "[["),
				it(itemLink, "amazing"),
				it(itemLinkDelim, "|"),
				it(itemText, "html "),
				it(itemOpenTagLeft, "<"),
				it(itemTagName, "span"),
				it(itemTagAttrName, "class"),
				it(itemTagAttrValue, "test"),
				it(itemTagAttrName, "unquoted"),
				it(itemTagAttrValue, "blah"),
				it(itemTagRight, ">"),
				it(itemText, "text"),
				it(itemCloseTagLeft, "</"),
				it(itemTagName, "span"),
				it(itemTagRight, ">"),
				it(itemRightLink, "]]"),
			}, false},
			{"[[amazing|html <span unquoted=blah class=\"test\">text</span>]]", "Balanced HTML unquoted and quoted attributes in balanced link", []item{
				it(itemLeftLink, "[["),
				it(itemLink, "amazing"),
				it(itemLinkDelim, "|"),
				it(itemText, "html "),
				it(itemOpenTagLeft, "<"),
				it(itemTagName, "span"),
				it(itemTagAttrName, "unquoted"),
				it(itemTagAttrValue, "blah"),
				it(itemTagAttrName, "class"),
				it(itemTagAttrValue, "test"),
				it(itemTagRight, ">"),
				it(itemText, "text"),
				it(itemCloseTagLeft, "</"),
				it(itemTagName, "span"),
				it(itemCloseTagRight, ">"),
				it(itemRightLink, "]]"),
			}, false},
			{"[[amazing|html <span unquoted class=\"test>text</span]]", "Unbalanced HTML unquoted and quoted attributes in balanced link", []item{
				it(itemLeftLink, "[["),
				it(itemLink, "amazing"),
				it(itemLinkDelim, "|"),
				it(itemText, "html "),
				it(itemOpenTagLeft, "<"),
				it(itemTagName, "span"),
				it(itemTagAttrName, "unquoted"),
				it(itemTagAttrValue, "blah"),
				it(itemTagAttrName, "class"),
				it(itemTagAttrValue, "test"),
				it(itemTagRight, ">"),
				it(itemText, "text"),
				it(itemCloseTagLeft, "</"),
				it(itemTagName, "span"),
				it(itemCloseTagRight, ">"),
				it(itemRightLink, "]]"),
			}, false},
			{"[[amazing|html <span class=\"test unquoted>text</span]]", "Unbalanced HTML quoted and unquoted attributes in balanced link", []item{
				it(itemLeftLink, "[["),
				it(itemLink, "amazing"),
				it(itemLinkDelim, "|"),
				it(itemText, "html "),
				it(itemOpenTagLeft, "<"),
				it(itemTagName, "span"),
				it(itemTagAttrName, "unquoted"),
				it(itemTagAttrValue, "blah"),
				it(itemTagAttrName, "class"),
				it(itemTagAttrValue, "test unquoted"),
				it(itemTagRight, ">"),
				it(itemText, "text"),
				it(itemCloseTagLeft, "</"),
				it(itemTagName, "span"),
				it(itemCloseTagRight, ">"),
				it(itemRightLink, "]]"),
			}, false},

			// Balance tests - HTML blank attributes

			{"[[amazing|html <span class=>text</span>]]", "Missing attribute value", []item{
				it(itemLeftLink, "[["),
				it(itemLink, "amazing"),
				it(itemLinkDelim, "|"),
				it(itemText, "html "),
				it(itemOpenTagLeft, "<"),
				it(itemTagName, "span"),
				it(itemTagAttrName, "class"),
				it(itemTagRight, ">"),
				it(itemText, "text"),
				it(itemCloseTagLeft, "</"),
				it(itemTagName, "span"),
				it(itemCloseTagRight, ">"),
				it(itemRightLink, "]]"),
			}, false},
			{"[[amazing|html <span class= style=\"color: red\">text</span>", "Initial missing attribute value", []item{
				it(itemLeftLink, "[["),
				it(itemLink, "amazing"),
				it(itemLinkDelim, "|"),
				it(itemText, "html "),
				it(itemOpenTagLeft, "<"),
				it(itemTagName, "span"),
				it(itemTagAttrName, "class"),
				it(itemTagAttrName, "style"),
				it(itemTagAttrValue, "color: red"),
				it(itemTagRight, ">"),
				it(itemText, "text"),
				it(itemCloseTagLeft, "</"),
				it(itemTagName, "span"),
				it(itemCloseTagRight, ">"),
				it(itemRightLink, "]]"),
			}, false},
			{"[[amazing|html <span class=\"test\" style=>text</span>", "Trailing missing attribute value", []item{
				it(itemLeftLink, "[["),
				it(itemLink, "amazing"),
				it(itemLinkDelim, "|"),
				it(itemText, "html "),
				it(itemOpenTagLeft, "<"),
				it(itemTagName, "span"),
				it(itemTagAttrName, "class"),
				it(itemTagAttrValue, "test"),
				it(itemTagAttrName, "style"),
				it(itemTagRight, ">"),
				it(itemText, "text"),
				it(itemCloseTagLeft, "</"),
				it(itemTagName, "span"),
				it(itemCloseTagRight, ">"),
				it(itemRightLink, "]]"),
			}, false},
		*/

		// TODO: Add void tag (self-closing and not) balance tests. With quoted / non-quoted attributes.
		/*
			// Balance tests -- HTML void tags
				{"<ref name=\"SOED\" />", "Balanced HTML quoted void tags", []item{
					it(itemOpenTagLeft, "<"),
					it(itemTagName, "ref"),
					it(itemTagAttrName, "name"),
					it(itemTagAttrValue, "SOED"),
					it(itemCloseTagRight, "/>"),
				}, false},
				{"<ref name=SOED />", "HTML void tag unquoted attribute whitespace", []item{
					it(itemOpenTagLeft, "<"),
					it(itemTagName, "ref"),
					it(itemTagAttrName, "name"),
					it(itemTagAttrValue, "SOED"),
					it(itemCloseTagRight, "/>"),
				}, false},
		*/
		// TODO: Matching / non-matching HTML tag names.

		// Balance tests - HTML comments and links

		/*
			{"[[amazing|html <!-- comment -->]]", "Balanced HTML comment in balanced link", []item{
				it(itemLeftLink, "[["),
				it(itemLink, "amazing"),
				it(itemLinkDelim, "|"),
				it(itemText, "html "),
				it(itemTagCommentLeft, "<!--"),
				it(itemTagComment, " comment "),
				it(itemTagCommentRight, "-->"),
				it(itemRightLink, "]]"),
			}, false},
			{"[[amazing|html <!-- comment -->", "Balanced HTML comment in unbalanced link", []item{
				it(itemText, "[[amazing|html "),
				it(itemTagCommentLeft, "<!--"),
				it(itemTagComment, " comment "),
				it(itemTagCommentRight, "-->"),
			}, false},
			{"[[amazing|html <!-- comment]]", "Unbalanced HTML comment in balanced link", []item{
				it(itemLeftLink, "[["),
				it(itemLink, "amazing"),
				it(itemLinkDelim, "|"),
				it(itemText, "html <!-- comment"),
				it(itemRightLink, "]]"),
			}, false},
			{"[[amazing|html <!-- comment", "Unbalanced HTML comment in unbalanced link", []item{
				it(itemText, "[[amazing|html <!-- comment"),
			}, false},
		*/
	}

	for _, tt := range tests {
		// TODO: Remove
		if len(tt.want) == 0 {
			continue
		}

		l := NewLexer2(tt.input, true)
		l.debug = tt.debug

		want := l.String(tt.want)
		got := l.String(l.Items())

		if want != got {
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(want, got, false)
			diffText := dmp.DiffPrettyText(diffs)
			t.Errorf("%s: NewLexer(%q) diff:\n%s", tt.desc, tt.input, diffText)
		}
	}
}

func it(typ itemType, val string) item {
	return item{typ: typ, val: val}
}

func ih(typ itemType, val string, depth int) item {
	return item{typ: typ, val: val, depth: depth}
}
