package ml

import (
	"fmt"
	"strings"
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
	Templates []Template
}

type Template struct {
	Action     string
	Parameters []Parameter
}

type Parameter struct {
	Name  string
	Value string
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
	var tpl *Template
	var param *Parameter

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
				language = &Language{
					Etymology: Etymology{
						Templates: []Template{},
					},
				}
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
				tpl = &Template{Parameters: []Parameter{}}
			}
		case itemRightTemplate:
			if sectionType == etymologySection {
				if language != nil {
					language.Etymology.Templates = append(language.Etymology.Templates,
						*tpl,
					)
				}
			}
		case itemAction:
			if sectionType == etymologySection {
				tpl.Action = i.val
			}
		case itemParam:
			if sectionType == etymologySection {
				tpl.Parameters = append(tpl.Parameters, Parameter{Value: i.val})
			}
		case itemParamName:
			if sectionType == etymologySection {
				param = &Parameter{Name: i.val}
			}
		case itemParamValue:
			if sectionType == etymologySection {
				param.Value = i.val
				tpl.Parameters = append(tpl.Parameters, *param)
			}
		}
	}

	return w, nil
}
