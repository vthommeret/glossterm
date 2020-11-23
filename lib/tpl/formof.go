package tpl

import (
	"fmt"
	"reflect"
	"strings"
)

// https://en.wiktionary.org/wiki/Module:form_of
// https://en.wiktionary.org/wiki/Module:form_of/templates
// https://en.wiktionary.org/wiki/Template:plural_of
// https://en.wiktionary.org/wiki/Category:Form-of_templates
type FormOf struct {
	Lang string `lang:"true" json:"lang,omitempty" firestore:"lang,omitempty"`
	Word string `json:"word,omitempty" firestore:"word,omitempty"`
	Alt  string `json:"alt,omitempty" firestore:"alt,omitempty"`

	Tag1  string `json:"tag1,omitempty" firestore:"tag1,omitempty"`
	Tag2  string `json:"tag2,omitempty" firestore:"tag2,omitempty"`
	Tag3  string `json:"tag3,omitempty" firestore:"tag3,omitempty"`
	Tag4  string `json:"tag4,omitempty" firestore:"tag4,omitempty"`
	Tag5  string `json:"tag5,omitempty" firestore:"tag5,omitempty"`
	Tag6  string `json:"tag6,omitempty" firestore:"tag6,omitempty"`
	Tag7  string `json:"tag7,omitempty" firestore:"tag7,omitempty"`
	Tag8  string `json:"tag8,omitempty" firestore:"tag8,omitempty"`
	Tag9  string `json:"tag9,omitempty" firestore:"tag9,omitempty"`
	Tag10 string `json:"tag10,omitempty" firestore:"tag10,omitempty"`

	Form string `json:"form,omityempty" firestore:"form,omitempty"`
}

// https://en.wiktionary.org/wiki/Template:form_of
type FormOfGeneric struct {
	Lang       string `lang:"true" json:"lang,omitempty" firestore:"lang,omitempty"`
	Definition string `json:"definition,omitempty" firestore:"definition,omitempty"`
	Word       string `json:"word,omitempty" firestore:"word,omitempty"`
	Alt        string `json:"alt,omitempty" firestore:"alt,omitempty"`
}

type GlossaryEntry struct {
	Type         string
	Name         string
	Shortcuts    []string
	NoLeftSpace  bool
	NoRightSpace bool
}

var formOfGlossary = map[string]*GlossaryEntry{
	// Person

	"first-person": &GlossaryEntry{
		Type:      "person",
		Shortcuts: []string{"1"},
	},

	"second-person": &GlossaryEntry{
		Type:      "person",
		Shortcuts: []string{"2"},
	},

	"third-person": &GlossaryEntry{
		Type:      "person",
		Shortcuts: []string{"3"},
	},

	"impersonal": &GlossaryEntry{
		Type:      "person",
		Shortcuts: []string{"impers"},
	},

	// Number

	"singular": &GlossaryEntry{
		Type:      "number",
		Shortcuts: []string{"s", "sg"},
	},

	"dual": &GlossaryEntry{
		Type:      "number",
		Shortcuts: []string{"d", "du"},
	},

	"plural": &GlossaryEntry{
		Type:      "number",
		Shortcuts: []string{"p", "pl"},
	},

	"single-possession": &GlossaryEntry{
		Type:      "number",
		Shortcuts: []string{"spos"},
	},

	"multiple-possession": &GlossaryEntry{
		Type:      "number",
		Shortcuts: []string{"mpos"},
	},

	// Gender

	"masculine": &GlossaryEntry{
		Type:      "gender",
		Shortcuts: []string{"m"},
	},

	"natural masculine": &GlossaryEntry{
		Type:      "gender",
		Shortcuts: []string{"natm"},
	},

	"feminine": &GlossaryEntry{
		Type:      "gender",
		Shortcuts: []string{"f"},
	},

	"neuter": &GlossaryEntry{
		Type:      "gender",
		Shortcuts: []string{"n"},
	},

	"common": &GlossaryEntry{
		Type:      "gender",
		Shortcuts: []string{"c"},
	},

	"nonvirile": &GlossaryEntry{
		Type:      "gender",
		Shortcuts: []string{"nv"},
	},

	// Animacy

	"animate": &GlossaryEntry{
		Type:      "animacy",
		Shortcuts: []string{"an"},
	},

	"inanimate": &GlossaryEntry{
		Type:      "animacy",
		Shortcuts: []string{"in", "inan"},
	},

	"personal": &GlossaryEntry{
		Type:      "animacy",
		Shortcuts: []string{"pr", "pers"},
	},

	// Tense/aspect

	"present": &GlossaryEntry{
		Type:      "tense-aspect",
		Shortcuts: []string{"pres"},
	},

	"past": &GlossaryEntry{
		Type: "tense-aspect",
	},

	"future": &GlossaryEntry{
		Type:      "tense-aspect",
		Shortcuts: []string{"fut", "futr"},
	},

	"non-past": &GlossaryEntry{
		Type:      "tense-aspect",
		Shortcuts: []string{"npast"},
	},

	"progressive": &GlossaryEntry{
		Type:      "tense-aspect",
		Shortcuts: []string{"prog"},
	},

	"preterite": &GlossaryEntry{
		Type:      "tense-aspect",
		Shortcuts: []string{"pret"},
	},

	"perfect": &GlossaryEntry{
		Type:      "tense-aspect",
		Shortcuts: []string{"perf"},
	},

	"imperfect": &GlossaryEntry{
		Type:      "tense-aspect",
		Shortcuts: []string{"impf", "imperf"},
	},

	"pluperfect": &GlossaryEntry{
		Type:      "tense-aspect",
		Shortcuts: []string{"plup", "pluperf"},
	},

	"aorist": &GlossaryEntry{
		Type:      "tense-aspect",
		Shortcuts: []string{"aor", "aori"},
	},

	"past historic": &GlossaryEntry{
		Type:      "tense-aspect",
		Shortcuts: []string{"phis"},
	},

	"imperfective": &GlossaryEntry{
		Type:      "tense-aspect",
		Shortcuts: []string{"impfv", "imperfv"},
	},

	"perfective": &GlossaryEntry{
		Type:      "tense-aspect",
		Shortcuts: []string{"pfv", "perfv"},
	},

	// Mood

	"imperative": &GlossaryEntry{
		Type:      "mood",
		Shortcuts: []string{"imp", "impr", "impv"},
	},

	"indicative": &GlossaryEntry{
		Type:      "mood",
		Shortcuts: []string{"ind", "indc", "indic"},
	},

	"subjunctive": &GlossaryEntry{
		Type:      "mood",
		Shortcuts: []string{"sub", "subj"},
	},

	"conditional": &GlossaryEntry{
		Type:      "mood",
		Shortcuts: []string{"cond"},
	},

	"optative": &GlossaryEntry{
		Type:      "mood",
		Shortcuts: []string{"opta", "opt"},
	},

	"jussive": &GlossaryEntry{
		Type:      "mood",
		Shortcuts: []string{"juss"},
	},

	// Voice / valence

	"active": &GlossaryEntry{
		Type:      "voice-valence",
		Shortcuts: []string{"act", "actv"},
	},

	"middle": &GlossaryEntry{
		Type:      "voice-valence",
		Shortcuts: []string{"mid", "midl"},
	},

	"passive": &GlossaryEntry{
		Type:      "voice-valence",
		Shortcuts: []string{"pass", "pasv"},
	},

	"mediopassive": &GlossaryEntry{
		Type:      "voice-valence",
		Shortcuts: []string{"mp", "mpass", "mpasv", "mpsv"},
	},

	"reflexive": &GlossaryEntry{
		Type:      "voice-valence",
		Shortcuts: []string{"refl"},
	},

	"transitive": &GlossaryEntry{
		Type:      "voice-valence",
		Shortcuts: []string{"tr", "vt"},
	},

	"intransitive": &GlossaryEntry{
		Type:      "voice-valence",
		Shortcuts: []string{"intr", "vi"},
	},

	"ditransitive": &GlossaryEntry{
		Type:      "voice-valence",
		Shortcuts: []string{"ditr"},
	},

	"causative": &GlossaryEntry{
		Type:      "voice-valence",
		Shortcuts: []string{"caus"},
	},

	// Non-finite

	"infinitive": &GlossaryEntry{
		Type:      "non-finite",
		Shortcuts: []string{"inf"},
	},

	"personal infinitive": &GlossaryEntry{
		Type:      "non-finite",
		Shortcuts: []string{"pinf"},
	},

	"participle": &GlossaryEntry{
		Type:      "non-finite",
		Shortcuts: []string{"part", "ptcp"},
	},

	"verbal noun": &GlossaryEntry{
		Type:      "non-finite",
		Shortcuts: []string{"vnoun"},
	},

	"gerund": &GlossaryEntry{
		Type:      "non-finite",
		Shortcuts: []string{"ger"},
	},

	"supine": &GlossaryEntry{
		Type:      "non-finite",
		Shortcuts: []string{"sup"},
	},

	"transgressive": &GlossaryEntry{
		Type: "non-finite",
	},

	// Case

	"ablative": &GlossaryEntry{
		Type:      "case",
		Shortcuts: []string{"abl"},
	},

	"accusative": &GlossaryEntry{
		Type:      "case",
		Shortcuts: []string{"acc"},
	},

	"dative": &GlossaryEntry{
		Type:      "case",
		Shortcuts: []string{"dat"},
	},

	"genitive": &GlossaryEntry{
		Type:      "case",
		Shortcuts: []string{"gen"},
	},

	"instrumental": &GlossaryEntry{
		Type:      "case",
		Shortcuts: []string{"ins"},
	},

	"locative": &GlossaryEntry{
		Type:      "case",
		Shortcuts: []string{"loc"},
	},

	"nominative": &GlossaryEntry{
		Type:      "case",
		Shortcuts: []string{"nom"},
	},

	"prepositional": &GlossaryEntry{
		Type:      "case",
		Shortcuts: []string{"pre", "prep"},
	},

	"vocative": &GlossaryEntry{
		Type:      "case",
		Shortcuts: []string{"voc"},
	},

	// State

	"construct": &GlossaryEntry{
		Type:      "state",
		Shortcuts: []string{"cons", "construct state"},
	},

	"definite": &GlossaryEntry{
		Type:      "state",
		Shortcuts: []string{"def", "defn", "definite state"},
	},

	"indefinite": &GlossaryEntry{
		Type:      "state",
		Shortcuts: []string{"indef", "indf", "indefinite state"},
	},

	"strong": &GlossaryEntry{
		Type:      "state",
		Shortcuts: []string{"str"},
	},

	"weak": &GlossaryEntry{
		Type:      "state",
		Shortcuts: []string{"wk"},
	},

	"mixed": &GlossaryEntry{
		Type:      "state",
		Shortcuts: []string{"mix"},
	},

	"attributive": &GlossaryEntry{
		Type:      "state",
		Shortcuts: []string{"attr"},
	},

	"predicative": &GlossaryEntry{
		Type:      "state",
		Shortcuts: []string{"pred"},
	},

	// Degrees of comparison

	"positive degree": &GlossaryEntry{
		Type:      "comparison",
		Shortcuts: []string{"posd", "positive"},
	},

	"comparative degree": &GlossaryEntry{
		Type:      "comparison",
		Shortcuts: []string{"comd", "comparative"},
	},

	"superlative degree": &GlossaryEntry{
		Type:      "comparison",
		Shortcuts: []string{"supd", "superlative"},
	},

	// Inflectional class

	"pronominal": &GlossaryEntry{
		Type:      "class",
		Shortcuts: []string{"pron"},
	},

	// Attitude

	"augmentative": &GlossaryEntry{
		Type:      "attitude",
		Shortcuts: []string{"aug"},
	},

	"diminutive": &GlossaryEntry{
		Type:      "attitude",
		Shortcuts: []string{"dim"},
	},

	"pejorative": &GlossaryEntry{
		Type:      "attitude",
		Shortcuts: []string{"pej"},
	},

	// Sound changes

	"contracted": &GlossaryEntry{
		Type: "sound change",
	},

	// Misc grammar

	"simple": &GlossaryEntry{
		Type:      "grammar",
		Shortcuts: []string{"sim"},
	},

	"short": &GlossaryEntry{
		Type: "grammar",
	},

	"long": &GlossaryEntry{
		Type: "grammar",
	},

	"form": &GlossaryEntry{
		Type: "grammar",
	},

	"adjectival": &GlossaryEntry{
		Type:      "grammar",
		Shortcuts: []string{"adj"},
	},

	"adverbial": &GlossaryEntry{
		Type:      "grammar",
		Shortcuts: []string{"adv"},
	},

	"negative": &GlossaryEntry{
		Type:      "grammar",
		Shortcuts: []string{"neg"},
	},

	"possessive": &GlossaryEntry{
		Type:      "non-finite",
		Shortcuts: []string{"poss"},
	},

	"nominalized": &GlossaryEntry{
		Type:      "grammar",
		Shortcuts: []string{"nomz"},
	},

	"nominalization": &GlossaryEntry{
		Type:      "grammar",
		Shortcuts: []string{"nomzn"},
	},

	"root": &GlossaryEntry{
		Type: "grammar",
	},

	"stem": &GlossaryEntry{
		Type: "grammar",
	},

	"dependent": &GlossaryEntry{
		Type:      "grammar",
		Shortcuts: []string{"dep"},
	},

	"independent": &GlossaryEntry{
		Type:      "grammar",
		Shortcuts: []string{"indep"},
	},

	// Other tags

	"and": &GlossaryEntry{
		Type: "other",
	},

	",": &GlossaryEntry{
		Type:        "other",
		NoLeftSpace: true,
	},

	":": &GlossaryEntry{
		Type:        "other",
		NoLeftSpace: true,
	},

	"/": &GlossaryEntry{
		Type:         "other",
		NoLeftSpace:  true,
		NoRightSpace: true,
	},

	"(": &GlossaryEntry{
		Type:         "other",
		NoRightSpace: true,
	},

	")": &GlossaryEntry{
		Type:        "other",
		NoLeftSpace: true,
	},

	"[": {
		Type:         "other",
		NoRightSpace: true,
	},

	"]": {
		Type:        "other",
		NoLeftSpace: true,
	},

	"-": {
		Type:         "other",
		NoLeftSpace:  true,
		NoRightSpace: true,
	},
}

var formOfShortcuts = map[string][]string{
	// Person

	"12":  {"1//2"},
	"13":  {"1//3"},
	"23":  {"2//3"},
	"123": {"1//2//3"},

	// Number

	"1s": {"1", "s"},
	"2s": {"2", "s"},
	"3s": {"3", "s"},
	"1d": {"1", "d"},
	"2d": {"2", "d"},
	"3d": {"3", "d"},
	"1p": {"1", "p"},
	"2p": {"2", "p"},
	"3p": {"3", "p"},

	// Gender

	"mf":  {"m//f"},
	"mn":  {"m//n"},
	"fn":  {"f//n"},
	"mfn": {"m//f//n"},

	// Tense / aspect

	"spast":          {"simple", "past"},
	"simple past":    {"simple", "past"},
	"spres":          {"simple", "present"},
	"simple present": {"simple", "present"},
}

func init() {
	for name := range formOfGlossary {
		formOfGlossary[name].Name = name
		for _, shortcut := range formOfGlossary[name].Shortcuts {
			formOfShortcuts[shortcut] = []string{name}
		}
	}
}

func (tpl *Template) ToFormOf(form string, tags ...string) FormOf {
	fo := FormOf{}

	ts := len(tags)
	if ts > 0 {
		fo.Tag1 = tags[0]
	}
	if ts > 1 {
		fo.Tag2 = tags[1]
	}
	if ts > 2 {
		fo.Tag3 = tags[2]
	}
	if ts > 3 {
		fo.Tag4 = tags[3]
	}
	if ts > 4 {
		fo.Tag5 = tags[4]
	}
	if ts > 5 {
		fo.Tag6 = tags[5]
	}
	if ts > 6 {
		fo.Tag7 = tags[6]
	}
	if ts > 7 {
		fo.Tag8 = tags[7]
	}
	if ts > 8 {
		fo.Tag9 = tags[8]
	}
	if ts > 9 {
		fo.Tag10 = tags[9]
	}

	tpl.toConcrete(reflect.TypeOf(fo), reflect.ValueOf(&fo))
	fo.Word = toEntryName(fo.Lang, fo.Word)
	fo.Form = form
	return fo
}

func (fo *FormOf) DisplayWord() string {
	if fo.Alt != "" {
		return fo.Alt
	}
	return fo.Word
}

func (fo *FormOf) Text() string {
	desc := GetFormOfDescription(
		fo.Form, fo.Tag1, fo.Tag2, fo.Tag3, fo.Tag4, fo.Tag5, fo.Tag6, fo.Tag7, fo.Tag8, fo.Tag9, fo.Tag10,
	)
	return fmt.Sprintf("%s %s", desc, fo.DisplayWord())
}

func GetFormOfDescription(form string, tags ...string) string {
	expandedTags := expandFormOfTags(tags)

	var entries []*GlossaryEntry
	for _, tag := range expandedTags {
		if entry, ok := formOfGlossary[tag]; ok {
			entries = append(entries, entry)
		}
	}

	if len(entries) == 0 {
		return form
	}

	parts := []string{}
	total := len(entries)
	for i, entry := range entries {
		var space string
		if (i != total-1) && !entry.NoRightSpace && !entries[i+1].NoLeftSpace {
			space = " "
		}
		parts = append(parts, fmt.Sprintf("%s%s", entry.Name, space))
	}

	return fmt.Sprintf("%s of", strings.Join(parts, ""))
}

func expandFormOfTags(tags []string) []string {
	var expanded []string
	for _, tag := range tags {
		if _, ok := formOfGlossary[tag]; ok {
			expanded = append(expanded, tag)
		} else if tags, ok := formOfShortcuts[tag]; ok {
			expanded = append(expanded, expandFormOfTags(tags)...)
		}
	}
	return expanded
}

func (tpl *Template) ToFormOfGeneric() FormOfGeneric {
	fog := FormOfGeneric{}
	tpl.toConcrete(reflect.TypeOf(fog), reflect.ValueOf(&fog))
	fog.Word = toEntryName(fog.Lang, fog.Word)
	return fog
}

func (fog *FormOfGeneric) DisplayWord() string {
	if fog.Alt != "" {
		return fog.Alt
	}
	return fog.Word
}

func (fog *FormOfGeneric) Text() string {
	return fmt.Sprintf("%s of %s", fog.Definition, fog.DisplayWord())
}
