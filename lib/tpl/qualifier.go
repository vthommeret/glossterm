package tpl

import (
	"fmt"
	"reflect"
)

// https://en.wiktionary.org/wiki/Template:qualifier
type Qualifier struct {
	Definition string `json:"text,omitempty" firestore:"text,omitempty"`
}

func (tpl *Template) ToQualifier() Qualifier {
	q := Qualifier{}
	tpl.toConcrete(reflect.TypeOf(q), reflect.ValueOf(&q))
	return q
}

func (q *Qualifier) Text() string {
	return fmt.Sprintf("(%s)", q.Definition)
}
