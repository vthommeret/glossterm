package tpl

import (
	"fmt"
	"reflect"
)

// https://en.wiktionary.org/wiki/Template:initialism_of
type Initialism struct {
	Lang       string `lang:"true" json:"lang,omitempty" firestore:"lang,omitempty"`
	Definition string `json:"definition,omitempty" firestore:"definition,omitempty"`
	Alt        string `json:"alt,omitempty" firestore:"alt,omitempty"`
	Gloss      string `names:"t" json:"gloss,omitempty" firestore:"gloss,omitempty"`
}

func (tpl *Template) ToInitialism() Initialism {
	i := Initialism{}
	tpl.toConcrete(reflect.TypeOf(i), reflect.ValueOf(&i))
	return i
}

func (i *Initialism) Text() string {
	var defn string
	if i.Alt != "" {
		defn = i.Alt
	} else {
		defn = i.Definition
	}
	var gloss string
	if i.Gloss != "" {
		gloss = fmt.Sprintf(" (%s)", i.Gloss)
	}
	return fmt.Sprintf("initialism of %s%s", defn, gloss)
}
