package tpl

import (
	"fmt"
	"reflect"
	"strings"
)

// https://en.wiktionary.org/wiki/Template:es-verb_form_of
type SpanishVerb struct {
	Word       string `names:"verb,inf,infinitive" json:"word,omitempty" firestore:"word,omitempty"`
	Ending     string `names:"ending,end" json:"ending" firestore:"ending"`
	Mood       string `names:"mood" json:"mood" firestore:"mood"`
	Tense      string `names:"tense" json:"tense" firestore:"tense"`
	Number     string `names:"number,num" json:"number" firestore:"number"`
	Person     string `names:"person,pers" json:"person" firestore:"person"`
	Formal     string `names:"formal" json:"formal" firestore:"formal"`
	Sense      string `names:"sense" json:"sense" firestore:"sense"`
	Sera       string `names:"sera" json:"sera" firestore:"sera"`
	Gender     string `names:"gender,gen" json:"gender,omitempty" firestore:"gender,omitempty"`
	Participle string `names:"participle,par,part" json:"participle,omitempty" firestore:"participle,omitempty"`
	Voseo      string `names:"voseo" json:"voseo,omitempty" firestore:"voseo,omitempty"`
	Region     string `names:"region" json:"region,omitempty" firestore:"region,omitempty"`
	NoDot      string `names:"nodot" json:"nodot,omitempty" firestore:"nodot,omitempty"`
}

const spanishLang = "es"

// https://en.wiktionary.org/wiki/Template:es-verb_form_of
func (tpl *Template) ToSpanishVerb() SpanishVerb {
	esv := SpanishVerb{}
	tpl.toConcrete(reflect.TypeOf(esv), reflect.ValueOf(&esv))
	esv.Normalize()
	return esv
}

func (esv *SpanishVerb) Normalize() {
	esv.Word = toEntryName(spanishLang, esv.Word)

	// Normalize mood
	switch esv.Mood {
	case "ind", "indicative":
		esv.Mood = "indicative"
	case "subj", "subjunctive":
		esv.Mood = "subjunctive"
	case "imp", "imperative":
		esv.Mood = "imperative"
	case "cond", "conditional=conditional":
	case "par", "part", "participle", "past participle", "past-participle":
		esv.Mood = "participle"
	case "adv", "adverbial", "ger", "gerund", "gerundive", "gerundio", "present participle", "present-participle":
		esv.Mood = "adverbial"
	}

	// Normalize tense
	switch esv.Tense {
	case "pres", "present":
		esv.Tense = "present"
	case "imp", "imperfect":
		esv.Tense = "imperfect"
	case "pret", "preterit", "preterite":
		esv.Tense = "preterite"
	case "fut", "future":
		esv.Tense = "future"
	case "cond", "conditional":
		esv.Tense = "conditional"
	}

	// Normalize number
	switch esv.Number {
	case "s", "sg", "sing", "singular":
		esv.Number = "singular"
	case "p", "pl", "plural":
		esv.Number = "plural"
	}

	// Normalize person
	switch esv.Person {
	case "1", "first", "first person", "first-person":
		esv.Person = "first"
	case "2", "second", "second person", "second-person":
		esv.Person = "second"
	case "3", "third", "third person", "third-person":
		esv.Person = "third"
	case "0", "-", "imp", "impersonal":
		esv.Person = "impersonal"
	}

	// Normalize formal
	switch esv.Formal {
	case "y", "yes":
		esv.Formal = "yes"
	case "n", "no":
		esv.Formal = "no"
	}

	// Normalize gender
	switch esv.Gender {
	case "m", "masc", "masculine":
		esv.Gender = "masculine"
	case "f", "fem", "feminine":
		esv.Gender = "feminine"
	}

	// Normalize sense
	switch esv.Sense {
	case "+", "aff", "affirmative":
		esv.Sense = "affirmative"
	case "-", "neg", "negative":
		esv.Sense = "negative"
	}

	// Normalize ending
	switch esv.Ending {
	case "ar", "-ar":
		esv.Ending = "ar"
	case "er", "-er":
		esv.Ending = "er"
	case "ir", "-ir":
		esv.Ending = "ir"
	}

	/*
				TODO: Support participle
			   -->|participle=<!--
		     -->{{{ms|{{{par|{{{part|{{{participle|}}}}}}}}}}}}<!--
	*/

	// Normalize voseo
	switch esv.Voseo {
	case "yes", "no":
		esv.Voseo = esv.Voseo
	default:
		esv.Voseo = "no"
	}
}

func (esv *SpanishVerb) Text() string {
	var text string

	switch esv.Mood {
	case "adverbial":
		text = esv.adverbialText()
	case "conditional":
		text = esv.conditionalText()
	case "imperative":
		text = esv.imperativeText()
	case "indicative":
		text = esv.indicativeText()
	case "participle":
		text = esv.participleText()
	case "subjunctive":
		text = esv.subjunctiveText()
	}

	var region string
	if esv.Region != "" {
		lbl := Label{Lang: spanishLang, Label1: esv.Region}
		region = lbl.Text()
	}

	var dot string
	if esv.NoDot == "" {
		dot = "."
	}

	return strings.Join(nonEmptyParts(region, text), " ") + dot
}

// https://en.wiktionary.org/wiki/Template:es-verb_form_of/adverbial
func (esv *SpanishVerb) adverbialText() string {
	return fmt.Sprintf("adverbial present participle of %s", esv.Word)
}

// https://en.wiktionary.org/wiki/Template:es-verb_form_of/conditional
func (esv *SpanishVerb) conditionalText() string {
	var incomplete bool

	switch esv.Person {
	case "first", "third":
		switch esv.Number {
		case "singular", "plural":
		default:
			incomplete = true
		}
	case "second":
		switch esv.Number {
		case "singular":
			switch esv.Formal {
			case "no", "yes":
			default:
				incomplete = true
			}
		case "plural":
		default:
			incomplete = true
		}
	case "impersonal":
	default:
		incomplete = true
	}

	var incompleteString string
	if incomplete {
		incompleteString = "a(n)"
	}

	subtenseName := esv.subtenseName()
	subtensePronoun := esv.subtensePronoun()

	return fmt.Sprintf("%s conditional form of %s", strings.Join(nonEmptyParts(incompleteString, subtenseName, subtensePronoun), " "), esv.Word)
}

// https://en.wiktionary.org/wiki/Template:es-verb_form_of/imperative
func (esv *SpanishVerb) imperativeText() string {
	var prefix string

	pattern := fmt.Sprintf("%s-%s-%s", esv.Person, esv.Number, esv.Formal)
	switch pattern {
	case "first-plural-", "first-plural-yes":
		prefix = strings.Join(nonEmptyParts("first-person plural (nosotros, nosotras)", esv.Sense, "imperative"), " ")
	case "second-singular-no":
		switch esv.Sense {
		case "affirmative", "negative":
			var pronoun string
			switch esv.Voseo {
			case "yes":
				pronoun = "(voseo)"
			case "no":
				pronoun = "(tú)"
			}
			prefix = strings.Join(nonEmptyParts(pronoun, "imperative"), " ")
		default:
			prefix = "a"
		}
	case "second-singular-yes":
		prefix = "formal second-person singular (usted) imperative"
	case "second-plural-no":
		switch esv.Sense {
		case "affirmative", "negative":
			prefix = fmt.Sprintf("informal second-person plural (vosotros or vosotras) %s imperative", esv.Sense)
		default:
			prefix = "a"
		}
	case "second-plural-", "second-plural-yes":
		prefix = "second-person plural (ustedes) imperative"
	default:
		prefix = "a"
	}

	return fmt.Sprintf("%s form of %s", prefix, esv.Word)
}

// https://en.wiktionary.org/wiki/Template:es-verb_form_of/indicative
func (esv *SpanishVerb) indicativeText() string {
	var incomplete bool

	switch esv.Tense {
	case "present", "imperfect", "preterite", "future", "conditional":
		switch esv.Person {
		case "first", "third":
			switch esv.Number {
			case "singular", "plural":
			default:
				incomplete = true
			}
		case "second":
			switch esv.Number {
			case "singular":
				switch esv.Formal {
				case "no", "yes":
				default:
					incomplete = true
				}
			case "plural":
			default:
				incomplete = true
			}
		case "impersonal":
		default:
			incomplete = true
		}
	default:
		incomplete = true
	}

	var incompleteString string
	if incomplete {
		incompleteString = "a(n)"
	}

	var prefix string

	subtenseName := esv.subtenseName()
	subtensePronoun := esv.subtensePronoun()

	switch esv.Tense {
	case "present", "imperfect", "preterite", "future":
		prefix = fmt.Sprintf("%s indicative", esv.Tense)
	case "conditional":
		prefix = esv.Tense
	}

	return fmt.Sprintf("%s form of %s", strings.Join(nonEmptyParts(incompleteString, subtenseName, subtensePronoun, prefix), " "), esv.Word)
}

// https://en.wiktionary.org/wiki/Template:es-verb_form_of/participle
func (esv *SpanishVerb) participleText() string {
	return fmt.Sprintf("%s %s past participle of %s", esv.Gender, esv.Number, esv.Word)
}

// https://en.wiktionary.org/wiki/Template:es-verb_form_of/subjunctive
func (esv *SpanishVerb) subjunctiveText() string {
	var incomplete bool

	switch esv.Person {
	case "first", "third":
		switch esv.Number {
		case "singular", "plural":
		default:
			incomplete = true
		}
	case "second":
		switch esv.Number {
		case "singular":
			switch esv.Formal {
			case "no", "yes":
			default:
				incomplete = true
			}
		case "plural":
		default:
			incomplete = true
		}
	case "impersonal":
	default:
		incomplete = true
	}

	switch esv.Tense {
	case "present", "imperfect", "future":
	default:
		incomplete = true
	}

	var incompleteString string
	if incomplete {
		incompleteString = "a(n)"
	}

	var prefix string

	subtenseName := esv.subtenseName()
	subtensePronoun := esv.subtensePronoun()

	switch esv.Tense {
	case "present", "imperfect", "future":
		prefix = fmt.Sprintf("%s subjunctive", esv.Tense)
	}

	return fmt.Sprintf("%s form of %s", strings.Join(nonEmptyParts(incompleteString, subtenseName, subtensePronoun, prefix), " "), esv.Word)
}

// https://en.wiktionary.org/wiki/Template:es-verb_form_of/subtense-name
func (esv *SpanishVerb) subtenseName() string {
	var person string
	var number string

	switch esv.Number {
	case "singular", "plural":
		number = esv.Number
	}

	switch esv.Person {
	case "first", "third":
		person = fmt.Sprintf("%s-person", esv.Person)
	case "second":
		var formal string
		switch esv.Formal {
		case "no":
			formal = "informal"
		case "yes":
			if esv.Number != "plural" {
				formal = "formal"
			}
		}
		person = fmt.Sprintf(strings.Join(nonEmptyParts(formal, "second-person"), " "))
	case "impersonal":
	}

	return strings.Join(nonEmptyParts(person, number), " ")
}

// https://en.wiktionary.org/wiki/Template:es-verb_form_of/subtense-pronoun
func (esv *SpanishVerb) subtensePronoun() string {
	var pronoun string
	switch esv.Voseo {
	case "yes":
		pronoun = "vos"
	case "no":
		pattern := fmt.Sprintf("%s-%s-%s", esv.Number, esv.Person, esv.Formal)
		switch pattern {
		case "singular-first-":
			pronoun = "yo"
		case "singular-second-", "singular-second-no":
			pronoun = "tú"
		case "singular-second-yes":
			pronoun = "usted"
		case "singular-third-":
			pronoun = "él, ella, also used with usted"
		case "plural-first-":
			pronoun = "nosotros, nosotras"
		case "plural-second-no":
			pronoun = "vosotros, vosotras"
		case "plural-second-yes":
			pronoun = "ustedes"
		case "plural-third-":
			pronoun = "ellos, ellas, also used with ustedes"
		default:
			return ""
		}
	}
	return fmt.Sprintf("(%s)", pronoun)
}

func nonEmptyParts(parts ...string) []string {
	var ret []string
	for _, part := range parts {
		if part != "" {
			ret = append(ret, part)
		}
	}
	return ret
}
