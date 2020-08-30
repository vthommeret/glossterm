package tpl

import "reflect"

// https://en.wiktionary.org/wiki/Template:link
type Link struct {
	Lang         string `lang:"true" firestore:"lang,omitempty"`
	Word         string `firestore:"word,omitempty"`
	Alt          string `firestore:"alt,omitempty"`
	Gloss        string `names:"t" firestore:"gloss,omitempty"`
	PartOfSpeech string `names:"pos" firestore:"partOfSpeech,omitempty"`
	Literal      string `names:"lit" firestore:"literal,omitempty"`
}

func (tpl *Template) ToLink() Link {
	l := Link{}
	tpl.toConcrete(reflect.TypeOf(l), reflect.ValueOf(&l))
	l.Word = toEntryName(l.Lang, l.Word)
	return l
}
