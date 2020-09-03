package tpl

import (
	"fmt"
	"reflect"
	"strings"
)

// https://en.wiktionary.org/wiki/Template:inflection_of
type Inflection struct {
	Lang  string `lang:"true" json:"lang,omitempty" firestore:"lang,omitempty"`
	Lemma string `json:"lemma,omitempty" firestore:"lemma,omitempty"`
	Alt   string `json:"alt,omitempty" firestore:"alt,omitempty"`

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
}

type GlossaryEntry struct {
	Type         string
	Name         string
	Shortcuts    []string
	NoLeftSpace  bool
	NoRightSpace bool
}

// See:
//   https://en.wiktionary.org/wiki/Module:form_of/templates#inflection_of_t
//   https://en.wiktionary.org/wiki/Module:form_of/data
// Not yet implemented:
//   https://en.wiktionary.org/wiki/Module:form_of/data2
var glossary = map[string]*GlossaryEntry{
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
		Type:      "tense-aspect",
		Shortcuts: []string{},
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
		Type: "mood",
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

	// Voice/valence

	"active": &GlossaryEntry{
		Type:      "voice-valence",
		Shortcuts: []string{"act", "actv"},
	},

	"middle": &GlossaryEntry{
		Type:      "voice-valence",
		Shortcuts: []string{"mid", "midl"},
	},

	"passive": &GlossaryEntry{
		Type: "voice-valence",
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
		Type: "voice-valence",
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
		Type: "non-finite",
	},

	"supine": &GlossaryEntry{
		Type:      "non-finite",
		Shortcuts: []string{"sup"},
	},

	"transgressive": &GlossaryEntry{
		Type:      "non-finite",
		Shortcuts: []string{},
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
		Type: "case",
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
		Type: "state",
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

	// Register

	// Deixis

	// Clusivity

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
		Type:      "sound change",
		Shortcuts: []string{},
	},

	// Misc grammar

	"simple": &GlossaryEntry{
		Type:      "grammar",
		Shortcuts: []string{"sim"},
	},

	"short": &GlossaryEntry{
		Type:      "grammar",
		Shortcuts: []string{},
	},
	"long": &GlossaryEntry{
		Type:      "grammar",
		Shortcuts: []string{},
	},

	"form": &GlossaryEntry{
		Type:      "grammar",
		Shortcuts: []string{},
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
		Type:      "grammar",
		Shortcuts: []string{},
	},

	"stem": &GlossaryEntry{
		Type:      "grammar",
		Shortcuts: []string{},
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
		Type:      "other",
		Shortcuts: []string{},
	},

	",": &GlossaryEntry{
		Type:        "other",
		Shortcuts:   []string{},
		NoLeftSpace: true,
	},

	":": &GlossaryEntry{
		Type:        "other",
		Shortcuts:   []string{},
		NoLeftSpace: true,
	},

	"/": &GlossaryEntry{
		Type:         "other",
		Shortcuts:    []string{},
		NoLeftSpace:  true,
		NoRightSpace: true,
	},

	"(": &GlossaryEntry{
		Type:         "other",
		Shortcuts:    []string{},
		NoRightSpace: true,
	},

	")": &GlossaryEntry{
		Type:        "other",
		Shortcuts:   []string{},
		NoLeftSpace: true,
	},

	"[": &GlossaryEntry{
		Type:         "other",
		Shortcuts:    []string{},
		NoRightSpace: true,
	},

	"]": &GlossaryEntry{
		Type:        "other",
		Shortcuts:   []string{},
		NoLeftSpace: true,
	},

	"-": &GlossaryEntry{
		Type:         "other",
		Shortcuts:    []string{},
		NoLeftSpace:  true,
		NoRightSpace: true,
	},
}

var shortcuts = map[string][]string{
	"12":  {"1", "/", "2"},
	"13":  {"1", "/", "3"},
	"23":  {"2", "/", "3"},
	"123": {"1", "/", "2", "/", "3"},

	"1s": {"1", "s"},
	"2s": {"2", "s"},
	"3s": {"3", "s"},
	"1d": {"1", "d"},
	"2d": {"2", "d"},
	"3d": {"3", "d"},
	"1p": {"1", "p"},
	"2p": {"2", "p"},
	"3p": {"3", "p"},

	"mf":  {"m", "/", "f"},
	"mn":  {"m", "/", "n"},
	"fn":  {"f", "/", "n"},
	"mfn": {"m", "/", "f", "/", "n"},

	"spast":          {"simple", "past"},
	"simple past":    {"simple", "past"},
	"spres":          {"simple", "present"},
	"simple present": {"simple", "present"},
}

func init() {
	for name := range glossary {
		glossary[name].Name = name
		for _, shortcut := range glossary[name].Shortcuts {
			shortcuts[shortcut] = []string{name}
		}
	}
}

func (tpl *Template) ToInflection() Inflection {
	i := Inflection{}
	tpl.toConcrete(reflect.TypeOf(i), reflect.ValueOf(&i))
	i.Lemma = toEntryName(i.Lang, i.Lemma)
	return i
}

func (i *Inflection) Text() string {
	desc := GetDescription([]string{
		i.Tag1, i.Tag2, i.Tag3, i.Tag4, i.Tag5, i.Tag6, i.Tag7, i.Tag8, i.Tag9, i.Tag10,
	})
	return fmt.Sprintf("%s of %s", desc, i.Lemma)
}

func GetDescription(tags []string) string {
	expandedTags := expandTags(tags)

	var entries []*GlossaryEntry
	for _, tag := range expandedTags {
		if entry, ok := glossary[tag]; ok {
			entries = append(entries, entry)
		}
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

	return strings.Join(parts, "")
}

func expandTags(tags []string) []string {
	var expanded []string
	for _, tag := range tags {
		if _, ok := glossary[tag]; ok {
			expanded = append(expanded, tag)
		} else if tags, ok := shortcuts[tag]; ok {
			expanded = append(expanded, expandTags(tags)...)
		}
	}
	return expanded
}
