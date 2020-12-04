package tpl

import (
	"reflect"
)

// https://en.wiktionary.org/wiki/Template:desctree
type DescTree struct {
	Lang string `lang:"true" json:"lang,omitempty" firestore:"lang,omitempty"`
	Word string `json:"word,omitempty" firestore:"word,omitempty"`
}

func (tpl *Template) ToDescTree() DescTree {
	dt := DescTree{}
	tpl.toConcrete(reflect.TypeOf(dt), reflect.ValueOf(&dt))
	dt.Word = toEntryName(dt.Lang, dt.Word)
	return dt
}
