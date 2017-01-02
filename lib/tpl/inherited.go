package tpl

import "reflect"

// https://en.wiktionary.org/wiki/Template:inherited
type Inherited struct {
	Lang         string `lang:"true"`
	FromLang     string `lang:"true"`
	FromWord     string
	Alt          string `names:"alt"`
	Gloss        string `names:"t,gloss"`
	PartOfSpeech string `names:"pos"`
	Literal      string `names:"lit"`
}

func (tpl *Template) ToInherited() Inherited {
	i := Inherited{}
	tpl.toConcrete(reflect.TypeOf(i), reflect.ValueOf(&i))
	i.FromWord = toEntryName(i.FromLang, i.FromWord)
	return i
}
