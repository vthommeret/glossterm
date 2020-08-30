package tpl

import "reflect"

// https://en.wiktionary.org/wiki/Template:inherited
type Inherited struct {
	Lang         string `lang:"true" firestore:"lang,omitempty"`
	FromLang     string `lang:"true" firestore:"fromLang,omitempty"`
	FromWord     string `firestore:"fromWord,omitempty"`
	Alt          string `names:"alt" firestore:"alt,omitempty"`
	Gloss        string `names:"t,gloss" firestore:"gloss,omitempty"`
	PartOfSpeech string `names:"pos" firestore:"partOfSpeech,omitempty"`
	Literal      string `names:"lit" firestore:"literal,omitempty"`
}

func (tpl *Template) ToInherited() Inherited {
	i := Inherited{}
	tpl.toConcrete(reflect.TypeOf(i), reflect.ValueOf(&i))
	i.FromWord = toEntryName(i.FromLang, i.FromWord)
	return i
}
