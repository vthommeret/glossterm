package gt

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/vthommeret/glossterm/lib/lang"
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
				Languages: map[string]*Language{
					"en": {
						Code: "en",
						Etymology: &Etymology{
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
				Languages: map[string]*Language{
					"en": {
						Code: "en",
						Etymology: &Etymology{
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
				Languages: map[string]*Language{
					"en": {
						Code: "en",
						Etymology: &Etymology{
							Mentions: []tpl.Mention{
								{Lang: "la", Word: "dictio"},
							},
						},
					},
				},
			},
		},
		{
			"Language list",
			"papyrus",
			"==Latin==\n\n====Descendants====\n* English: {{l|en|papyrus}}, [[paper]]\n* French: {{l|fr|papyrus}}, {{l|fr|papier}}, [[Category:paper]]",
			Word{
				Name: "papyrus",
				Languages: map[string]*Language{
					"la": {
						Code: "la",
						Links: []tpl.Link{
							{Lang: "en", Word: "papyrus"},
							{Lang: "en", Word: "paper"},
							{Lang: "fr", Word: "papyrus"},
							{Lang: "fr", Word: "papier"},
						},
					},
				},
			},
		},
	}

	ignoreUnexported := cmpopts.IgnoreUnexported(Language{})

	for _, tt := range tests {
		got, err := ParseWord(Page{Title: tt.word, Text: tt.text}, lang.DefaultLangMap)
		if err != nil {
			t.Errorf("%s: gt.ParseWord(%q, %q) got error: %s.", tt.desc, tt.word, tt.text, err)
			continue
		}
		if !cmp.Equal(tt.want, got, ignoreUnexported) {
			if diff := cmp.Diff(tt.want, got, ignoreUnexported); diff != "" {
				t.Errorf("%s: gt.ParseWord(%q, %q) diff: %s", tt.desc, tt.word, tt.text, diff)
			}
		}
	}
}
