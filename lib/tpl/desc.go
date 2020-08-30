package tpl

import "reflect"

// https://en.wiktionary.org/wiki/Template:descendant
type Descendant struct {
	Lang string `lang:"true" firestore:"lang,omitempty"`
	Word string `firestore:"word,omitempty"`
}

func (tpl *Template) ToDescendant() Descendant {
	d := Descendant{}
	tpl.toConcrete(reflect.TypeOf(d), reflect.ValueOf(&d))
	d.Word = toEntryName(d.Lang, d.Word)
	return d
}
