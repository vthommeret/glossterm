package tpl

import "reflect"

// https://en.wiktionary.org/wiki/Template:derived
type Derived struct {
	Lang         string `lang:"true" json:"lang,omitempty" firestore:"lang,omitempty"`
	FromLang     string `lang:"true" json:"fromLang,omitempty" firestore:"fromLang,omitempty"`
	FromWord     string `json:"fromWord,omitempty" firestore:"fromWord,omitempty"`
	Alt          string `names:"alt" json:"alt,omitempty" firestore:"alt,omitempty"`
	Gloss        string `names:"t,gloss" json:"gloss,omitempty" firestore:"gloss,omitempty"`
	PartOfSpeech string `names:"pos" json:"partOfSpeech,omitempty" firestore:"partOfSpeech,omitempty"`
	Literal      string `names:"lit" json:"literal,omitempty" firestore:"literal,omitempty"`
}

func (tpl *Template) ToDerived() Derived {
	d := Derived{}
	tpl.toConcrete(reflect.TypeOf(d), reflect.ValueOf(&d))
	d.FromWord = toEntryName(d.FromLang, d.FromWord)
	return d
}
