package tpl

import "reflect"

// https://en.wiktionary.org/wiki/Template:cognate
type Cognate struct {
	Lang string `lang:"true" json:"lang,omitempty" firestore:"lang,omitempty"`
	Word string `json:"word,omitempty" firestore:"word,omitempty"`
}

func (tpl *Template) ToCognate() Cognate {
	c := Cognate{}
	tpl.toConcrete(reflect.TypeOf(c), reflect.ValueOf(&c))
	c.Word = toEntryName(c.Lang, c.Word)
	return c
}
