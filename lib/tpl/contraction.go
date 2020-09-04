package tpl

import (
	"fmt"
	"reflect"
)

// https://en.wiktionary.org/wiki/Template:contraction_of
type Contraction struct {
	Lang       string `lang:"true" json:"lang,omitempty" firestore:"lang,omitempty"`
	Definition string `json:"definition,omitempty" firestore:"definition,omitempty"`
	Alt        string `json:"alt,omitempty" firestore:"alt,omitempty"`
	Gloss      string `names:"t" json:"gloss,omitempty" firestore:"gloss,omitempty"`
}

func (tpl *Template) ToContraction() Contraction {
	c := Contraction{}
	tpl.toConcrete(reflect.TypeOf(c), reflect.ValueOf(&c))
	return c
}

func (c *Contraction) Text() string {
	var defn string
	if c.Alt != "" {
		defn = c.Alt
	} else {
		defn = c.Definition
	}
	var gloss string
	if c.Gloss != "" {
		gloss = fmt.Sprintf(" (%s)", c.Gloss)
	}
	return fmt.Sprintf("contraction of %s%s", defn, gloss)
}
