package tpl

import "reflect"

// https://en.wiktionary.org/wiki/Template:mention
type Mention struct {
	Lang         string
	Word         string
	Alt          string
	Gloss        string `names:"t"`
	PartOfSpeech string `names:"pos"`
	Literal      string `names:"lit"`
}

func (tpl *Template) ToMention() Mention {
	m := Mention{}
	tpl.toConcrete(reflect.TypeOf(m), reflect.ValueOf(&m))
	return m
}