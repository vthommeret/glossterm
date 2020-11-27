package tpl

import (
	"fmt"
	"reflect"
	"strings"
)

// https://en.wiktionary.org/wiki/Template:es-compound_of
type SpanishCompound struct {
	VerbStem         string `json:"verbStem,omitempty" firestore:"verbStem,omitempty"`
	InfinitiveEnding string `json:"infinitiveEnding,omitempty" firestore:"infinitiveEnding,omitempty"`
	VerbForm         string `json:"verbForm,omitempty" firestore:"verbForm,omitempty"`
	FirstPronoun     string `json:"firstPronoun,omitempty" firestore:"firstPronoun,omitempty"`
	SecondPronoun    string `json:"secondPronoun,omitempty" firestore:"secondPronoun,omitempty"`
	Mood             string `names:"mood" json:"mood,omitempty" firestore:"mood,omitempty"`
	Person           string `names:"person" json:"person,omitempty" firestore:"person,omitempty"`
	Tense            string `names:"tense" json:"tense,omitempty" firestore:"tense,omitempty"`
	NoDot            string `names:"nodot" json:"nodot,omitempty" firestore:"nodot,omitempty"`
}

func (tpl *Template) ToSpanishCompound() SpanishCompound {
	esc := SpanishCompound{}
	tpl.toConcrete(reflect.TypeOf(esc), reflect.ValueOf(&esc))
	return esc
}

func (esc *SpanishCompound) Word() string {
	switch esc.Mood {
	case "inf", "infinitive":
		return esc.VerbStem + esc.InfinitiveEnding
	case "subjunctive":
		return esc.VerbStem + esc.InfinitiveEnding
	case "pret", "preterite":
		return esc.VerbForm
	case "pres", "present":
		return esc.VerbForm
	case "imperfect", "impf":
		return esc.VerbForm
	case "part", "participle", "adv", "adverbial", "ger", "gerund", "gerundive", "gerundio", "present participle":
		return esc.VerbStem + esc.InfinitiveEnding
	case "imp", "imperative":
		return esc.VerbStem + esc.InfinitiveEnding
	}

	return esc.VerbStem + esc.InfinitiveEnding
}

func (esc *SpanishCompound) Text() string {
	var verbText string

	var word = esc.Word()

	var ending string
	switch esc.InfinitiveEnding {
	case "ar":
		ending = "ar"
	case "er":
		ending = "er"
	case "ir", "ír":
		ending = "ir"
	}

	var nodot = "yes"

	switch esc.Mood {
	case "inf", "infinitive":
		verbText = fmt.Sprintf("the infinitive %s", word)
	case "subjunctive":
		var tense string
		if esc.Tense != "" {
			tense = esc.Tense
		} else {
			tense = "present"
		}
		esv := SpanishVerb{
			Word:   fmt.Sprintf("%s, %s", word, esc.VerbForm),
			Mood:   "subjunctive",
			Tense:  tense,
			Person: "third",
			Number: esc.Person,
			Ending: ending,
			NoDot:  nodot,
		}
		esv.Normalize()
		verbText = fmt.Sprintf("the %s", esv.Text())
	case "pret", "preterite":
		verbText = fmt.Sprintf("the preterite %s", word)
	case "pres", "present":
		verbText = fmt.Sprintf("the present indicative %s", word)
	case "imperfect", "impf":
		verbText = fmt.Sprintf("the imperfect %s", word)
	case "part", "participle", "adv", "adverbial", "ger", "gerund", "gerundive", "gerundio", "present participle":
		esv := SpanishVerb{
			Word:   fmt.Sprintf("%s, %s", word, esc.VerbForm),
			Mood:   "adverbial",
			Ending: ending,
			NoDot:  nodot,
		}
		esv.Normalize()
		verbText = fmt.Sprintf("the %s", esv.Text())
	case "imp", "imperative":
		var person string
		switch esc.Person {
		case "f", "inf", "tú", "tu", "vos", "usted", "ud", "f-pl", "inf-pl", "ustedes", "uds", "vosotros", "v":
			person = "2"
		case "nosotros":
			person = "1"
		}
		var formal string
		switch esc.Person {
		case "usted", "ud", "ustedes", "uds", "f", "f-pl", "nosotros":
			formal = "yes"
		case "tú", "tu", "vos", "inf", "inf-pl", "vosotros", "v":
			formal = "no"
		}
		var number string
		switch esc.Person {
		case "f", "inf", "tú", "tu", "vos", "usted", "ud":
			number = "s"
		case "f-pl", "inf-pl", "ustedes", "uds", "vosotros", "v", "nosotros":
			number = "p"
		}
		var voseo string
		switch esc.Person {
		case "usted", "ud", "ustedes", "uds", "f", "f-pl", "nosotros", "tú", "tu", "inf", "inf-pl", "vosotros", "v":
			voseo = "no"
		case "vos":
			voseo = "yes"
		}
		esv := SpanishVerb{
			Word:   word,
			Sense:  "affirmative",
			Person: person,
			Mood:   "imperative",
			Formal: formal,
			Number: number,
			Ending: ending,
			Voseo:  voseo,
			NoDot:  nodot,
		}
		esv.Normalize()
		verbText = fmt.Sprintf("the %s, %s", esv.Text(), esc.VerbForm)
	default:
		verbText = "a form of the verb"
	}

	var pronounText string
	if esc.SecondPronoun != "" {
		pronounText = fmt.Sprintf("and the pronouns %s and %s", esc.FirstPronoun, esc.SecondPronoun)
	} else {
		var reflexiveText string
		if esc.Mood == "subjunctive" {
			reflexiveText = "reflexive"
		}
		pronounText = strings.Join(nonEmptyParts("and the", reflexiveText, "pronoun", esc.FirstPronoun), " ")
	}

	var dot string
	if esc.NoDot == "" {
		dot = "."
	}

	return fmt.Sprintf("compound of %s %s%s", verbText, pronounText, dot)
}
