package ml

import "fmt"

type Word struct {
	Value    string
	Sections []Section
}

type Section struct {
	Name  string
	Depth int
	Text  string
}

func Parse(p Page) Word {
	w := Word{
		Value: p.Title,
	}

	var inHeader bool
	var section *Section

	l := newLexer(p.Text)

Parse:
	for {
		i := l.NextItem()
		switch i.typ {
		case itemError:
			fmt.Printf("Error: %s\n", i.val)
			break Parse
		case itemEOF:
			if section != nil {
				w.Sections = append(w.Sections, *section)
			}
			break Parse
		case itemHeaderStart:
			if section != nil {
				w.Sections = append(w.Sections, *section)
			}
			section = &Section{Depth: i.depth}
			inHeader = true
		case itemHeaderEnd:
			inHeader = false
		case itemText:
			if inHeader {
				section.Name = i.val
			} else {
				section.Text = i.val
			}
		}
	}

	return w
}
