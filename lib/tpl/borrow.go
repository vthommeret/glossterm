package tpl

import "reflect"

// https://en.wiktionary.org/wiki/Template:borrowing
type Borrow struct {
	Lang         string `lang:"true" json:"lang,omitempty" firestore:"lang,omitempty"`
	FromLang     string `lang:"true" json:"fromLang,omitempty" firestore:"fromLang,omitempty"`
	FromWord     string `json:"fromWord,omitempty" firestore:"fromWord,omitempty"`
	Alt          string `names:"alt" json:"alt,omitempty" firestore:"alt,omitempty"`
	Gloss        string `names:"t,gloss" json:"gloss,omitempty" firestore:"gloss,omitempty"`
	PartOfSpeech string `names:"pos" json:"partOfSpeech,omitempty" firestore:"partOfSpeech,omitempty"`
	Literal      string `names:"lit" json:"literal,omitempty" firestore:"literal,omitempty"`
}

func (tpl *Template) ToBorrow() Borrow {
	b := Borrow{}
	tpl.toConcrete(reflect.TypeOf(b), reflect.ValueOf(&b))
	b.FromWord = toEntryName(b.FromLang, b.FromWord)
	return b
}
