package lang

import "testing"

func TestMakeEntryName(t *testing.T) {
	tests := []struct {
		lang string
		name string
		want string
	}{
		{"en", "hello", "hello"},
		{"fr", "étudiant", "étudiant"},
		{"la", "laxō", "laxo"},
		{"es", "¿cuánto?", "cuánto"},
	}
	for _, tt := range tests {
		l, ok := Langs[tt.lang]
		if !ok {
			t.Errorf("Unknown language: %s", l)
		}
		if got := l.MakeEntryName(tt.name); got != tt.want {
			t.Errorf("For %s, l.MakeEntryName(%q) = %q, want %q.", tt.lang, tt.name, got, tt.want)
		}
	}
}
