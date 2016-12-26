package gt

import (
	"fmt"
	"strings"

	"github.com/vthommeret/glossterm/lib/tpl"
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
	Derived  []tpl.Derived
	Prefixes []tpl.Prefix
	Suffixes []tpl.Suffix
}

func (w *Word) IsEmpty() bool {
	return w.Languages == nil
}

func (l *Language) IsEmpty() bool {
	if l.Descendants != nil {
		return false
	}
	if l.Etymology.Mentions != nil {
		return false
	}
	if l.Etymology.Borrows != nil {
		return false
	}
	if l.Etymology.Derived != nil {
		return false
	}
	if l.Etymology.Prefixes != nil {
		return false
	}
	if l.Etymology.Suffixes != nil {
		return false
	}
	return true
}

type sectionType int

const (
	unknownSection sectionType = iota
	etymologySection
	descendantsSection
)

func Parse(p Page, langMap map[string]bool) (Word, error) {
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
			if language != nil && !language.IsEmpty() {
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
				if language != nil && !language.IsEmpty() {
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
					if _, ok := langMap[lang.Code]; ok {
						language.Code = lang.Code
					} else {
						language = nil
					}
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
						mention := template.ToMention()
						if _, ok := langMap[mention.Lang]; ok {
							language.Etymology.Mentions =
								append(language.Etymology.Mentions, mention)
						}
					case "bor", "borrowing":
						borrow := template.ToBorrow()
						if _, ok := langMap[borrow.Lang]; ok {
							if _, ok := langMap[borrow.FromLang]; ok {
								language.Etymology.Borrows =
									append(language.Etymology.Borrows, borrow)
							}
						}
					case "der", "derived":
						derived := template.ToDerived()
						if _, ok := langMap[derived.Lang]; ok {
							if _, ok := langMap[derived.FromLang]; ok {
								language.Etymology.Derived =
									append(language.Etymology.Derived, derived)
							}
						}
					case "prefix":
						prefix := template.ToPrefix()
						if _, ok := langMap[prefix.Lang]; ok {
							language.Etymology.Prefixes =
								append(language.Etymology.Prefixes, prefix)
						}
					case "suffix":
						suffix := template.ToSuffix()
						if _, ok := langMap[suffix.Lang]; ok {
							language.Etymology.Suffixes =
								append(language.Etymology.Suffixes, suffix)
						}
					}
				}
			} else if subSection == descendantsSection {
				switch template.Action {
				case "l", "link":
					if language != nil {
						link := template.ToLink()
						if _, ok := langMap[link.Lang]; ok {
							language.Descendants =
								append(language.Descendants, link)
						}
					}
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