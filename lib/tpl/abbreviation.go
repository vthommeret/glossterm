package tpl

import (
	"fmt"
	"reflect"
)

// https://en.wiktionary.org/wiki/Template:abbreviation_of
type Abbreviation struct {
	Lang       string `lang:"true" json:"lang,omitempty" firestore:"lang,omitempty"`
	Definition string `json:"definition,omitempty" firestore:"definition,omitempty"`
	Alt        string `json:"alt,omitempty" firestore:"alt,omitempty"`
	Gloss      string `names:"t" json:"gloss,omitempty" firestore:"gloss,omitempty"`
}

func (tpl *Template) ToAbbreviation() Abbreviation {
	a := Abbreviation{}
	tpl.toConcrete(reflect.TypeOf(a), reflect.ValueOf(&a))
	return a
}

func (a *Abbreviation) Text() string {
	var defn string
	if a.Alt != "" {
		defn = a.Alt
	} else {
		defn = a.Definition
	}
	var gloss string
	if a.Gloss != "" {
		gloss = fmt.Sprintf(" (%s)", a.Gloss)
	}
	return fmt.Sprintf("abbreviation of %s%s", defn, gloss)
}
