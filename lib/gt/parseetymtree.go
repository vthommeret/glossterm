package gt

import (
	"fmt"
	"strings"
	"vthommeret/glossterm/lib/tpl"
)

type Descendants struct {
	Word  string
	Links []tpl.Link
}

// Parses a given etymology tree (e.g. https://en.wiktionary.org/wiki/Template:etymtree/la/germanus)
func ParseEtymTree(p Page, langMap map[string]bool) (*Descendants, error) {
	descendants := Descendants{
		Word: p.Title,
	}

	var template *tpl.Template
	var namedParam *tpl.Parameter
	var listItem *ListItem

	var depth int

	l := NewLexer(p.Text)

Parse:
	for {
		i := l.NextItem()

		switch i.typ {
		case itemError:
			return nil, fmt.Errorf("unable to parse: %s", i.val)
		case itemEOF:
			if listItem != nil {
				descendants.Links =
					append(descendants.Links, listItem.TplLinks(langMap)...)
			}
			break Parse
		case itemListItemStart:
			listItem = &ListItem{}
		case itemListItemPrefix:
			if listItem != nil {
				listItem.Prefix = i.val
			}
		case itemLink:
			if listItem != nil {
				listItem.Links = append(listItem.Links, i.val)
			}
		case itemText:
			// TODO: More intelligently handle whitespace in lexer. Emit newline
			// tokens and ignore otherwise unimportant whitespace.
			if listItem != nil && strings.Contains(i.val, "\n") {
				descendants.Links =
					append(descendants.Links, listItem.TplLinks(langMap)...)
				listItem = nil
			}
		case itemLeftTemplate:
			depth++
			if depth == 1 {
				template = &tpl.Template{}
			} else {
				// Nested templates aren't supported for now.
				template = nil
			}
		case itemRightTemplate:
			depth--
			if template == nil {
				break
			}
			switch template.Action {
			case "l", "link":
				link := template.ToLink()
				if _, ok := langMap[link.Lang]; ok {
					descendants.Links = append(descendants.Links, link)
				}
			}
		case itemAction:
			if template == nil {
				break
			}
			template.Action = i.val
		case itemParamText:
			if template == nil {
				break
			}
			if namedParam != nil {
				namedParam.Value = i.val
				template.NamedParameters = append(template.NamedParameters,
					*namedParam)
				namedParam = nil
			} else {
				template.Parameters = append(template.Parameters, i.val)
			}
		case itemParamName:
			if template == nil {
				break
			}
			namedParam = &tpl.Parameter{Name: i.val}
		}
	}

	return &descendants, nil
}
