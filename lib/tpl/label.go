package tpl

import (
	"fmt"
	"reflect"
	"strings"
)

// https://en.wiktionary.org/wiki/Template:label
type Label struct {
	Lang string `lang:"true" json:"lang,omitempty" firestore:"lang,omitempty"`

	Label1  string `json:"label1,omitempty" firestore:"label1,omitempty"`
	Label2  string `json:"label2,omitempty" firestore:"label2,omitempty"`
	Label3  string `json:"label3,omitempty" firestore:"label3,omitempty"`
	Label4  string `json:"label4,omitempty" firestore:"label4,omitempty"`
	Label5  string `json:"label5,omitempty" firestore:"label5,omitempty"`
	Label6  string `json:"label6,omitempty" firestore:"label6,omitempty"`
	Label7  string `json:"label7,omitempty" firestore:"label7,omitempty"`
	Label8  string `json:"label8,omitempty" firestore:"label8,omitempty"`
	Label9  string `json:"label9,omitempty" firestore:"label9,omitempty"`
	Label10 string `json:"label10,omitempty" firestore:"label10,omitempty"`
}

type LabelEntry struct {
	Markup       string
	Text         string
	NoLeftComma  bool
	NoRightComma bool
	NoLeftSpace  bool
	NoText       bool
	Alias        string
}

var labelMap = map[string]*LabelEntry{
	"_": &LabelEntry{
		Markup:       "",
		NoLeftComma:  true,
		NoRightComma: true,
		NoLeftSpace:  true,
		NoText:       true,
	},

	"also": &LabelEntry{
		NoRightComma: true,
	},

	"and": &LabelEntry{
		NoLeftComma:  true,
		NoRightComma: true,
	},
	"&": &LabelEntry{Alias: "and"},

	"or": &LabelEntry{
		NoLeftComma:  true,
		NoRightComma: true,
	},

	";": &LabelEntry{
		NoLeftComma:  true,
		NoRightComma: true,
		NoLeftSpace:  true,
	},

	"by": &LabelEntry{
		NoLeftComma:  true,
		NoRightComma: true,
	},

	"with": &LabelEntry{
		NoLeftComma:  true,
		NoRightComma: true,
	},
	"+": &LabelEntry{Alias: "with"},

	"except": &LabelEntry{
		NoLeftComma:  true,
		NoRightComma: true,
	},

	"outside": &LabelEntry{
		NoLeftComma:  true,
		NoRightComma: true,
	},
	"except in": &LabelEntry{Alias: "outside"},

	// Qualifier labels

	"chiefly": &LabelEntry{
		NoRightComma: true,
	},
	"mainly":    &LabelEntry{Alias: "chiefly"},
	"mostly":    &LabelEntry{Alias: "chiefly"},
	"primarily": &LabelEntry{Alias: "chiefly"},

	"especially": &LabelEntry{
		NoRightComma: true,
	},

	"particularly": &LabelEntry{
		NoRightComma: true,
	},

	"excluding": &LabelEntry{
		NoRightComma: true,
	},

	"extremely": &LabelEntry{
		NoRightComma: true,
	},

	"frequently": &LabelEntry{
		NoRightComma: true,
	},

	"humorously": &LabelEntry{
		NoRightComma: true,
	},

	"including": &LabelEntry{
		NoRightComma: true,
	},

	"many": &LabelEntry{
		NoRightComma: true,
	},

	"markedly": &LabelEntry{
		NoRightComma: true,
	},

	"mildly": &LabelEntry{
		NoRightComma: true,
	},

	"now": &LabelEntry{
		NoRightComma: true,
	},
	"nowadays": &LabelEntry{Alias: "now"},
	"Now":      &LabelEntry{Alias: "now"},

	"of": &LabelEntry{
		NoRightComma: true,
	},

	"of a": &LabelEntry{
		NoRightComma: true,
	},

	"of an": &LabelEntry{
		NoRightComma: true,
	},

	"often": &LabelEntry{
		NoRightComma: true,
	},

	"originally": &LabelEntry{
		NoRightComma: true,
	},

	"possibly": &LabelEntry{
		NoRightComma: true,
	},

	"rarely": &LabelEntry{
		NoRightComma: true,
	},

	"slightly": &LabelEntry{
		NoRightComma: true,
	},

	"sometimes": &LabelEntry{
		NoRightComma: true,
	},

	"somewhat": &LabelEntry{
		NoRightComma: true,
	},

	"strongly": &LabelEntry{
		NoRightComma: true,
	},

	"then": &LabelEntry{
		NoRightComma: true,
	},

	"typically": &LabelEntry{
		NoRightComma: true,
	},

	"usually": &LabelEntry{
		NoRightComma: true,
	},

	"very": &LabelEntry{
		NoRightComma: true,
	},

	// Grammatical labels

	"abbreviation": &LabelEntry{
		Markup: "[[abbreviation]]",
	},

	"acronym": &LabelEntry{
		Markup: "[[acronym]]",
	},

	"active voice":  &LabelEntry{Alias: "active"},
	"in the active": &LabelEntry{Alias: "active"},

	"ambitransitive": &LabelEntry{
		Markup: "[[transitive]], [[intransitive]]",
		Text:   "transitive, intransitive",
	},

	"in the indicative": &LabelEntry{Alias: "indicative"},
	"indicative mood":   &LabelEntry{Alias: "indicative"},

	"in the subjunctive": &LabelEntry{Alias: "subjunctive"},
	"subjunctive mood":   &LabelEntry{Alias: "subjunctive"},

	"in the imperative": &LabelEntry{Alias: "imperative"},
	"imperative mood":   &LabelEntry{Alias: "imperative"},

	"in the jussive": &LabelEntry{Alias: "jussive"},
	"jussive mood":   &LabelEntry{Alias: "jussive"},

	"attributive": &LabelEntry{
		Markup: "[[Appendix:English nouns#Attributive|attributive]]",
	},

	"attributively": &LabelEntry{
		Markup: "[[Appendix:English nouns#Attributive|attributively]]",
	},

	"cardinal": &LabelEntry{
		Markup: "[[cardinal number|cardinal]]",
	},

	"causative": &LabelEntry{
		Markup: "[[causative]]",
	},

	"cognate object": &LabelEntry{
		Markup: "with [[w:Cognate object|cognate object]]",
		Text:   "with cognate object",
	},
	"with cognate object": &LabelEntry{Alias: "cognate object"},

	"control": &LabelEntry{Alias: "control verb"},

	"copulative": &LabelEntry{
		Markup: "[[copular verb|copulative]]",
	},
	"copular": &LabelEntry{Alias: "copulative"},

	"dysphemism": &LabelEntry{Alias: "dysphemistic"},

	"by ellipsis": &LabelEntry{
		Markup: "by [[ellipsis]]",
	},

	"hence": &LabelEntry{Alias: "by extension"},

	"hedges": &LabelEntry{Alias: "hedge"},

	"ideophone": &LabelEntry{Alias: "ideophonic"},

	"idiom":         &LabelEntry{Alias: "idiomatic"},
	"idiomatically": &LabelEntry{Alias: "idiomatic"},

	"in the singular": &LabelEntry{
		Markup: "in the [[singular]]",
	},
	"in singular": &LabelEntry{Alias: "in the singular"},
	"singular":    &LabelEntry{Alias: "in the singular"},

	"in the dual": &LabelEntry{
		Markup: "in the [[dual]]",
	},
	"in dual": &LabelEntry{Alias: "in the dual"},
	"dual":    &LabelEntry{Alias: "in the dual"},

	"in the plural": &LabelEntry{
		Markup: "in the [[Appendix:Glossary#plural|plural]]",
	},
	"in plural": &LabelEntry{Alias: "in the plural"},
	"plural":    &LabelEntry{Alias: "in the plural"},

	"in the mediopassive": &LabelEntry{
		Markup: "in the [[mediopassive]]",
	},
	"in mediopassive": &LabelEntry{Alias: "in the mediopassive"},
	"mediopassive":    &LabelEntry{Alias: "in the mediopassive"},

	"indef": &LabelEntry{Alias: "indefinite"},

	"initialism": &LabelEntry{
		Markup: "[[initialism]]",
	},

	"International Phonetic Alphabet": &LabelEntry{Alias: "IPA"},

	"middle voice":        &LabelEntry{Alias: "middle"},
	"in the middle":       &LabelEntry{Alias: "middle"},
	"in the middle voice": &LabelEntry{Alias: "middle"},

	"mnemonic": &LabelEntry{
		Markup: "[[mnemonic]]",
	},

	"negative polarity":       &LabelEntry{Alias: "chiefly in the negative"},
	"negative polarity item":  &LabelEntry{Alias: "chiefly in the negative"},
	"usually in the negative": &LabelEntry{Alias: "chiefly in the negative"},

	"not comparable": &LabelEntry{
		Markup: "[[Appendix:Glossary#uncomparable|not comparable]]",
	},
	"notcomp":      &LabelEntry{Alias: "not comparable"},
	"uncomparable": &LabelEntry{Alias: "not comparable"},

	"onomatopoeia": &LabelEntry{
		Markup: "[[onomatopoeia]]",
	},

	"passive voice":  &LabelEntry{Alias: "passive"},
	"in the passive": &LabelEntry{Alias: "passive"},

	"pluralonly":     &LabelEntry{Alias: "plural only"},
	"plurale tantum": &LabelEntry{Alias: "plural only"},

	"possessive pronoun": &LabelEntry{
		Text: "possessive",
	},

	"predicative": &LabelEntry{
		Markup: "[[Appendix:Glossary#predicative|predicative]]",
	},

	"predicatively": &LabelEntry{
		Markup: "[[Appendix:Glossary#predicative|predicatively]]",
	},

	"productive": &LabelEntry{
		Markup: "[[productive]]"},

	"pronominal": &LabelEntry{
		Markup: "takes a [[Appendix:Glossary#reflexive|reflexive pronoun]]",
		Text:   "takes a reflexive pronoun",
	},

	"reciprocal": &LabelEntry{
		Markup: "[[Appendix:Glossary#reciprocal|reciprocal]]",
	},

	"reflexive": &LabelEntry{
		Markup: "[[Appendix:Glossary#reflexive|reflexive]]",
	},

	"reflexive pronoun": &LabelEntry{
		Markup: "[[Appendix:Glossary#reflexive|reflexive]]",
		Text:   "reflexive",
	},

	"relational": &LabelEntry{
		Markup: "[[Appendix:Glossary#relational|relational]]",
	},

	"set phrase": &LabelEntry{
		Markup: "[[set phrase]]"},

	"singular only": &LabelEntry{
		Markup: "singular only",
	},
	"singulare tantum": &LabelEntry{Alias: "singular only"},
	"no plural":        &LabelEntry{Alias: "singular only"},

	"stative verb": &LabelEntry{Alias: "stative"},

	"narrowly": &LabelEntry{Alias: "strictly"},

	"usually plural": &LabelEntry{
		Markup: "usually in the [[plural]]",
		Text:   "usually in the plural",
	},
	"usually in the plural": &LabelEntry{Alias: "usually plural"},
	"usually in plural":     &LabelEntry{Alias: "usually plural"},

	// Usage labels

	"ACG": &LabelEntry{
		Markup: "[[ACG]]",
	},

	"ad slang": &LabelEntry{Alias: "advertising slang"},
	"cosmo":    &LabelEntry{Alias: "advertising slang"},

	"endearing": &LabelEntry{
		Markup: "[[endearing]]",
	},
	"affectionate": &LabelEntry{Alias: "endearing"},

	"pre-classical": &LabelEntry{
		Text: "pre-Classical",
	},
	"Pre-classical":  &LabelEntry{Alias: "pre-classical"},
	"pre-Classical":  &LabelEntry{Alias: "pre-classical"},
	"Pre-Classical":  &LabelEntry{Alias: "pre-classical"},
	"Preclassical":   &LabelEntry{Alias: "pre-classical"},
	"preclassical":   &LabelEntry{Alias: "pre-classical"},
	"ante-classical": &LabelEntry{Alias: "pre-classical"},
	"Ante-classical": &LabelEntry{Alias: "pre-classical"},
	"ante-Classical": &LabelEntry{Alias: "pre-classical"},
	"Ante-Classical": &LabelEntry{Alias: "pre-classical"},
	"Anteclassical":  &LabelEntry{Alias: "pre-classical"},
	"anteclassical":  &LabelEntry{Alias: "pre-classical"},

	"back slang": &LabelEntry{
		Markup: "[[Appendix:Glossary#backslang|back slang]]",
	},
	"backslang":  &LabelEntry{Alias: "back slang"},
	"back-slang": &LabelEntry{Alias: "back slang"},

	"UK slang": &LabelEntry{Alias: "British slang"},

	"buzzword": &LabelEntry{
		Markup: "[[buzzword]]",
	},

	"cant": &LabelEntry{
		Markup: "[[cant]]",
	},
	"argot":      &LabelEntry{Alias: "cant"},
	"cryptolect": &LabelEntry{Alias: "cant"},

	"capitalized": &LabelEntry{
		Markup: "[[capitalisation|capitalized]]"},

	"childish": &LabelEntry{
		Markup: "[[childish]]",
	},
	"baby talk":      &LabelEntry{Alias: "childish"},
	"child language": &LabelEntry{Alias: "childish"},

	"chu Nom": &LabelEntry{
		Markup: "[[Vietnamese]] [[chữ Nôm]]",
		Text:   "Vietnamese chữ Nôm",
	},

	"Classic 1811 Dictionary of the i Tongue": &LabelEntry{
		Markup: "[[Appendix:Glossary#archaic|archaic]], [[Appendix:Glossary#slang|slang]]",
		Text:   "archaic, slang",
	},
	"1811": &LabelEntry{Alias: "Classic 1811 Dictionary of the Vulgar Tongue"},

	"Cockney rhyming slang": &LabelEntry{
		Markup: "[[Cockney rhyming slang]]",
	},

	"colloquially": &LabelEntry{Alias: "colloquial"},

	"costermongers": &LabelEntry{
		Markup: "[[Appendix:Costermongers' back slang|costermongers]]",
	},
	"coster":                    &LabelEntry{Alias: "costermongers"},
	"costers":                   &LabelEntry{Alias: "costermongers"},
	"costermonger":              &LabelEntry{Alias: "costermongers"},
	"costermongers back slang":  &LabelEntry{Alias: "costermongers"},
	"costermongers' back slang": &LabelEntry{Alias: "costermongers"},

	"derogatory": &LabelEntry{
		Markup: "[[derogatory]]",
	},
	"pejorative":  &LabelEntry{Alias: "derogatory"},
	"disparaging": &LabelEntry{Alias: "derogatory"},

	"dialect": &LabelEntry{
		Markup: "[[Appendix:Glossary#dialectal|dialect]]",
	},

	"dialects": &LabelEntry{
		Markup: "[[Appendix:Glossary#dialectal|dialects]]",
	},

	"dismissal": &LabelEntry{
		Markup: "[[dismissal]]",
	},

	"elevated": &LabelEntry{Alias: "solemn"},

	"ethnic slur": &LabelEntry{
		Markup: "[[ethnic]] [[slur]]",
	},
	"racial slur": &LabelEntry{Alias: "ethnic slur"},

	"euphemism": &LabelEntry{Alias: "euphemistic"},

	"eye dialect": &LabelEntry{
		Markup: "[[eye dialect]]",
	},

	"fandom slang": &LabelEntry{
		Markup: "[[fandom]] [[slang]]",
	},
	"fandom": &LabelEntry{Alias: "fandom slang"},

	"figuratively":   &LabelEntry{Alias: "figurative"},
	"metaphorically": &LabelEntry{Alias: "figurative"},
	"metaphorical":   &LabelEntry{Alias: "figurative"},
	"metaphor":       &LabelEntry{Alias: "figurative"},

	"gay slang": &LabelEntry{
		Markup: "[[gay]] [[slang]]",
	},

	"hapax legomenon": &LabelEntry{
		Markup: "hapax",
		Text:   "hapax",
	},
	"hapax": &LabelEntry{Alias: "hapax"},

	"historic": &LabelEntry{Alias: "historical"},
	"history":  &LabelEntry{Alias: "historical"},

	"non-native speakers": &LabelEntry{
		Markup: "[[non-native speaker]]s",
	},
	"NNS": &LabelEntry{Alias: "non-native speakers"},

	"non-native speakers' English": &LabelEntry{
		Markup: "[[non-native speaker]]s' English",
	},
	"NNES": &LabelEntry{Alias: "non-native speakers' English"},
	"NNSE": &LabelEntry{Alias: "non-native speakers' English"},

	"Homeric epithet": &LabelEntry{
		Markup:       "[[Homeric Greek|Homeric]] [[w:Homeric epithets|epithet]]",
		NoRightComma: true,
	},

	"humble": &LabelEntry{
		Markup: "[[humble]]",
	},

	"humorous": &LabelEntry{
		Markup: "[[humorous]]",
	},
	"jocular": &LabelEntry{Alias: "humorous"},

	"hyperbole": &LabelEntry{Alias: "hyperbolic"},

	"informally": &LabelEntry{Alias: "informal"},

	"Internet slang": &LabelEntry{
		Markup: "[[Internet]] [[slang]]",
	},

	"internet slang": &LabelEntry{Alias: "Internet slang"},

	"IRC": &LabelEntry{
		Markup: "[[IRC]]",
	},

	"leet": &LabelEntry{
		Markup: "[[leetspeak]]",
		Text:   "leetspeak",
	},
	"leetspeak": &LabelEntry{Alias: "leet"},

	"literal": &LabelEntry{Alias: "literally"},

	"bookish": &LabelEntry{Alias: "literary"},

	"Lubunyaca": &LabelEntry{
		Markup: "[[Lubunyaca]]",
	},

	"medical slang": &LabelEntry{
		Markup: "[[medical]] [[slang]]",
	},

	"male speech": &LabelEntry{Alias: "men's speech"},

	"metonymic": &LabelEntry{Alias: "metonymically"},
	"metonymy":  &LabelEntry{Alias: "metonymically"},
	"metonym":   &LabelEntry{Alias: "metonymically"},

	"military slang": &LabelEntry{
		Markup: "[[military]] [[slang]]",
	},

	"minced oath": &LabelEntry{
		Markup: "[[minced oath]]",
	},

	"neologistic": &LabelEntry{Alias: "neologism"},

	"neopronoun": &LabelEntry{
		Markup: "[[neopronoun]]",
	},

	"no longer productive": &LabelEntry{
		Markup: "no longer [[Appendix:Glossary#productive|productive]]",
	},

	"nonce word": &LabelEntry{
		Markup: "[[Appendix:Glossary#nonce word|nonce word]]",
	},
	"nonce": &LabelEntry{Alias: "nonce word"},

	"non-standard": &LabelEntry{Alias: "nonstandard"},

	"offensive": &LabelEntry{
		Markup: "[[offensive]]",
	},

	"officialese": &LabelEntry{
		Markup: "[[officialese]]",
	},

	"Oxbridge slang": &LabelEntry{
		Markup: "[[Oxbridge]] [[slang]]",
	},

	"poetic": &LabelEntry{
		Markup: "[[poetic]]",
	},

	"Polari": &LabelEntry{
		Markup: "[[Polari]]",
	},

	"post-classical": &LabelEntry{
		Text: "post-Classical",
	},
	"Post-classical": &LabelEntry{Alias: "post-classical"},
	"post-Classical": &LabelEntry{Alias: "post-classical"},
	"Post-Classical": &LabelEntry{Alias: "post-classical"},
	"Postclassical":  &LabelEntry{Alias: "post-classical"},
	"postclassical":  &LabelEntry{Alias: "post-classical"},

	"radio slang": &LabelEntry{
		Markup: "[[radio]] [[slang]]",
	},

	"rare sense": &LabelEntry{Alias: "rare"},

	"rare term": &LabelEntry{
		Text: "rare",
	},

	"religious slur": &LabelEntry{
		Markup: "[[religious]] [[slur]]",
	},
	"sectarian slur": &LabelEntry{Alias: "religious slur"},

	"sarcastic": &LabelEntry{
		Markup: "[[sarcastic]]",
	},

	"school slang": &LabelEntry{
		Markup: "[[school]] [[slang]]",
	},
	"public school slang": &LabelEntry{Alias: "school slang"},

	"self-deprecatory": &LabelEntry{
		Markup: "[[self-deprecatory]]",
	},
	"self-deprecating": &LabelEntry{Alias: "self-deprecatory"},

	"college slang": &LabelEntry{
		Markup: "[[college]] [[slang]]",
	},
	"university slang": &LabelEntry{Alias: "college slang"},
	"student slang":    &LabelEntry{Alias: "college slang"},

	"profanity": &LabelEntry{Alias: "swear word"},
	"expletive": &LabelEntry{Alias: "swear word"},

	"text messaging": &LabelEntry{
		Markup: "[[text messaging]]",
	},
	"texting": &LabelEntry{Alias: "text messaging"},

	"thieves cant": &LabelEntry{Alias: "thieves' cant"},
	"thieves'":     &LabelEntry{Alias: "thieves' cant"},
	"thieves":      &LabelEntry{Alias: "thieves' cant"},

	"trademark": &LabelEntry{
		Markup: "[[trademark]]",
	},

	"transferred senses": &LabelEntry{
		Markup: "[[transferred sense#English|transferred senses]]",
	},

	"transgender slang": &LabelEntry{
		Markup: "[[transgender]] [[slang]]",
	},

	"Twitch-speak": &LabelEntry{
		Markup: "[[Twitch-speak]]",
	},

	"uds.": &LabelEntry{
		Markup: "[[Appendix:Spanish pronouns#Ustedes and vosotros|used formally in Spain]]",
		Text:   "used formally in Spain",
	},

	"verlan": &LabelEntry{
		Markup: "[[Appendix:Glossary#verlan|verlan]]",
	},

	"coarse":  &LabelEntry{Alias: "vulgar"},
	"obscene": &LabelEntry{Alias: "vulgar"},
	"profane": &LabelEntry{Alias: "vulgar"},

	"2channel slang": &LabelEntry{
		Markup: "[[w:2channel|2channel]] [[slang]]",
	},

	"2ch slang": &LabelEntry{Alias: "2channel slang"},

	"female speech": &LabelEntry{Alias: "women's speech"},
}

func (tpl *Template) ToLabel() Label {
	l := Label{}
	tpl.toConcrete(reflect.TypeOf(l), reflect.ValueOf(&l))
	return l
}

func (l *Label) Text() string {
	var labels []string
	if l.Label1 != "" {
		labels = append(labels, l.Label1)
	}
	if l.Label2 != "" {
		labels = append(labels, l.Label2)
	}
	if l.Label3 != "" {
		labels = append(labels, l.Label3)
	}
	if l.Label4 != "" {
		labels = append(labels, l.Label4)
	}
	if l.Label5 != "" {
		labels = append(labels, l.Label5)
	}
	if l.Label6 != "" {
		labels = append(labels, l.Label6)
	}
	if l.Label7 != "" {
		labels = append(labels, l.Label7)
	}
	if l.Label8 != "" {
		labels = append(labels, l.Label8)
	}
	if l.Label9 != "" {
		labels = append(labels, l.Label9)
	}
	if l.Label10 != "" {
		labels = append(labels, l.Label10)
	}
	return fmt.Sprintf("(%s)", strings.Join(expandLabels(labels), ""))
}

func expandLabels(labels []string) []string {
	var entries []*LabelEntry

	// Create entries by following aliases
	for _, label := range labels {
		var entry *LabelEntry
		if e, ok := labelMap[label]; ok {
			entry = e
			if e.Alias != "" {
				if ea, ok := labelMap[e.Alias]; ok {
					entry = ea
					if entry.Text == "" {
						entry.Text = e.Alias
					}
				} else {
					entry = &LabelEntry{Text: e.Alias}
				}
			}
			if entry.Text == "" && !entry.NoText {
				entry.Text = label
			}
		}
		if entry == nil {
			entry = &LabelEntry{Text: label}
		}
		entries = append(entries, entry)
	}

	var parts []string
	n := len(entries)

	// Generate parts using comma / spacing rules
	for i, entry := range entries {
		var part string
		var suffix string

		part = entry.Text
		if i != n-1 {
			var comma string
			var space string
			if !entry.NoRightComma && !entries[i+1].NoLeftComma {
				comma = ","
			}
			if !entries[i+1].NoLeftSpace {
				space = " "
			}
			suffix = fmt.Sprintf("%s%s", comma, space)
		}

		parts = append(parts, fmt.Sprintf("%s%s", part, suffix))
	}

	return parts
}
