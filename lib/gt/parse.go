package gt

import (
	"fmt"
	"strings"

	"github.com/vthommeret/glossterm/lib/lang"
	"github.com/vthommeret/glossterm/lib/tpl"
)

type Word struct {
	Name      string
	Languages []Language
}

type Language struct {
	Code            string
	Etymology       Etymology
	Descendants     []tpl.Link
	DescendantTrees []tpl.EtymTree
}

type Etymology struct {
	Mentions  []tpl.Mention
	Borrows   []tpl.Borrow
	Derived   []tpl.Derived
	Inherited []tpl.Inherited
	Prefixes  []tpl.Prefix
	Suffixes  []tpl.Suffix
}

func (w *Word) IsEmpty() bool {
	return w.Languages == nil
}

func (l *Language) IsEmpty() bool {
	if l.Descendants != nil {
		return false
	}
	if l.DescendantTrees != nil {
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
	if l.Etymology.Inherited != nil {
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

type ListItem struct {
	Prefix string
	Links  []string
}

func Parse(p Page, langMap map[string]bool) (Word, error) {
	return ParseWord(p.Title, p.Text, langMap)
}

func ParseWord(name, text string, langMap map[string]bool) (Word, error) {
	w := Word{
		Name: name,
	}

	var inLanguageHeader bool
	var inSectionHeader bool

	var section sectionType
	var subSection sectionType
	var sectionDepth int = -1

	var language *Language
	var template *tpl.Template
	var namedParam *tpl.Parameter
	var listItem *ListItem

	var depth int

	l := NewLexer(text)

Parse:
	for {
		i := l.NextItem()

		switch i.typ {
		case itemError:
			return Word{}, fmt.Errorf("unable to parse: %s", i.val)
		case itemEOF:
			if language != nil {
				if listItem != nil {
					language.Descendants =
						append(language.Descendants, listItem.TplLinks(langMap)...)
				}
				if !language.IsEmpty() {
					w.Languages = append(w.Languages, *language)
				}
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
			if language != nil && listItem != nil {
				language.Descendants =
					append(language.Descendants, listItem.TplLinks(langMap)...)
			}
			listItem = nil
		case itemHeaderEnd:
			if i.depth == 2 {
				inLanguageHeader = false
			} else if i.depth > 2 {
				inSectionHeader = false
			}
		case itemListItemStart:
			if subSection == descendantsSection {
				listItem = &ListItem{}
			}
		case itemListItemPrefix:
			if listItem != nil {
				listItem.Prefix = i.val
			}
		case itemLink:
			if listItem != nil && !strings.Contains(i.val, ":") {
				listItem.Links = append(listItem.Links, i.val)
			}
		case itemText:
			// TODO: More intelligently handle whitespace in lexer. Emit newline
			// tokens and ignore otherwise unimportant whitespace.
			if language != nil && listItem != nil && strings.Contains(i.val, "\n") {
				language.Descendants =
					append(language.Descendants, listItem.TplLinks(langMap)...)
				listItem = nil
			}
			if inLanguageHeader {
				if l, ok := lang.CanonicalLangs[i.val]; ok {
					if _, ok := langMap[l.Code]; ok {
						language.Code = l.Code
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
					// This will exclude subsections named "Etymology" for now, e.g. https://en.wiktionary.org/wiki/taco#Noun_4
					section = unknownSection

					if sectionDepth >= 3 && i.val == "Descendants" {
						subSection = descendantsSection
					} else {
						subSection = unknownSection
					}
				}
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
					case "inh", "inherited":
						inherited := template.ToInherited()
						if _, ok := langMap[inherited.Lang]; ok {
							if _, ok := langMap[inherited.FromLang]; ok {
								language.Etymology.Inherited =
									append(language.Etymology.Inherited, inherited)
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
				case "etymtree":
					etymTree := template.ToEtymTree()
					if _, ok := langMap[etymTree.Lang]; ok {
						if _, ok := langMap[etymTree.RootLang]; ok || etymTree.RootLang == "" {
							if etymTree.Word == "" {
								etymTree.Word = w.Name
							}
							language.DescendantTrees =
								append(language.DescendantTrees, etymTree)
						}
					}
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

	return w, nil
}

func (li *ListItem) TplLinks(langMap map[string]bool) (ls []tpl.Link) {
	for _, l := range li.Links {
		var c string
		if l, ok := lang.CanonicalLangs[li.Prefix]; ok {
			if _, ok := langMap[l.Code]; ok {
				c = l.Code
			} else {
				continue
			}
		} else {
			continue
		}
		ls = append(ls, tpl.Link{Lang: c, Word: l})
	}
	return ls
}
