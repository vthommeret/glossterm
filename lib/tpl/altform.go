package tpl

import (
	"fmt"
	"reflect"
)

// https://en.wiktionary.org/wiki/Template:alternative_form_of
type AltForm struct {
	Lang  string `lang:"true" json:"lang,omitempty" firestore:"lang,omitempty"`
	Word  string `json:"word,omitempty" firestore:"word,omitempty"`
	Alt   string `json:"alt,omitempty" firestore:"alt,omitempty"`
	Gloss string `json:"gloss,omitempty" firestore:"gloss,omitempty"`
}

func (tpl *Template) ToAltForm() AltForm {
	af := AltForm{}
	tpl.toConcrete(reflect.TypeOf(af), reflect.ValueOf(&af))
	af.Word = toEntryName(af.Lang, af.Word)
	return af
}

func (af *AltForm) Text() string {
	var word string
	if af.Alt != "" {
		word = af.Alt
	} else {
		word = af.Word
	}
	var gloss string
	if af.Gloss != "" {
		gloss = fmt.Sprintf(" (%s)", af.Gloss)
	}
	return fmt.Sprintf("alternative form of %s%s", word, gloss)
}
