package tpl

import (
	"reflect"
	"strings"

	"github.com/vthommeret/glossterm/lib/lang"
)

type Template struct {
	Action          string
	Parameters      []string
	NamedParameters []Parameter
}

type Parameter struct {
	Name  string
	Value string
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

func toEntryName(langName string, name string) string {
	l, ok := lang.Langs[langName]
	if !ok {
		return name
	}
	return l.MakeEntryName(name)
}
