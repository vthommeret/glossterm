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
			str := tpl.Parameters[i]
			if str != "" {
				v.Field(i).SetString(str)
			}
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
		tf := t.Field(i)
		vf := v.Field(i)
		for _, p := range strings.Split(tf.Tag.Get("names"), ",") {
			if val, ok := paramMap[p]; ok && val != "" {
				vf.SetString(val)
			}
		}
		if tf.Tag.Get("lang") != "" {
			str := vf.String()
			if str != "" {
				vf.SetString(lang.ToParent(str))
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
