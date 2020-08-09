package tpl

import "reflect"

// https://en.wiktionary.org/wiki/Template:suffix
type Suffix struct {
	Lang   string `names:"lang", lang:"true"`
	Root   string
	Suffix string
}

func (tpl *Template) ToSuffix() Suffix {
	s := Suffix{}
	tpl.toConcrete(reflect.TypeOf(s), reflect.ValueOf(&s))
	return s
}
