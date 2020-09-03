package gt

import (
	"fmt"
	"strings"
	"time"

	"github.com/vthommeret/glossterm/lib/lang"
	"github.com/vthommeret/glossterm/lib/tpl"
)

type Word struct {
	Name      string               `json:"name"`
	Languages *map[string]Language `json:"languages"`
	Indexed   *time.Time           `json:"indexed,omitempty"`
}

type Language struct {
	Code            string           `firestore:"code"`
	Definitions     *Definitions     `json:"definitions,omitempty" firestore:"definitions,omitempty"`
	Etymology       *Etymology       `json:"etymology,omitempty" firestore:"etymology,omitempty"`
	Links           []tpl.Link       `json:"links,omitempty" firestore:"links,omitempty"`
	Descendants     []tpl.Descendant `json:"descendants,omitempty" firestore:"descendants,omitempty"`
	DescendantTrees []tpl.EtymTree   `json:"descendantTrees,omitempty" firestore:"descendantTrees,omitempty"`

	section      sectionType
	subSection   sectionType
	sectionDepth int

	listItem             *ListItem
	listItemDepth        int
	inListItemDefinition bool
	inListItemSublist    bool

	definitionBuffer TextBuffer
}

type Definitions struct {
	Nouns         []string `json:"nouns,omitempty" firestore:"nouns,omitempty"`
	Adjectives    []string `json:"adjectives,omitempty" firestore:"adjectives,omitempty"`
	Verbs         []string `json:"verbs,omitempty" firestore:"verbs,omitempty"`
	Interjections []string `json:"interjections,omitempty" firestore:"interjections,omitempty"`
	Numerals      []string `json:"numerals,omitempty" firestore:"numerals,omitempty"`
}

type Etymology struct {
	Cognates  []tpl.Cognate   `json:"cognates,omitempty" firestore:"cognates,omitempty"`
	Mentions  []tpl.Mention   `json:"mentions,omitempty" firestore:"mentions,omitempty"`
	Borrows   []tpl.Borrow    `json:"borrows,omitempty" firestore:"borrows,omitempty"`
	Derived   []tpl.Derived   `json:"derived,omitempty" firestore:"derived,omitempty"`
	Inherited []tpl.Inherited `json:"inherited,omitempty" firestore:"inherited,omitempty"`
	Prefixes  []tpl.Prefix    `json:"prefixes,omitempty" firestore:"prefixes,omitempty"`
	Suffixes  []tpl.Suffix    `json:"suffixes,omitempty" firestore:"suffixes,omitempty"`
}

func (w *Word) IsEmpty() bool {
	return w.Languages == nil
}

func (l *Language) IsEmpty() bool {
	if l.Definitions != nil {
		if l.Definitions.Nouns != nil {
			return false
		}
		if l.Definitions.Adjectives != nil {
			return false
		}
		if l.Definitions.Verbs != nil {
			return false
		}
		if l.Definitions.Interjections != nil {
			return false
		}
		if l.Definitions.Numerals != nil {
			return false
		}
	}
	if l.Etymology != nil {
		if l.Etymology.Cognates != nil {
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
	}
	if l.Links != nil {
		return false
	}
	if l.Descendants != nil {
		return false
	}
	if l.DescendantTrees != nil {
		return false
	}
	return true
}

func (l *Language) flushDefinition() {
	definition := strings.TrimSpace(strings.Join(l.definitionBuffer, ""))
	if definition != "" {
		switch l.section {
		case nounSection:
			if l.Definitions == nil {
				l.Definitions = &Definitions{}
			}
			l.Definitions.Nouns =
				append(l.Definitions.Nouns, definition)
		case adjectiveSection:
			if l.Definitions == nil {
				l.Definitions = &Definitions{}
			}
			l.Definitions.Adjectives =
				append(l.Definitions.Adjectives, definition)
		case verbSection:
			if l.Definitions == nil {
				l.Definitions = &Definitions{}
			}
			l.Definitions.Verbs =
				append(l.Definitions.Verbs, definition)
		case interjectionSection:
			if l.Definitions == nil {
				l.Definitions = &Definitions{}
			}
			l.Definitions.Interjections =
				append(l.Definitions.Interjections, definition)
		case numeralSection:
			if l.Definitions == nil {
				l.Definitions = &Definitions{}
			}
			l.Definitions.Numerals =
				append(l.Definitions.Numerals, definition)
		}
	}
	l.definitionBuffer = nil
	l.inListItemDefinition = false
	l.inListItemSublist = false
}

type sectionType int

const (
	unknownSection sectionType = iota
	etymologySection

	nounSection
	adjectiveSection
	verbSection
	interjectionSection
	numeralSection

	descendantsSection
)

const linkCategoryPrefix = "Category:"
const linkReferencePrefix = "#"

var wordTypeMap = map[string]sectionType{
	"Noun":         nounSection,
	"Adjective":    adjectiveSection,
	"Verb":         verbSection,
	"Interjection": interjectionSection,
	"Numeral":      numeralSection,
}

type ListItem struct {
	Prefix string
	Links  []string
}

type TextBuffer []string

func definitionSection(section sectionType) bool {
	return section == nounSection || section == adjectiveSection || section == verbSection || section == interjectionSection || section == numeralSection
}

// Parses a given word (e.g. https://en.wiktionary.org/wiki/hombre).
func ParseWord(p Page, langMap map[string]bool) (Word, error) {
	name := p.Title
	text := p.Text

	w := Word{
		Name: name,
	}

	var inLanguageHeader bool
	var inSectionHeader bool

	var language *Language
	var template *tpl.Template
	var namedParam *tpl.Parameter

	var templateDepth int

	l := NewLexer(text)

Parse:
	for {
		i := l.NextItem()

		switch i.typ {
		case itemError:
			return Word{}, fmt.Errorf("unable to parse: %s", i.val)
		case itemEOF:
			if language != nil {
				if language.listItem != nil {
					language.Links =
						append(language.Links, language.listItem.TplLinks(langMap)...)
				}
				if !language.IsEmpty() {
					if w.Languages == nil {
						w.Languages = &map[string]Language{}
					}
					(*w.Languages)[language.Code] = *language
				}
			}
			break Parse
		case itemHeaderStart:
			if i.depth == 1 {
				language = nil
				inLanguageHeader = false
				inSectionHeader = false
				if language != nil {
					language.section = unknownSection
					language.subSection = unknownSection
					language.sectionDepth = -1
				}
			} else if i.depth == 2 {
				if language != nil && !language.IsEmpty() {
					if w.Languages == nil {
						w.Languages = &map[string]Language{}
					}
					(*w.Languages)[language.Code] = *language
				}
				language = &Language{sectionDepth: -1}
				inLanguageHeader = true
			} else if i.depth > 2 {
				inSectionHeader = true
				if language != nil {
					language.sectionDepth = i.depth - 1
				}
			}
			if language != nil {
				if language.listItem != nil {
					language.Links =
						append(language.Links, language.listItem.TplLinks(langMap)...)
				}
			}
		case itemHeaderEnd:
			if i.depth == 2 {
				inLanguageHeader = false
			} else if i.depth > 2 {
				inSectionHeader = false
			}
		case itemUnorderedListItemStart:
			if language != nil {
				language.listItemDepth = i.depth
			}
		case itemOrderedListItemStart:
			if language != nil {
				language.listItemDepth = i.depth
			}
		case itemOrderedDefinitionStart:
			if language != nil {
				language.inListItemDefinition = true
			}
		case itemOrderedUnorderedStart:
			if language != nil {
				language.inListItemSublist = true
			}
		case itemUnorderedOrderedStart:
			if language != nil {
				language.inListItemSublist = true
			}
		case itemListItemPrefix:
			if language != nil && language.listItem != nil {
				language.listItem.Prefix = i.val
			}
		case itemListItemEnd:
			if language != nil && language.definitionBuffer != nil {
				language.flushDefinition()
			}
		case itemLink:
			if language != nil {
				if language.listItem != nil {
					language.listItem.Links = append(language.listItem.Links, i.val)
				} else if language.definitionBuffer != nil && language.listItemDepth == 1 && !language.inListItemDefinition && !language.inListItemSublist && !strings.HasPrefix(i.val, linkCategoryPrefix) && !strings.HasPrefix(i.val, linkReferencePrefix) {
					language.definitionBuffer = append(language.definitionBuffer, i.val)
				}
			}
		case itemText:
			// TODO: More intelligently handle whitespace in lexer. Emit newline
			// tokens and ignore otherwise unimportant whitespace.
			if language != nil && language.listItem != nil && strings.Contains(i.val, "\n") {
				language.Links =
					append(language.Links, language.listItem.TplLinks(langMap)...)
			}
			if inLanguageHeader {
				if l, ok := lang.CanonicalLangs[i.val]; ok {
					if _, ok := langMap[l.Code]; ok {
						language.Code = l.Code
						language.section = unknownSection
						language.subSection = unknownSection
					} else {
						language = nil
					}
				} else {
					language = nil
				}
			} else if inSectionHeader && language != nil {
				if language.sectionDepth == 2 && strings.HasPrefix(i.val, "Etymology") {
					language.section = etymologySection
				} else if wordSection, ok := wordTypeMap[i.val]; (language.sectionDepth == 2 || language.sectionDepth == 3) && ok {
					language.section = wordSection
				} else {
					// This will exclude subsections named "Etymology" for now, e.g. https://en.wiktionary.org/wiki/taco#Noun_4
					language.section = unknownSection

					if language.sectionDepth >= 3 && i.val == "Descendants" {
						language.subSection = descendantsSection
					} else {
						language.subSection = unknownSection
					}
				}
			} else if language != nil && definitionSection(language.section) && language.listItemDepth == 1 && !language.inListItemDefinition && !language.inListItemSublist {
				language.definitionBuffer = append(language.definitionBuffer, i.val)
			}
		case itemLeftTemplate:
			templateDepth++
			if templateDepth == 1 {
				template = &tpl.Template{}
			} else {
				// Nested templates aren't supported for now.
				template = nil
			}
		case itemRightTemplate:
			templateDepth--
			if language == nil || template == nil {
				break
			}
			if language.section == etymologySection {
				switch template.Action {
				case "cog", "cognate":
					cognate := template.ToCognate()
					if _, ok := langMap[cognate.Lang]; ok {
						if language.Etymology == nil {
							language.Etymology = &Etymology{}
						}
						language.Etymology.Cognates =
							append(language.Etymology.Cognates, cognate)
					}
				case "m", "mention":
					mention := template.ToMention()
					if _, ok := langMap[mention.Lang]; ok {
						if language.Etymology == nil {
							language.Etymology = &Etymology{}
						}
						language.Etymology.Mentions =
							append(language.Etymology.Mentions, mention)
					}
				case "bor", "borrowing":
					borrow := template.ToBorrow()
					if _, ok := langMap[borrow.Lang]; ok {
						if _, ok := langMap[borrow.FromLang]; ok {
							if language.Etymology == nil {
								language.Etymology = &Etymology{}
							}
							language.Etymology.Borrows =
								append(language.Etymology.Borrows, borrow)
						}
					}
				case "der", "derived":
					derived := template.ToDerived()
					if _, ok := langMap[derived.Lang]; ok {
						if _, ok := langMap[derived.FromLang]; ok {
							if language.Etymology == nil {
								language.Etymology = &Etymology{}
							}
							language.Etymology.Derived =
								append(language.Etymology.Derived, derived)
						}
					}
				case "inh", "inherited":
					inherited := template.ToInherited()
					if _, ok := langMap[inherited.Lang]; ok {
						if _, ok := langMap[inherited.FromLang]; ok {
							if language.Etymology == nil {
								language.Etymology = &Etymology{}
							}
							language.Etymology.Inherited =
								append(language.Etymology.Inherited, inherited)
						}
					}
				case "prefix":
					prefix := template.ToPrefix()
					if _, ok := langMap[prefix.Lang]; ok {
						if language.Etymology == nil {
							language.Etymology = &Etymology{}
						}
						language.Etymology.Prefixes =
							append(language.Etymology.Prefixes, prefix)
					}
				case "suffix":
					suffix := template.ToSuffix()
					if _, ok := langMap[suffix.Lang]; ok {
						if language.Etymology == nil {
							language.Etymology = &Etymology{}
						}
						language.Etymology.Suffixes =
							append(language.Etymology.Suffixes, suffix)
					}
				}
			} else if language.subSection == descendantsSection {
				switch template.Action {
				case "desc", "descendant":
					desc := template.ToDescendant()
					if _, ok := langMap[desc.Lang]; ok {
						language.Descendants =
							append(language.Descendants, desc)
					}
				case "l", "link":
					link := template.ToLink()
					if _, ok := langMap[link.Lang]; ok {
						language.Links =
							append(language.Links, link)
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
			if definitionSection(language.section) {
				switch template.Action {
				case "l", "link":
					link := template.ToLink()
					if language.definitionBuffer != nil {
						language.definitionBuffer = append(language.definitionBuffer, link.Text())
					}
				case "gloss":
					gloss := template.ToGloss()
					if language.definitionBuffer != nil {
						language.definitionBuffer = append(language.definitionBuffer, gloss.Text())
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
	for _, link := range li.Links {
		if strings.Contains(link, ":") {
			continue
		}
		parts := strings.Split(link, "#")
		link = parts[0]
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
		ls = append(ls, tpl.Link{Lang: c, Word: link})
	}
	return ls
}
