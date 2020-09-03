package tpl

import (
	"fmt"
	"reflect"
)

// https://en.wiktionary.org/wiki/Template:gloss
type Gloss struct {
	Gloss string `json:"gloss,omitempty" firestore:"gloss,omitempty"`
}

func (tpl *Template) ToGloss() Gloss {
	g := Gloss{}
	tpl.toConcrete(reflect.TypeOf(g), reflect.ValueOf(&g))
	return g
}

func (g *Gloss) Text() string {
	return fmt.Sprintf("(%s)", g.Gloss)
}
