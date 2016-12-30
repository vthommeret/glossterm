package gt

import (
	"reflect"
	"testing"

	"github.com/vthommeret/glossterm/lib/tpl"
)

func TestParse(t *testing.T) {
	tests := []struct {
		desc string
		word string
		text string
		want Word
	}{
		{
			"Unstructured text",
			"dictionary",
			"definition",
			Word{
				Name: "dictionary",
			},
		},
		{
			"Simple mention",
			"dictionary",
			"==English==\n\n===Etymology===\n{{m|la|dictio||speaking}}",
			Word{
				Name: "dictionary",
				Languages: []Language{
					{
						Code: "en",
						Etymology: Etymology{
							Mentions: []tpl.Mention{
								{Lang: "la", Word: "dictio", Gloss: "speaking"},
							},
						},
					},
				},
			},
		},
		{
			"Named parameter",
			"dictionary",
			"==English==\n\n===Etymology===\n{{m|la|dictio|t=speaking}}",
			Word{
				Name: "dictionary",
				Languages: []Language{
					{
						Code: "en",
						Etymology: Etymology{
							Mentions: []tpl.Mention{
								{Lang: "la", Word: "dictio", Gloss: "speaking"},
							},
						},
					},
				},
			},
		},
		{
			"Nested templates (should be ignored)",
			"dictionary",
			"==English==\n\n===Etymology===\n{{m|la|dictio}}\n{{m|la|dictus{{m|la|dictus}}}}",
			Word{
				Name: "dictionary",
				Languages: []Language{
					{
						Code: "en",
						Etymology: Etymology{
							Mentions: []tpl.Mention{
								{Lang: "la", Word: "dictio"},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		got, err := ParseWord(tt.word, tt.text, DefaultLangMap)
		if err != nil {
			t.Errorf("%s: gt.ParseWord(%q, %q) got error: %s.", tt.desc, tt.word, tt.text, err)
		} else if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s: gt.ParseWord(%q, %q) = %+v, want %+v.", tt.desc, tt.word, tt.text, got, tt.want)
		}
	}
}
