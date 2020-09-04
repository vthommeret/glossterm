package tpl

import (
	"fmt"
	"reflect"
)

// https://en.wiktionary.org/wiki/Template:past_participle_of
type PastParticiple struct {
	Lang string `lang:"true" json:"lang,omitempty" firestore:"lang,omitempty"`
	Word string `json:"word,omitempty" firestore:"word,omitempty"`
}

func (tpl *Template) ToPastParticiple() PastParticiple {
	pp := PastParticiple{}
	tpl.toConcrete(reflect.TypeOf(pp), reflect.ValueOf(&pp))
	pp.Word = toEntryName(pp.Lang, pp.Word)
	return pp
}

func (pp *PastParticiple) Text() string {
	return fmt.Sprintf("past participle of %s", pp.Word)
}
