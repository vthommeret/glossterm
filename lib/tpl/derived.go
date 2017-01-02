package tpl

import "reflect"

// https://en.wiktionary.org/wiki/Template:derived
type Derived struct {
	Lang         string
	FromLang     string
	FromWord     string
	Alt          string `names:"alt"`
	Gloss        string `names:"t,gloss"`
	PartOfSpeech string `names:"pos"`
	Literal      string `names:"lit"`
}

func (tpl *Template) ToDerived() Derived {
	d := Derived{}
	tpl.toConcrete(reflect.TypeOf(d), reflect.ValueOf(&d))
	d.FromWord = toEntryName(d.FromLang, d.FromWord)
	return d
}
