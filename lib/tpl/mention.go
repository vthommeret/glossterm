package tpl

import (
	"fmt"
	"reflect"
)

// https://en.wiktionary.org/wiki/Template:mention
type Mention struct {
	Lang         string `lang:"true" json:"lang,omitempty" firestore:"lang,omitempty"`
	Word         string `json:"word,omitempty" firestore:"word,omitempty"`
	Alt          string `json:"alt,omitempty" firestore:"alt,omitempty"`
	Gloss        string `names:"t" json:"gloss,omitempty" firestore:"gloss,omitempty"`
	PartOfSpeech string `names:"pos" json:"partOfSpeech,omitempty" firestore:"partOfSpeech,omitempty"`
	Literal      string `names:"lit" json:"literal,omitempty" firestore:"literal,omitempty"`
}

func (tpl *Template) ToMention() Mention {
	m := Mention{}
	tpl.toConcrete(reflect.TypeOf(m), reflect.ValueOf(&m))
	m.Word = toEntryName(m.Lang, m.Word)
	return m
}

func (m *Mention) Text() string {
	var word string
	if m.Alt != "" {
		word = m.Alt
	} else {
		word = m.Word
	}
	var gloss string
	if m.Gloss != "" {
		gloss = fmt.Sprintf(" (%s)", m.Gloss)
	}
	return fmt.Sprintf("%s%s", word, gloss)
}
