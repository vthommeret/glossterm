package ml

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

// Filter out diacritics for search.
func Normalize(w string) string {
	s, _, _ := transform.String(normalizer, w)
	return s
}

func GetSplitFiles(pathTemplate string) (files []*os.File, err error) {
	ext := filepath.Ext(pathTemplate)
	base := strings.TrimSuffix(pathTemplate, ext)

	paths, err := filepath.Glob(fmt.Sprintf("%s-*%s", base, ext))
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}
