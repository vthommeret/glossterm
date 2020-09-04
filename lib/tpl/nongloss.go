package tpl

import (
	"reflect"
)

// https://en.wiktionary.org/wiki/Template:non-gloss
type NonGloss struct {
	Definition string `json:"definition,omitempty" firestore:"definition,omitempty"`
}

func (tpl *Template) ToNonGloss() NonGloss {
	ng := NonGloss{}
	tpl.toConcrete(reflect.TypeOf(ng), reflect.ValueOf(&ng))
	return ng
}

func (ng *NonGloss) Text() string {
	return ng.Definition
}
