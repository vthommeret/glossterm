package ml

import "fmt"

type Word struct {
	Value     string
	Languages []Language
	Text      string
	Sections  []Section
}

type Language struct {
	Name      string
	Text      string
	Etymology string
	Sections  []Section
}

type Section struct {
	Name  string
	Depth int
	Text  string
}

type sectionType int

const (
	unknownSection sectionType = iota
	etymologySection
)

func Parse(p Page) (Word, error) {
	w := Word{
		Value: p.Title,
	}

	var inLanguage bool
	var inSection bool
	var inSectionType sectionType

	var language *Language
	var section *Section

	l := newLexer(p.Text)

Parse:
	for {
		i := l.NextItem()
		switch i.typ {
		case itemError:
			return Word{}, fmt.Errorf("unable to parse: %s", i.val)
		case itemEOF:
			if language == nil {
				if section != nil {
					w.Sections = append(w.Sections, *section)
				}
			} else {
				if section != nil {
					language.Sections = append(language.Sections, *section)
				}
				w.Languages = append(w.Languages, *language)
			}
			break Parse
		case itemHeaderStart:
			if i.depth == 2 {
				if language != nil {
					w.Languages = append(w.Languages, *language)
				}
				language = &Language{Sections: []Section{}}
				inLanguage = true
			} else if i.depth > 2 {
				if section != nil {
					if language == nil {
						w.Sections = append(w.Sections, *section)
					} else {
						language.Sections = append(language.Sections, *section)
					}
				}
				section = &Section{Depth: i.depth - 1}
				inSection = true
			}
		case itemHeaderEnd:
			if i.depth == 2 {
				inLanguage = false
			} else if i.depth > 2 {
				inSection = false
			}
		case itemText:
			if inLanguage {
				language.Name = i.val
			} else if inSection {
				section.Name = i.val
				if section.Depth == 2 {
					switch section.Name {
					case "Etymology":
						inSectionType = etymologySection
					default:
						inSectionType = unknownSection
					}
				} else if section.Depth < 3 {
					inSectionType = unknownSection
				}
			} else {
				if language == nil {
					w.Text = i.val
				} else if section == nil {
					language.Text = i.val
				} else {
					section.Text = i.val
				}
				switch inSectionType {
				case etymologySection:
					language.Etymology = i.val
				}
			}
		}
	}

	return w, nil
}
