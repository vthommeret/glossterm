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
	Code        string
	Etymology   Etymology
	Descendants []tpl.Link
}

type Etymology struct {
	Mentions []tpl.Mention
	Borrows  []tpl.Borrow
	Prefixes []tpl.Prefix
	Suffixes []tpl.Suffix
}

type Descendant struct {
	Language string
	Word     string
}

type sectionType int

const (
	unknownSection sectionType = iota
	etymologySection
	descendantsSection
)

func Parse(p Page) (Word, error) {
	w := Word{
		Name: p.Title,
	}

	var inLanguageHeader bool
	var inSectionHeader bool

	var section sectionType
	var subSection sectionType
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
				section = unknownSection
				subSection = unknownSection
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
				if lang, ok := CanonicalLangs[i.val]; ok {
					language.Code = lang.Code
				} else {
					language = nil
				}
				section = unknownSection
				subSection = unknownSection
			} else if inSectionHeader {
				if sectionDepth == 2 {
					if strings.HasPrefix(i.val, "Etymology") {
						section = etymologySection
					} else {
						section = unknownSection
					}
					subSection = unknownSection
				} else {
					// This will exclude subsections of "Etymology" for now, e.g. https://en.wiktionary.org/wiki/taco#Noun_4
					section = unknownSection

					if sectionDepth == 3 && i.val == "Descendants" {
						subSection = descendantsSection
					} else {
						subSection = unknownSection
					}
				}
			}
		case itemLeftTemplate:
			template = &tpl.Template{}
		case itemRightTemplate:
			if section == etymologySection {
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
			} else if subSection == descendantsSection {
				switch template.Action {
				case "l", "link":
					language.Descendants = append(language.Descendants,
						template.ToLink(),
					)
				}
			}
		case itemAction:
			template.Action = i.val
		case itemParam:
			template.Parameters = append(template.Parameters, i.val)
		case itemParamName:
			param = &tpl.Parameter{Name: i.val}
		case itemParamValue:
			param.Value = i.val
			template.NamedParameters = append(template.NamedParameters, *param)
		}
	}

	return w, nil
}

// FilterLangs filters specific languages.
func (w *Word) FilterLangs(filters []string) {
	langMap := make(map[string]bool)
	for _, l := range filters {
		langMap[l] = true
	}
	var langs []Language
	for _, l := range w.Languages {
		if _, ok := langMap[l.Code]; ok {
			var descendants []tpl.Link
			for _, d := range l.Descendants {
				if _, ok := langMap[d.Lang]; ok {
					descendants = append(descendants, d)
				}
			}
			l.Descendants = descendants
			langs = append(langs, l)
		}
	}
	w.Languages = langs

}
