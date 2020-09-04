package tpl

import (
	"fmt"
	"reflect"
)

// https://en.wiktionary.org/wiki/Template:synonym_of
type Synonym struct {
	Lang string `lang:"true" json:"lang,omitempty" firestore:"lang,omitempty"`
	Word string `json:"word,omitempty" firestore:"word,omitempty"`
}

func (tpl *Template) ToSynonym() Synonym {
	s := Synonym{}
	tpl.toConcrete(reflect.TypeOf(s), reflect.ValueOf(&s))
	s.Word = toEntryName(s.Lang, s.Word)
	return s
}

func (s *Synonym) Text() string {
	return fmt.Sprintf("synonym of %s", s.Word)
}
