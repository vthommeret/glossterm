package ml

import (
	"reflect"
	"strings"
)

// https://en.wiktionary.org/wiki/Template:mention
type Mention struct {
	Lang         string
	Word         string
	Alt          string
	Gloss        string `names:"t"`
	PartOfSpeech string `names:"pos"`
	Literal      string `names:"lit"`
}

func (tpl *Template) ToMention() Mention {
	m := Mention{}
	tpl.toConcrete(reflect.TypeOf(m), reflect.ValueOf(&m))
	return m
}

// https://en.wiktionary.org/wiki/Template:borrowing
type Borrow struct {
	Lang         string
	FromLang     string
	FromWord     string
	Alt          string `names:"alt"`
	Gloss        string `names:"t,gloss"`
	PartOfSpeech string `names:"pos"`
	Literal      string `names:"lit"`
}

func (tpl *Template) ToBorrow() Borrow {
	b := Borrow{}
	tpl.toConcrete(reflect.TypeOf(b), reflect.ValueOf(&b))
	return b
}

// toConcrete turns a generic template into a concrete struct.
func (tpl *Template) toConcrete(t reflect.Type, v reflect.Value) {
	v = v.Elem()

	// Set positional parameters.
	n := len(tpl.Parameters)
	for i := 0; i < t.NumField(); i++ {
		if n > i {
			v.Field(i).SetString(tpl.Parameters[i])
		}
	}

	// Create named parameter map.
	paramMap := make(map[string]string)
	for _, p := range tpl.NamedParameters {
		if p.Name != "" {
			paramMap[p.Name] = p.Value
		}
	}

	// Set named parameters.
	for i := 0; i < t.NumField(); i++ {
		for _, p := range strings.Split(t.Field(i).Tag.Get("names"), ",") {
			if val, ok := paramMap[p]; ok {
				v.Field(i).SetString(val)
			}
		}
	}
}
