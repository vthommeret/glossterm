package tpl

import "reflect"

// https://en.wiktionary.org/wiki/Template:derived
type Derived struct {
	Lang         string `lang:"true" firestore:"lang,omitempty"`
	FromLang     string `lang:"true" firestore:"fromLang,omitempty"`
	FromWord     string `firestore:"fromWord,omitempty"`
	Alt          string `names:"alt" firestore:"alt,omitempty"`
	Gloss        string `names:"t,gloss" firestore:"gloss,omitempty"`
	PartOfSpeech string `names:"pos" firestore:"partOfSpeech,omitempty"`
	Literal      string `names:"lit" firestore:"literal,omitempty"`
}

func (tpl *Template) ToDerived() Derived {
	d := Derived{}
	tpl.toConcrete(reflect.TypeOf(d), reflect.ValueOf(&d))
	d.FromWord = toEntryName(d.FromLang, d.FromWord)
	return d
}
