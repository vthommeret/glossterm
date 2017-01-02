package tpl

import "reflect"

// https://en.wiktionary.org/wiki/Template:prefix
type Prefix struct {
	Prefix string
	Root   string
	Lang   string `names:"lang", lang:"true"`
}

func (tpl *Template) ToPrefix() Prefix {
	p := Prefix{}
	tpl.toConcrete(reflect.TypeOf(p), reflect.ValueOf(&p))
	return p
}
