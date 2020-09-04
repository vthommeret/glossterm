package tpl

import (
	"fmt"
	"reflect"
)

// https://en.wiktionary.org/wiki/Template:apocopic_form_of
type Apocopic struct {
	Lang  string `lang:"true" json:"lang,omitempty" firestore:"lang,omitempty"`
	Word  string `json:"word,omitempty" firestore:"word,omitempty"`
	Alt   string `json:"alt,omitempty" firestore:"alt,omitempty"`
	Gloss string `names:"t" json:"gloss,omitempty" firestore:"gloss,omitempty"`
}

func (tpl *Template) ToApocopic() Apocopic {
	a := Apocopic{}
	tpl.toConcrete(reflect.TypeOf(a), reflect.ValueOf(&a))
	a.Word = toEntryName(a.Lang, a.Word)
	return a
}

func (a *Apocopic) Text() string {
	var word string
	if a.Alt != "" {
		word = a.Alt
	} else {
		word = a.Word
	}
	var gloss string
	if a.Gloss != "" {
		gloss = fmt.Sprintf(" (%s)", a.Gloss)
	}
	return fmt.Sprintf("apocopic form of %s%s", word, gloss)
}
