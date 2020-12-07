package tpl

import "reflect"

// https://en.wiktionary.org/wiki/Template:etyl
type Etyl struct {
	Lang string `names:"lang" lang:"true" json:"lang,omitempty" firestore:"lang,omitempty"`
}

func (tpl *Template) ToEtyl() Etyl {
	e := Etyl{}
	tpl.toConcrete(reflect.TypeOf(e), reflect.ValueOf(&e))
	return e
}
