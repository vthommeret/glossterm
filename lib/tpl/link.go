package tpl

import "reflect"

// https://en.wiktionary.org/wiki/Template:link
type Link struct {
	Lang         string
	Word         string
	Alt          string
	Gloss        string `names:"t"`
	PartOfSpeech string `names:"pos"`
	Literal      string `names:"lit"`
}

func (tpl *Template) ToLink() Link {
	l := Link{}
	tpl.toConcrete(reflect.TypeOf(l), reflect.ValueOf(&l))
	l.Word = toEntryName(l.Lang, l.Word)
	return l
}
