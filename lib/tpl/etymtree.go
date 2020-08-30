package tpl

import (
	"fmt"
	"reflect"
)

// https://en.wiktionary.org/wiki/Template:etymtree
type EtymTree struct {
	Lang     string `lang:"true" firestore:"lang,omitempty"`
	RootLang string `lang:"true" firestore:"rootLang,omitempty"`
	Word     string `names:"branch_term" firestore:"word,omitempty"`
}

func (tpl *Template) ToEtymTree() EtymTree {
	et := EtymTree{}
	tpl.toConcrete(reflect.TypeOf(et), reflect.ValueOf(&et))
	if et.RootLang != "" {
		et.Word = toEntryName(et.RootLang, et.Word)
	} else {
		et.Word = toEntryName(et.Lang, et.Word)
	}
	return et
}

func (et *EtymTree) ToEntryName() string {
	var lang string
	if et.RootLang != "" {
		lang = et.RootLang
	} else {
		lang = et.Lang
	}
	return fmt.Sprintf("Template:etymtree/%s/%s", lang, et.Word)
}
