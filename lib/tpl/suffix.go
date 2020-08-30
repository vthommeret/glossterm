package tpl

import "reflect"

// https://en.wiktionary.org/wiki/Template:suffix
type Suffix struct {
	Lang   string `names:"lang" lang:"true" firestore:"lang,omitempty"`
	Root   string `firestore:"root,omitempty"`
	Suffix string `firestore:"suffix,omitempty"`
}

func (tpl *Template) ToSuffix() Suffix {
	s := Suffix{}
	tpl.toConcrete(reflect.TypeOf(s), reflect.ValueOf(&s))
	return s
}
