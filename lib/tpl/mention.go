package tpl

import "reflect"

// https://en.wiktionary.org/wiki/Template:mention
type Mention struct {
	Lang         string `lang:"true" firestore:"lang,omitempty"`
	Word         string `firestore:"word,omitempty"`
	Alt          string `firestore:"alt,omitempty"`
	Gloss        string `names:"t" firestore:"gloss,omitempty"`
	PartOfSpeech string `names:"pos" firestore:"partOfSpeech,omitempty"`
	Literal      string `names:"lit" firestore:"literal,omitempty"`
}

func (tpl *Template) ToMention() Mention {
	m := Mention{}
	tpl.toConcrete(reflect.TypeOf(m), reflect.ValueOf(&m))
	m.Word = toEntryName(m.Lang, m.Word)
	return m
}
