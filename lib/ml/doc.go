package ml

import (
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var normalizer transform.Transformer

func init() {
	// See https://blog.golang.org/normalization
	normalizer = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
}

// Document is a minimal document for the search index.
type Document struct {
	Name   string
	Normal string
}

// NewDocument returns a new document from a Word.
func NewDocument(w *Word) *Document {
	return &Document{
		Name:   w.Name,
		Normal: normalize(w.Name),
	}
}

// Type satisfies bleve.Classifier interface.
func Type() string {
	return "word"
}

// Filter out diacritics for search.
func normalize(w string) string {
	s, _, _ := transform.String(normalizer, w)
	return s
}
