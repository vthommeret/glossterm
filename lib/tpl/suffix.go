package tpl

import "reflect"

// https://en.wiktionary.org/wiki/Template:suffix
type Suffix struct {
	Root   string
	Suffix string
	Lang   string `names:"lang", lang:"true"`
}

func (tpl *Template) ToSuffix() Suffix {
	s := Suffix{}
	tpl.toConcrete(reflect.TypeOf(s), reflect.ValueOf(&s))
	return s
}
