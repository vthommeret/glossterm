package tpl

import (
	"fmt"
	"reflect"
)

// https://en.wiktionary.org/wiki/Template:frac
type Frac struct {
	Num1 string `json:"num1,omitempty" firestore:"num1,omitempty"`
	Num2 string `json:"num2,omitempty" firestore:"num2,omitempty"`
	Num3 string `json:"num3,omitempty" firestore:"num3,omitempty"`
}

func (tpl *Template) ToFrac() Frac {
	f := Frac{}
	tpl.toConcrete(reflect.TypeOf(f), reflect.ValueOf(&f))
	return f
}

func (f *Frac) Text() string {
	if f.Num3 != "" {
		i := f.Num1
		n := f.Num2
		d := f.Num3
		return fmt.Sprintf("%s and %s/%s", i, n, d)
	} else if f.Num2 != "" {
		n := f.Num1
		d := f.Num2
		return fmt.Sprintf("%s/%s", n, d)
	}
	d := f.Num1
	return fmt.Sprintf("1/%s", d)
}
