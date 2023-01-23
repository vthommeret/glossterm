package lang

import (
	"bytes"
	"strings"
)

const (
	acute      = '\u0301'
	macron     = '\u0304'
	breve      = '\u0306'
	diaer      = '\u0308'
	dotabove   = '\u0307'
	apostrophe = '\u0027'
)

var DefaultLangs = []string{"en", "ang", "enm", "es", "pt", "fr", "fro", "frm", "la", "LL"}
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
	"ang": {
		Code:      "ang",
		Canonical: "Old English",
		Other:     []string{"Anglo-Saxon"},
		EntryNameMap: map[rune]rune{
			'Ā': 'A', 'Á': 'A', 'ā': 'a', 'á': 'a', 'Ǣ': 'Æ', 'Ǽ': 'Æ', 'ǣ': 'æ', 'ǽ': 'æ', 'Ċ': 'C', 'ċ': 'c', 'Ē': 'E', 'É': 'E', 'ē': 'e', 'é': 'e', 'Ġ': 'G', 'ġ': 'g', 'Ī': 'I', 'Í': 'I', 'ī': 'i', 'í': 'i', 'Ō': 'O', 'Ó': 'O', 'ō': 'o', 'ó': 'o', 'Ū': 'U', 'Ú': 'U', 'ū': 'u', 'ú': 'u', 'Ȳ': 'Y', 'Ý': 'Y', 'ȳ': 'y', 'ý': 'y', 'Ƿ': 'W', 'ƿ': 'w',
		},
		EntryNameStrip: []rune{
			macron, acute, dotabove,
		},
	},
	"enm": {
		Code:      "enm",
		Canonical: "Middle English",
		Other:     []string{"Medieval English", "Mediaeval English"},
		EntryNameMap: map[rune]rune{
			'Ā': 'A', 'Á': 'A', 'ā': 'a', 'á': 'a', 'Ǣ': 'Æ', 'Ǽ': 'Æ', 'ǣ': 'æ', 'ǽ': 'æ', 'Ċ': 'C', 'ċ': 'c', 'Ē': 'E', 'É': 'E', 'Ė': 'E', 'ē': 'e', 'é': 'e', 'ė': 'e', 'Ġ': 'G', 'ġ': 'g', 'Ī': 'I', 'Í': 'I', 'ī': 'i', 'í': 'i', 'Ō': 'O', 'Ó': 'O', 'ō': 'o', 'ó': 'o', 'Ū': 'U', 'Ú': 'U', 'ū': 'u', 'ú': 'u', 'Ȳ': 'Y', 'Ý': 'Y', 'ȳ': 'y', 'ý': 'y',
		},
		EntryNameStrip: []rune{
			macron, acute, dotabove,
		},
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
	"fro": {
		Code:      "fro",
		Canonical: "Old French",
		EntryNameMap: map[rune]rune{
			'á': 'a', 'à': 'a', 'â': 'a', 'ä': 'a', 'é': 'e', 'è': 'e', 'ê': 'e', 'ë': 'e', 'í': 'i', 'ì': 'i', 'î': 'i', 'ï': 'i', 'ó': 'o', 'ò': 'o', 'ô': 'o', 'ö': 'o', 'ú': 'u', 'ù': 'u', 'û': 'u', 'ü': 'u', 'ý': 'y', 'ỳ': 'y', 'ŷ': 'y', 'ÿ': 'y', 'ç': 'c',
		},
		EntryNameStrip: []rune{
			apostrophe,
		},
	},
	"frm": {
		Code:      "frm",
		Canonical: "Middle French",
		EntryNameMap: map[rune]rune{
			'á': 'a', 'à': 'a', 'â': 'a', 'ä': 'a', 'é': 'e', 'è': 'e', 'ê': 'e', 'ë': 'e', 'í': 'i', 'ì': 'i', 'î': 'i', 'ï': 'i', 'ó': 'o', 'ò': 'o', 'ô': 'o', 'ö': 'o', 'ú': 'u', 'ù': 'u', 'û': 'u', 'ü': 'u', 'ý': 'y', 'ỳ': 'y', 'ŷ': 'y', 'ÿ': 'y', 'ç': 'c',
		},
		EntryNameStrip: []rune{
			apostrophe,
		},
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
