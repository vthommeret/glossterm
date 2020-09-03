package tpl

import (
	"fmt"
	"reflect"
)

// https://en.wiktionary.org/wiki/Template:feminine_noun_of
type FemNoun struct {
	Lang string `lang:"true" json:"lang,omitempty" firestore:"lang,omitempty"`
	Word string `json:"word,omitempty" firestore:"word,omitempty"`
}

func (tpl *Template) ToFemNoun() FemNoun {
	fn := FemNoun{}
	tpl.toConcrete(reflect.TypeOf(fn), reflect.ValueOf(&fn))
	fn.Word = toEntryName(fn.Lang, fn.Word)
	return fn
}

func (fn *FemNoun) Text() string {
	return fmt.Sprintf("feminine equivalent of %s", fn.Word)
}
