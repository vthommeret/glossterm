package ml

import (
	"fmt"
	"strings"

	"github.com/vthommeret/memory.limited/lib/tpl"
)

type Word struct {
	Name      string
	Languages []Language
}

type Language struct {
	Name      string
	Etymology Etymology
}

type Etymology struct {
	Mentions []tpl.Mention
	Borrows  []tpl.Borrow
	Prefixes []tpl.Prefix
	Suffixes []tpl.Suffix
}

type sectionType int

const (
	unknownSection sectionType = iota
	etymologySection
)

func Parse(p Page) (Word, error) {
	w := Word{
		Name: p.Title,
	}

	var inLanguageHeader bool
	var inSectionHeader bool

	var sectionType sectionType
	var sectionDepth int = -1

	var language *Language
	var template *tpl.Template
	var param *tpl.Parameter

	l := NewLexer(p.Text)

Parse:
	for {
		i := l.NextItem()
		switch i.typ {
		case itemError:
			return Word{}, fmt.Errorf("unable to parse: %s", i.val)
		case itemEOF:
			if language != nil {
				w.Languages = append(w.Languages, *language)
			}
			break Parse
		case itemHeaderStart:
			if i.depth == 1 {
				language = nil
				inLanguageHeader = false
				inSectionHeader = false
				sectionType = unknownSection
				sectionDepth = -1
			} else if i.depth == 2 {
				if language != nil {
					w.Languages = append(w.Languages, *language)
				}
				language = &Language{}
				inLanguageHeader = true
			} else if i.depth > 2 {
				inSectionHeader = true
				sectionDepth = i.depth - 1
			}
		case itemHeaderEnd:
			if i.depth == 2 {
				inLanguageHeader = false
			} else if i.depth > 2 {
				inSectionHeader = false
			}
		case itemText:
			if inLanguageHeader {
				language.Name = i.val
				sectionType = unknownSection
			} else if inSectionHeader {
				if sectionDepth == 2 {
					if strings.HasPrefix(i.val, "Etymology") {
						sectionType = etymologySection
					} else {
						sectionType = unknownSection
					}
				} else {
					// This will exclude subsections of "Etymology" for now, e.g. https://en.wiktionary.org/wiki/taco#Noun_4
					sectionType = unknownSection
				}
			}
		case itemLeftTemplate:
			if sectionType == etymologySection {
				template = &tpl.Template{}
			}
		case itemRightTemplate:
			if sectionType == etymologySection {
				if language != nil {
					switch template.Action {
					case "m", "mention":
						language.Etymology.Mentions = append(language.Etymology.Mentions,
							template.ToMention(),
						)
					case "bor", "borrowing":
						language.Etymology.Borrows = append(language.Etymology.Borrows,
							template.ToBorrow(),
						)
					case "prefix":
						language.Etymology.Prefixes = append(language.Etymology.Prefixes,
							template.ToPrefix(),
						)
					case "suffix":
						language.Etymology.Suffixes = append(language.Etymology.Suffixes,
							template.ToSuffix(),
						)
					}
				}
			}
		case itemAction:
			if sectionType == etymologySection {
				template.Action = i.val
			}
		case itemParam:
			if sectionType == etymologySection {
				template.Parameters = append(template.Parameters, i.val)
			}
		case itemParamName:
			if sectionType == etymologySection {
				param = &tpl.Parameter{Name: i.val}
			}
		case itemParamValue:
			if sectionType == etymologySection {
				param.Value = i.val
				template.NamedParameters = append(template.NamedParameters, *param)
			}
		}
	}

	return w, nil
}
