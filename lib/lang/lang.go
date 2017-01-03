package lang

import (
	"bytes"
	"strings"
)

const (
	macron = '\u0304'
	breve  = '\u0306'
	diaer  = '\u0308'
)

var DefaultLangs = []string{"en", "es", "fr", "la", "LL"}
var DefaultLangMap map[string]bool

var CanonicalLangs map[string]Lang

func init() {
	DefaultLangMap = make(map[string]bool)
	for _, l := range DefaultLangs {
		DefaultLangMap[l] = true
	}
	CanonicalLangs = make(map[string]Lang)
	for _, l := range Langs {
		CanonicalLangs[l.Canonical] = l
		if l.EntryNameStrip != nil {
			l.entryNameStripMap = make(map[rune]bool)
			for _, r := range l.EntryNameStrip {
				l.entryNameStripMap[r] = true
			}
		}
	}
}

// Based on https://en.wiktionary.org/wiki/Module:languages#Language:makeEntryName
func (l *Lang) MakeEntryName(s string) string {
	s = strings.TrimRight(strings.TrimLeft(s, "¿¡"), "؟?!;՛՜ ՞ ՟？！।")
	if l.EntryNameMap == nil && l.entryNameStripMap == nil {
		return s
	}
	var buffer bytes.Buffer
	for _, r := range s {
		if t, ok := l.EntryNameMap[r]; ok {
			buffer.WriteRune(t)
		} else if _, ok := l.entryNameStripMap[r]; !ok {
			buffer.WriteRune(r)
		}
	}
	return buffer.String()
}

type Lang struct {
	Code              string
	Canonical         string
	Other             []string
	EntryNameMap      map[rune]rune
	EntryNameStrip    []rune
	entryNameStripMap map[rune]bool
}

// From https://en.wiktionary.org/wiki/Category:Language_data_modules
// Only supporting a subset of languages for now.
var Langs = map[string]Lang{
	"de": {
		Code:      "de",
		Canonical: "German",
		Other:     []string{"High German", "New High German", "Deutsch"},
	},
	"en": {
		Code:      "en",
		Canonical: "English",
		Other:     []string{"Modern English", "New English", "Hawaiian Creole English", "Hawai'ian Creole English", "Hawaiian Creole", "Hawai'ian Creole", "Polari", "Yinglish"},
	},
	"es": {
		Code:      "es",
		Canonical: "Spanish",
		Other:     []string{"Castilian", "Amazonian Spanish", "Amazonic Spanish", "Loreto-Ucayali Spanish"},
	},
	"fr": {
		Code:      "fr",
		Canonical: "French",
		Other:     []string{"Modern French"},
	},
	"grc": {
		Code:      "grc",
		Canonical: "Ancient Greek",
		EntryNameMap: map[rune]rune{
			'Ᾰ': 'A', 'Ᾱ': 'A', 'ᾰ': 'α', 'ᾱ': 'α', 'Ῐ': 'I', 'Ῑ': 'I', 'ῐ': 'ι', 'ῑ': 'ι', 'Ῠ': 'Y', 'Ῡ': 'Y', 'ῠ': 'υ', 'ῡ': 'υ', 'µ': 'μ',
		},
	},
	"it": {
		Code:      "it",
		Canonical: "Italian",
	},
	"la": {
		Code:      "la",
		Canonical: "Latin",
		EntryNameMap: map[rune]rune{
			'Ā': 'A', 'Ă': 'A', 'ā': 'a', 'ă': 'a', 'Ē': 'E', 'Ĕ': 'E', 'ē': 'e', 'ĕ': 'e', 'ë': 'e', 'Ī': 'I', 'Ĭ': 'I', 'Ï': 'I', 'ī': 'i', 'ĭ': 'i', 'ï': 'i', 'Ō': 'O', 'Ŏ': 'O', 'ō': 'o', 'ŏ': 'o', 'Ū': 'U', 'Ŭ': 'U', 'Ü': 'U', 'ū': 'u', 'ŭ': 'u', 'ü': 'u', 'Ȳ': 'Y', 'ȳ': 'y',
		},
		EntryNameStrip: []rune{
			macron, breve, diaer,
		},
	},
	"pt": {
		Code:      "pt",
		Canonical: "Portuguese",
		Other:     []string{"Modern Portuguese"},
	},
}
