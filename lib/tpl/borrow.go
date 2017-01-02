package tpl

import "reflect"

// https://en.wiktionary.org/wiki/Template:borrowing
type Borrow struct {
	Lang         string
	FromLang     string
	FromWord     string
	Alt          string `names:"alt"`
	Gloss        string `names:"t,gloss"`
	PartOfSpeech string `names:"pos"`
	Literal      string `names:"lit"`
}

func (tpl *Template) ToBorrow() Borrow {
	b := Borrow{}
	tpl.toConcrete(reflect.TypeOf(b), reflect.ValueOf(&b))
	b.FromWord = toEntryName(b.FromLang, b.FromWord)
	return b
}
