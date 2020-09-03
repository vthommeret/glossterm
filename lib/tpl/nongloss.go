package tpl

import (
	"reflect"
)

// https://en.wiktionary.org/wiki/Template:non-gloss
type NonGloss struct {
	NonGloss string `json:"text,omitempty" firestore:"text,omitempty"`
}

func (tpl *Template) ToNonGloss() NonGloss {
	ng := NonGloss{}
	tpl.toConcrete(reflect.TypeOf(ng), reflect.ValueOf(&ng))
	return ng
}

func (ng *NonGloss) Text() string {
	return ng.NonGloss
}
