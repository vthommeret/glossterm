package tpl

import "reflect"

// https://en.wiktionary.org/wiki/Template:link
type Link struct {
	Lang         string `lang:"true" json:"lang,omitempty" firestore:"lang,omitempty"`
	Word         string `json:"word,omitempty" firestore:"word,omitempty"`
	Alt          string `json:"alt,omitempty" firestore:"alt,omitempty"`
	Gloss        string `names:"t" json:"gloss,omitempty" firestore:"gloss,omitempty"`
	PartOfSpeech string `names:"pos" json:"partOfSpeech,omitempty" firestore:"partOfSpeech,omitempty"`
	Literal      string `names:"lit" json:"literal,omitempty" firestore:"literal,omitempty"`
}

func (tpl *Template) ToLink() Link {
	l := Link{}
	tpl.toConcrete(reflect.TypeOf(l), reflect.ValueOf(&l))
	l.Word = toEntryName(l.Lang, l.Word)
	return l
}

func (l *Link) Text() string {
	return l.Word
}
