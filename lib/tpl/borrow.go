package tpl

import "reflect"

// https://en.wiktionary.org/wiki/Template:borrowing
type Borrow struct {
	Lang         string `lang:"true" firestore:"lang,omitempty"`
	FromLang     string `lang:"true" firestore:"fromLang,omitempty"`
	FromWord     string `firestore:"fromWord,omitempty"`
	Alt          string `names:"alt" firestore:"alt,omitempty"`
	Gloss        string `names:"t,gloss" firestore:"gloss,omitempty"`
	PartOfSpeech string `names:"pos" firestore:"partOfSpeech,omitempty"`
	Literal      string `names:"lit" firestore:"literal,omitempty"`
}

func (tpl *Template) ToBorrow() Borrow {
	b := Borrow{}
	tpl.toConcrete(reflect.TypeOf(b), reflect.ValueOf(&b))
	b.FromWord = toEntryName(b.FromLang, b.FromWord)
	return b
}
