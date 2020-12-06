package gt

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/vthommeret/glossterm/lib/lang"
	"github.com/vthommeret/glossterm/lib/tpl"
)

type Word struct {
	Name      string               `json:"name"`
	Languages map[string]*Language `json:"languages"`
	Indexed   *time.Time           `json:"indexed,omitempty"`
}

type Language struct {
	Code        string       `firestore:"code"`
	Definitions *Definitions `json:"definitions,omitempty" firestore:"definitions,omitempty"`
	Etymology   *Etymology   `json:"etymology,omitempty" firestore:"etymology,omitempty"`
	// TODO: Clean up Links / represent as descendants?
	// Partly used by legacy etymtree templates. Links are only used for descendants.
	Links       []tpl.Link       `json:"links,omitempty" firestore:"links,omitempty"`
	Descendants []tpl.Descendant `json:"descendants,omitempty" firestore:"descendants,omitempty"`
	// TODO: Rename EtymTrees. May need to re-write DB.
	DescendantTrees []tpl.EtymTree `json:"descendantTrees,omitempty" firestore:"descendantTrees,omitempty"`
	DescTrees       []tpl.DescTree `json:"descTrees,omitempty" firestore:"descTrees,omitempty"`
	Cognates        []*Cognate     `json:"cognates,omitempty" firestore:"cognates,omitempty"`

	section      sectionType
	subSection   sectionType
	sectionDepth int

	descendantLang *string

	listItem             *ListItem
	listItemDepth        int
	inListItemDefinition bool
	inListItemSublist    bool

	linkBuffer *LinkBuffer

	definitionBuffer TextBuffer
	definitionRoot   *RootWord
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

type LinkBuffer struct {
	Link string
	Name *string
}

type Definitions struct {
	Nouns         []Definition `json:"nouns,omitempty" firestore:"nouns,omitempty"`
	Adjectives    []Definition `json:"adjectives,omitempty" firestore:"adjectives,omitempty"`
	Verbs         []Definition `json:"verbs,omitempty" firestore:"verbs,omitempty"`
	Adverbs       []Definition `json:"adverbs,omitempty" firestore:"adverbs,omitempty"`
	Articles      []Definition `json:"articles,omitempty" firestore:"articles,omitempty"`
	Prepositions  []Definition `json:"prepositions,omitempty" firestore:"prepositions,omitempty"`
	Pronouns      []Definition `json:"pronouns,omitempty" firestore:"pronouns,omitempty"`
	Conjunctions  []Definition `json:"conjunctions,omitempty" firestore:"conjunctions,omitempty"`
	Interjections []Definition `json:"interjections,omitempty" firestore:"interjections,omitempty"`
	Numerals      []Definition `json:"numerals,omitempty" firestore:"numerals,omitempty"`
	Particles     []Definition `json:"particles,omitempty" firestore:"particles,omitempty"`
	Determiners   []Definition `json:"determiners,omitempty" firestore:"determiners,omitempty"`
}

type Definition struct {
	Text string    `json:"text" firestore:"text"`
	Root *RootWord `json:"root,omitempty" firestore:"root,omitempty"`
}

type RootWord struct {
	Lang string `json:"lang" firestore:"lang"`
	Name string `json:"name" firestore:"name"`
}

func (l *Language) AllDefinitions() [][]Definition {
	return [][]Definition{
		l.Definitions.Nouns,
		l.Definitions.Adjectives,
		l.Definitions.Verbs,
		l.Definitions.Adverbs,
		l.Definitions.Articles,
		l.Definitions.Prepositions,
		l.Definitions.Pronouns,
		l.Definitions.Conjunctions,
		l.Definitions.Interjections,
		l.Definitions.Numerals,
		l.Definitions.Particles,
		l.Definitions.Determiners,
	}
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
		if l.Definitions.Adverbs != nil {
			return false
		}
		if l.Definitions.Articles != nil {
			return false
		}
		if l.Definitions.Prepositions != nil {
			return false
		}
		if l.Definitions.Pronouns != nil {
			return false
		}
		if l.Definitions.Conjunctions != nil {
			return false
		}
		if l.Definitions.Interjections != nil {
			return false
		}
		if l.Definitions.Numerals != nil {
			return false
		}
		if l.Definitions.Particles != nil {
			return false
		}
		if l.Definitions.Determiners != nil {
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
	if l.definitionBuffer != nil {
		definition := strings.TrimSpace(strings.Join(l.definitionBuffer, ""))

		root := l.definitionRoot
		if definition != "" {
			switch l.section {
			case nounSection:
				if l.Definitions == nil {
					l.Definitions = &Definitions{}
				}
				l.Definitions.Nouns =
					append(l.Definitions.Nouns, Definition{Text: definition, Root: root})
			case adjectiveSection:
				if l.Definitions == nil {
					l.Definitions = &Definitions{}
				}
				l.Definitions.Adjectives =
					append(l.Definitions.Adjectives, Definition{Text: definition, Root: root})
			case verbSection:
				if l.Definitions == nil {
					l.Definitions = &Definitions{}
				}
				l.Definitions.Verbs =
					append(l.Definitions.Verbs, Definition{Text: definition, Root: root})
			case adverbSection:
				if l.Definitions == nil {
					l.Definitions = &Definitions{}
				}
				l.Definitions.Adverbs =
					append(l.Definitions.Adverbs, Definition{Text: definition, Root: root})
			case articleSection:
				if l.Definitions == nil {
					l.Definitions = &Definitions{}
				}
				l.Definitions.Articles =
					append(l.Definitions.Articles, Definition{Text: definition, Root: root})
			case prepositionSection:
				if l.Definitions == nil {
					l.Definitions = &Definitions{}
				}
				l.Definitions.Prepositions =
					append(l.Definitions.Prepositions, Definition{Text: definition, Root: root})
			case pronounSection:
				if l.Definitions == nil {
					l.Definitions = &Definitions{}
				}
				l.Definitions.Pronouns =
					append(l.Definitions.Pronouns, Definition{Text: definition, Root: root})
			case conjunctionSection:
				if l.Definitions == nil {
					l.Definitions = &Definitions{}
				}
				l.Definitions.Conjunctions =
					append(l.Definitions.Conjunctions, Definition{Text: definition, Root: root})
			case interjectionSection:
				if l.Definitions == nil {
					l.Definitions = &Definitions{}
				}
				l.Definitions.Interjections =
					append(l.Definitions.Interjections, Definition{Text: definition, Root: root})
			case numeralSection:
				if l.Definitions == nil {
					l.Definitions = &Definitions{}
				}
				l.Definitions.Numerals =
					append(l.Definitions.Numerals, Definition{Text: definition, Root: root})
			case particleSection:
				if l.Definitions == nil {
					l.Definitions = &Definitions{}
				}
				l.Definitions.Particles =
					append(l.Definitions.Particles, Definition{Text: definition, Root: root})
			case determinerSection:
				if l.Definitions == nil {
					l.Definitions = &Definitions{}
				}
				l.Definitions.Determiners =
					append(l.Definitions.Determiners, Definition{Text: definition, Root: root})
			}
		}
	}

	l.listItemDepth = 0
	l.inListItemDefinition = false
	l.inListItemSublist = false

	l.definitionBuffer = nil
	l.definitionRoot = nil
}

func (l *Language) shouldDefineLink() bool {
	return l.definitionBuffer != nil && l.listItemDepth == 1 && !l.inListItemDefinition && !l.inListItemSublist
}

type sectionType int

const (
	unknownSection sectionType = iota
	etymologySection

	nounSection
	adjectiveSection
	verbSection
	adverbSection
	articleSection
	prepositionSection
	pronounSection
	conjunctionSection
	interjectionSection
	numeralSection
	particleSection
	determinerSection

	descendantsSection
)

const linkCategoryPrefix = "Category:"
const linkReferencePrefix = "#"

const spanishLang = "es"

var wordTypeRegex *regexp.Regexp

var wordTypeMap = map[string]sectionType{
	"Noun":         nounSection,
	"Adjective":    adjectiveSection,
	"Verb":         verbSection,
	"Adverb":       adverbSection,
	"Article":      articleSection,
	"Preposition":  prepositionSection,
	"Pronoun":      pronounSection,
	"Conjunction":  conjunctionSection,
	"Interjection": interjectionSection,
	"Numeral":      numeralSection,
	"Particle":     particleSection,
	"Determiner":   determinerSection,
}

var formOfMap = map[string]*FormTemplate{
	"abbreviation of":                      &FormTemplate{Shortcuts: []string{"abbr of"}},
	"abstract noun of":                     &FormTemplate{},
	"accusative of":                        &FormTemplate{Tags: []string{"acc"}},
	"accusative plural of":                 &FormTemplate{Tags: []string{"acc", "p"}},
	"accusative singular of":               &FormTemplate{Tags: []string{"acc", "s"}},
	"acronym of":                           &FormTemplate{},
	"active participle of":                 &FormTemplate{Tags: []string{"act", "part"}},
	"adj form of":                          &FormTemplate{},
	"agent noun of":                        &FormTemplate{},
	"alternative case form of":             &FormTemplate{Shortcuts: []string{"alt case"}},
	"alternative form of":                  &FormTemplate{Shortcuts: []string{"alt form", "altform"}},
	"alternative plural of":                &FormTemplate{},
	"alternative reconstruction of":        &FormTemplate{Tags: []string{"alternative", "reconstruction"}},
	"alternative spelling of":              &FormTemplate{Shortcuts: []string{"alt sp"}},
	"alternative typography of":            &FormTemplate{},
	"aphetic form of":                      &FormTemplate{Tags: []string{"attr", "form"}},
	"apocopic form of":                     &FormTemplate{},
	"archaic form of":                      &FormTemplate{},
	"archaic inflection of":                &FormTemplate{},
	"archaic spelling of":                  &FormTemplate{},
	"aspirate mutation of":                 &FormTemplate{},
	"attributive form of":                  &FormTemplate{Tags: []string{"attr", "form"}},
	"augmentative of":                      &FormTemplate{},
	"broad form of":                        &FormTemplate{},
	"causative of":                         &FormTemplate{Tags: []string{"caus"}},
	"clipping of":                          &FormTemplate{Shortcuts: []string{"clip of"}},
	"combining form of":                    &FormTemplate{},
	"comparative of":                       &FormTemplate{},
	"construed with":                       &FormTemplate{},
	"contraction of":                       &FormTemplate{},
	"dated form of":                        &FormTemplate{},
	"dated spelling of":                    &FormTemplate{},
	"dative of":                            &FormTemplate{Tags: []string{"dat"}},
	"dative plural of":                     &FormTemplate{Tags: []string{"dat", "p"}},
	"dative singular of":                   &FormTemplate{Tags: []string{"dat", "s"}},
	"definite singular of":                 &FormTemplate{Tags: []string{"def", "s"}},
	"definite plural of":                   &FormTemplate{Tags: []string{"def", "p"}},
	"deliberate misspelling of":            &FormTemplate{},
	"diminutive of":                        &FormTemplate{Shortcuts: []string{"dim of"}},
	"dual of":                              &FormTemplate{Tags: []string{"d"}},
	"eclipsis of":                          &FormTemplate{Text: "eclipsed form of"},
	"eggcorn of":                           &FormTemplate{},
	"elative of":                           &FormTemplate{Tags: []string{"elad"}},
	"ellipsis of":                          &FormTemplate{},
	"elongated form of":                    &FormTemplate{},
	"endearing diminutive of":              &FormTemplate{},
	"endearing form of":                    &FormTemplate{},
	"equative of":                          &FormTemplate{Text: "equative degree of"},
	"euphemistic form of":                  &FormTemplate{},
	"euphemistic spelling of":              &FormTemplate{},
	"eye dialect of":                       &FormTemplate{Text: "eye dialect spelling of"},
	"female equivalent of":                 &FormTemplate{Tags: []string{"female", "equivalent"}},
	"feminine of":                          &FormTemplate{Tags: []string{"f"}},
	"feminine plural of":                   &FormTemplate{Tags: []string{"f", "p"}},
	"feminine plural past participle of":   &FormTemplate{Tags: []string{"f", "p", "of the", "past", "part"}},
	"feminine singular of":                 &FormTemplate{Tags: []string{"f", "s"}},
	"feminine singular past participle of": &FormTemplate{Tags: []string{"f", "s", "of the", "past", "part"}},
	"form of":                              &FormTemplate{}, // https://en.wiktionary.org/wiki/Template:form_of
	"former name of":                       &FormTemplate{},
	"frequentative of":                     &FormTemplate{Tags: []string{"freq"}},
	"future participle of":                 &FormTemplate{Tags: []string{"fut", "part"}},
	"genitive of":                          &FormTemplate{Tags: []string{"gen"}},
	"genitive plural of":                   &FormTemplate{Tags: []string{"gen", "p"}},
	"genitive singular of":                 &FormTemplate{Tags: []string{"gen", "s"}},
	"gerund of":                            &FormTemplate{Tags: []string{"gerund"}},
	"h-prothesis of":                       &FormTemplate{Text: "h-prothesized form of"},
	"hard mutation of":                     &FormTemplate{},
	"harmonic variant of":                  &FormTemplate{},
	"honorific alternative case form of":   &FormTemplate{Text: "honorific alternative letter-case form of", Shortcuts: []string{"honor alt case"}},
	"imperative of":                        &FormTemplate{Tags: []string{"impr"}},
	"imperfective form of":                 &FormTemplate{Tags: []string{"impfv"}},
	"indefinite plural of":                 &FormTemplate{Tags: []string{"indef", "p"}},
	"inflection of":                        &FormTemplate{Shortcuts: []string{"infl of"}, Inflection: true},
	"informal form of":                     &FormTemplate{},
	"informal spelling of":                 &FormTemplate{},
	"initialism of":                        &FormTemplate{Shortcuts: []string{"init of"}},
	"iterative of":                         &FormTemplate{Tags: []string{"iter"}},
	"lenition of":                          &FormTemplate{Text: "lenited form of"},
	"masculine noun of":                    &FormTemplate{Tags: []string{"m", "equivalent"}},
	"masculine of":                         &FormTemplate{Tags: []string{"m"}},
	"masculine plural of":                  &FormTemplate{Tags: []string{"m", "p"}},
	"masculine plural past participle of":  &FormTemplate{Tags: []string{"m", "p", "of the", "past", "part"}},
	"medieval spelling of":                 &FormTemplate{},
	"men's speech form of":                 &FormTemplate{},
	"misconstruction of":                   &FormTemplate{},
	"misromanization of":                   &FormTemplate{},
	"misspelling of":                       &FormTemplate{Shortcuts: []string{"missp"}},
	"mixed mutation of":                    &FormTemplate{},
	"nasal mutation of":                    &FormTemplate{},
	"negative of":                          &FormTemplate{Tags: []string{"neg", "form"}},
	"neuter plural of":                     &FormTemplate{Tags: []string{"n", "p"}},
	"neuter singular of":                   &FormTemplate{Tags: []string{"n", "s"}},
	"neuter singular past participle of":   &FormTemplate{Tags: []string{"n", "s", "of the", "past", "part"}},
	"nomen sacrum form of":                 &FormTemplate{},
	"nominalization of":                    &FormTemplate{Tags: []string{"nomzn"}},
	"nominative plural of":                 &FormTemplate{Tags: []string{"nom", "p"}},
	"nonstandard form of":                  &FormTemplate{},
	"nonstandard spelling of":              &FormTemplate{},
	"noun form of":                         &FormTemplate{},
	"nuqtaless form of":                    &FormTemplate{},
	"obsolete form of":                     &FormTemplate{Shortcuts: []string{"obs sp"}},
	"obsolete spelling of":                 &FormTemplate{},
	"obsolete typography of":               &FormTemplate{},
	"participle of":                        &FormTemplate{},
	"passive of":                           &FormTemplate{Tags: []string{"pasv"}},
	"passive participle of":                &FormTemplate{Tags: []string{"pass", "part"}},
	"passive past tense of":                &FormTemplate{Tags: []string{"pass", "past"}},
	"past active participle of":            &FormTemplate{Tags: []string{"pass", "actv", "ptcp"}},
	"past participle form of":              &FormTemplate{},
	"past participle of":                   &FormTemplate{Tags: []string{"past", "ptcp"}},
	"past passive participle of":           &FormTemplate{Tags: []string{"past", "pasv", "ptcp"}},
	"past tense of":                        &FormTemplate{Tags: []string{"past", "tense"}},
	"pejorative of":                        &FormTemplate{Tags: []string{"pej"}},
	"perfect participle of":                &FormTemplate{Tags: []string{"perf", "part"}},
	"perfective form of":                   &FormTemplate{Tags: []string{"pfv", "form"}},
	"plural of":                            &FormTemplate{Tags: []string{"p"}},
	"present active participle of":         &FormTemplate{Tags: []string{"pres", "act", "part"}},
	"present participle of":                &FormTemplate{Tags: []string{"pres", "ptcp"}},
	"present tense of":                     &FormTemplate{Tags: []string{"pres"}},
	"pronunciation spelling of":            &FormTemplate{},
	"pronunciation variant of":             &FormTemplate{},
	"rare form of":                         &FormTemplate{},
	"rare spelling of":                     &FormTemplate{Shortcuts: []string{"rare sp"}},
	"reflexive of":                         &FormTemplate{Tags: []string{"refl"}},
	"rfform":                               &FormTemplate{Text: "unknown form of"},
	"romanization of":                      &FormTemplate{},
	"short for":                            &FormTemplate{},
	"singular of":                          &FormTemplate{Tags: []string{"s"}},
	"singulative of":                       &FormTemplate{Tags: []string{"sgl"}},
	"slender form of":                      &FormTemplate{},
	"soft mutation of":                     &FormTemplate{},
	"spelling of":                          &FormTemplate{},
	"standard form of":                     &FormTemplate{},
	"standard spelling of":                 &FormTemplate{Shortcuts: []string{"standard sp"}},
	"superlative attributive of":           &FormTemplate{Tags: []string{"supd"}},
	"superlative of":                       &FormTemplate{Text: "superlative degree of"},
	"superlative predicative of":           &FormTemplate{Text: "superlative (when used predicatively) of"},
	"superseded spelling of":               &FormTemplate{Shortcuts: []string{"sup sp"}},
	"supine of":                            &FormTemplate{Tags: []string{"supine"}},
	"syncopic form of":                     &FormTemplate{},
	"synonym of":                           &FormTemplate{Shortcuts: []string{"syn of"}},
	"t-prothesis of":                       &FormTemplate{Text: "t-prothesized form of"},
	"uncommon form of":                     &FormTemplate{},
	"uncommon spelling of":                 &FormTemplate{},
	"verbal noun of":                       &FormTemplate{Tags: []string{"vnoun"}},
	"verb form of":                         &FormTemplate{Inflection: true},
	"vocative plural of":                   &FormTemplate{Tags: []string{"voc", "p"}},
	"vocative singular of":                 &FormTemplate{Tags: []string{"voc", "s"}},

	// English templates
	// https://en.wiktionary.org/wiki/Category:English_form-of_templates

	"en-archaic third-person singular of":       &FormTemplate{Language: "en", Text: "(archaic) third-person singular simple present indicative form of"},
	"en-comparative of":                         &FormTemplate{Language: "en", Text: "comparative form of"},
	"en-archaic second-person singular of":      &FormTemplate{Language: "en", Text: "second-person singular simple present form of"},
	"en-archaic second-person singular past of": &FormTemplate{Language: "en", Text: "second-person singular simple past form of"},
	"en-ing form of":                            &FormTemplate{Language: "en", Text: "present participle and gerund of"},
	"en-simple past of":                         &FormTemplate{Language: "en", Text: "simple past tense of"},
	"en-irregular plural of":                    &FormTemplate{Language: "en", Text: "plural of"},
	"en-past of":                                &FormTemplate{Language: "en", Text: "simple past tense and past participle of"},
	"en-superlative of":                         &FormTemplate{Language: "en", Text: "superlative form of"},
	"en-third-person singular of":               &FormTemplate{Language: "en", Text: "third-person singular simple present indicative form of"},

	// TODO: Support https://en.wiktionary.org/wiki/Category:Spanish_form-of_templates
}

var formOfShortcuts = map[string]*FormTemplate{}

type FormTemplateType string

type FormTemplate struct {
	Language   string
	Tags       []string
	Shortcuts  []string
	Text       string
	Inflection bool
}

type ListItem struct {
	Prefix string
	Links  []string
}

type TextBuffer []string

func init() {
	wordTypeRegex = regexp.MustCompile("^([^0-9]+)(?: [0-9]+)?$")

	for text, form := range formOfMap {
		if form.Text == "" {
			form.Text = text
		}
		for _, shortcut := range form.Shortcuts {
			formOfShortcuts[shortcut] = form
		}
	}
}

func definitionSection(section sectionType) bool {
	return section == nounSection || section == adjectiveSection || section == verbSection || section == adverbSection || section == articleSection || section == prepositionSection || section == pronounSection || section == conjunctionSection || section == interjectionSection || section == numeralSection || section == particleSection || section == determinerSection
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
						append(language.Links, language.listItem.TplLinks(langMap, w.Name)...)
				}
				if !language.IsEmpty() {
					if w.Languages == nil {
						w.Languages = map[string]*Language{}
					}
					w.Languages[language.Code] = language
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
						w.Languages = map[string]*Language{}
					}
					w.Languages[language.Code] = language
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
						append(language.Links, language.listItem.TplLinks(langMap, w.Name)...)
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
			if language != nil {
				language.flushDefinition()
				language.descendantLang = nil
			}
		case itemLeftLink:
			if language != nil {
				language.linkBuffer = &LinkBuffer{}
			}
		case itemLink:
			if language != nil {
				if language.listItem != nil {
					language.listItem.Links = append(language.listItem.Links, i.val)
				} else {
					language.linkBuffer.Link = i.val
				}
			}
		case itemLinkName:
			if language != nil {
				language.linkBuffer.Name = &i.val
			}
		case itemRightLink:
			if language != nil {
				if language.shouldDefineLink() {
					var link string
					if language.linkBuffer.Name != nil {
						link = *language.linkBuffer.Name
					} else {
						link = language.linkBuffer.Link
					}
					language.definitionBuffer = append(language.definitionBuffer, link)
				} else if language.subSection == descendantsSection {
					if language.descendantLang != nil {
						tplLink := toTplLink(langMap, *language.descendantLang, language.linkBuffer.Link, w.Name)
						if tplLink != nil {
							language.Links = append(language.Links, *tplLink)
						}
					}
				}
				language.linkBuffer = nil
			}
		case itemText:
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
				} else {
					sectionMatches := wordTypeRegex.FindStringSubmatch(i.val)
					var setWordSection = false

					if len(sectionMatches) > 1 {
						if wordSection, ok := wordTypeMap[sectionMatches[1]]; (language.sectionDepth == 2 || language.sectionDepth == 3) && ok {
							language.section = wordSection
							setWordSection = true
						}
					}

					if !setWordSection {
						// This will exclude subsections named "Etymology" for now, e.g. https://en.wiktionary.org/wiki/taco#Noun_4
						language.section = unknownSection

						if language.sectionDepth >= 3 && i.val == "Descendants" {
							language.subSection = descendantsSection
						} else {
							language.subSection = unknownSection
						}
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

			namedParam = nil

			// Don't support nested templates for now
			if language == nil || template == nil || templateDepth != 0 {
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
						language.descendantLang = &desc.Lang
					}
				case "l", "link":
					link := template.ToLink()
					if _, ok := langMap[link.Lang]; ok {
						language.Links =
							append(language.Links, link)
						language.descendantLang = &link.Lang
					}
				case "desctree", "descendants tree":
					descTree := template.ToDescTree()
					if _, ok := langMap[descTree.Lang]; ok {
						language.DescTrees = append(language.DescTrees, descTree)
						language.descendantLang = &descTree.Lang

						// Also add descendant tree as a descendant
						desc := tpl.Descendant{Lang: descTree.Lang, Word: descTree.Word}
						language.Descendants = append(language.Descendants, desc)
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
							language.descendantLang = &etymTree.Lang
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
				case "m", "mention":
					mention := template.ToMention()
					if language.definitionBuffer != nil {
						language.definitionBuffer = append(language.definitionBuffer, mention.Text())
					}
				case "gloss":
					gloss := template.ToGloss()
					if language.definitionBuffer != nil {
						language.definitionBuffer = append(language.definitionBuffer, gloss.Text())
					}
				case "non-gloss definition", "non-gloss", "non gloss", "ngd", "n-g":
					nonGloss := template.ToNonGloss()
					if language.definitionBuffer != nil {
						language.definitionBuffer = append(language.definitionBuffer, nonGloss.Text())
					}
				case "label", "lbl", "lb":
					label := template.ToLabel()
					if language.definitionBuffer != nil {
						language.definitionBuffer = append(language.definitionBuffer, label.Text())
					}
				case "qualifier", "qual", "q", "i":
					qualifier := template.ToQualifier()
					if language.definitionBuffer != nil {
						language.definitionBuffer = append(language.definitionBuffer, qualifier.Text())
					}
				case "frac":
					frac := template.ToFrac()
					if language.definitionBuffer != nil {
						language.definitionBuffer = append(language.definitionBuffer, frac.Text())
					}
					// TODO: Remove
				case "feminine noun of":
					femNoun := template.ToFemNoun()
					if language.definitionBuffer != nil {
						language.definitionBuffer = append(language.definitionBuffer, femNoun.Text())
					}

					// Spanish forms
				case "es-verb form of":
					spanishVerb := template.ToSpanishVerb()
					if language.definitionBuffer != nil {
						language.definitionBuffer = append(language.definitionBuffer, spanishVerb.Text())
						language.definitionRoot = &RootWord{Lang: spanishLang, Name: spanishVerb.Word}
					}
				case "es-compound of":
					spanishCompound := template.ToSpanishCompound()
					if language.definitionBuffer != nil {
						language.definitionBuffer = append(language.definitionBuffer, spanishCompound.Text())
						language.definitionRoot = &RootWord{Lang: spanishLang, Name: spanishCompound.Word()}
					}

				default:
					if template.Action == "form of" {
						formOf := template.ToFormOfGeneric()
						if language.definitionBuffer != nil {
							language.definitionBuffer = append(language.definitionBuffer, formOf.Text())
							language.definitionRoot = &RootWord{Lang: formOf.Lang, Name: formOf.DisplayWord()}
						}
					} else {
						var formTpl *FormTemplate

						if defn, ok := formOfMap[template.Action]; ok {
							formTpl = defn
						} else if defn, ok := formOfShortcuts[template.Action]; ok {
							formTpl = defn
						}

						if formTpl != nil {

							// If form-of template specifies language, manually inject it to template parameters
							// so word is still in second position.
							if formTpl.Language != "" {
								template.Parameters = append([]string{formTpl.Language}, template.Parameters...)
							}

							formOf := template.ToFormOf(formTpl.Text, formTpl.Tags...)
							if language.definitionBuffer != nil {
								language.definitionBuffer = append(language.definitionBuffer, formOf.Text())
								language.definitionRoot = &RootWord{Lang: formOf.Lang, Name: formOf.DisplayWord()}
							}
						}
					}
				}
			}
		case itemAction:
			// Don't support nested templates for now.
			if template == nil || templateDepth != 1 {
				break
			}
			template.Action = i.val
		case itemParamDelim:
			namedParam = nil
		case itemParamText:

			// Don't support nested templates for now.
			if template == nil || templateDepth != 1 {
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
			// Don't support nested templates for now
			if template == nil || templateDepth != 1 {
				break
			}
			namedParam = &tpl.Parameter{Name: i.val}
		}
	}

	return w, nil
}

func (li *ListItem) TplLinks(langMap map[string]bool, parent string) (ls []tpl.Link) {
	for _, link := range li.Links {
		if canonical, ok := lang.CanonicalLangs[li.Prefix]; ok {
			if _, ok := langMap[canonical.Code]; ok {
				tplLink := toTplLink(langMap, canonical.Code, link, parent)
				if tplLink == nil {
					continue
				}
				ls = append(ls, *tplLink)
			}
		}
	}
	return ls
}

func toTplLink(langMap map[string]bool, lang, linkText string, parent string) *tpl.Link {
	if strings.Contains(linkText, ":") {
		return nil
	}
	parts := strings.Split(linkText, "#")
	var link string
	if parts[0] != "" {
		link = parts[0]
	} else {
		link = parent
	}
	return &tpl.Link{Lang: lang, Word: link}
}
