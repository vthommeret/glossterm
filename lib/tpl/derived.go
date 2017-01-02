package tpl

import "reflect"

// https://en.wiktionary.org/wiki/Template:derived
type Derived struct {
	Lang         string `lang:"true"`
	FromLang     string `lang:"true"`
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
