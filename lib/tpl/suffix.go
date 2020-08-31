package tpl

import "reflect"

// https://en.wiktionary.org/wiki/Template:suffix
type Suffix struct {
	Lang   string `names:"lang" lang:"true" json:"lang,omitempty" firestore:"lang,omitempty"`
	Root   string `json:"root,omitempty" firestore:"root,omitempty"`
	Suffix string `json:"suffix,omitempty" firestore:"suffix,omitempty"`
}

func (tpl *Template) ToSuffix() Suffix {
	s := Suffix{}
	tpl.toConcrete(reflect.TypeOf(s), reflect.ValueOf(&s))
	return s
}
