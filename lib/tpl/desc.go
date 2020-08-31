package tpl

import "reflect"

// https://en.wiktionary.org/wiki/Template:descendant
type Descendant struct {
	Lang string `lang:"true" json:"lang,omitempty" firestore:"lang,omitempty"`
	Word string `json:"word,omitempty" firestore:"word,omitempty"`
}

func (tpl *Template) ToDescendant() Descendant {
	d := Descendant{}
	tpl.toConcrete(reflect.TypeOf(d), reflect.ValueOf(&d))
	d.Word = toEntryName(d.Lang, d.Word)
	return d
}
