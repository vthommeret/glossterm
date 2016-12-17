package ml

import "fmt"

type Word struct {
	Value     string
	Languages []Language
}

type Language struct {
	Name     string
	Sections []Section
}

type Section struct {
	Name  string
	Depth int
	Text  string
}

func Parse(p Page) (Word, error) {
	w := Word{
		Value: p.Title,
	}

	var inLanguage bool
	var inSection bool

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
			if section != nil {
				language.Sections = append(language.Sections, *section)
			}
			if language != nil {
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
					language.Sections = append(language.Sections, *section)
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
			} else {
				section.Text = i.val
			}
		}
	}

	return w, nil
}
