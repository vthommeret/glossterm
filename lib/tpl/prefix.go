package tpl

import "reflect"

// https://en.wiktionary.org/wiki/Template:prefix
type Prefix struct {
	Prefix string `json:"prefix,omitempty" firestore:"prefix,omitempty"`
	Root   string `json:"root,omitempty" firestore:"root,omitempty"`
	Lang   string `names:"lang" lang:"true" json:"lang,omitempty" firestore:"lang,omitempty"`
}

func (tpl *Template) ToPrefix() Prefix {
	p := Prefix{}
	tpl.toConcrete(reflect.TypeOf(p), reflect.ValueOf(&p))
	return p
}
