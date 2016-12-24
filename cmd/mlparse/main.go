package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/vthommeret/memory.limited/lib/ml"
)

var langs = []string{"en", "es", "fr", "la"}

func main() {
	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatalf("Unable to stat stdin.")
	}

	var f io.Reader

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		f = os.Stdin
	} else {
		if len(os.Args) < 2 {
			log.Fatalf("Must specify file.")
		}
		fp := os.Args[1]
		f, err = os.Open(fp)
		if err != nil {
			log.Fatalf("Unable to open fp: %s", err)
		}
	}

	dec := gob.NewDecoder(f)

	var p ml.Page
	err = dec.Decode(&p)
	if err != nil {
		log.Fatalf("Unable to unmarshal JSON: %s", err)
	}

	w, err := ml.Parse(p)
	if err != nil {
		log.Fatalf("Unable to parse word: %s", err)
	}
	filterLangs(&w, langs)

	b, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		log.Fatalf("Unable to marshal JSON: %s", err)
	}

	fmt.Println(string(b))
}

func filterLangs(w *ml.Word, langs []string) {
	langMap := make(map[string]bool)
	for _, l := range langs {
		langMap[l] = true
	}
	var filtered []ml.Language
	for _, l := range w.Languages {
		if _, ok := langMap[l.Code]; ok {
			filtered = append(filtered, l)
		}
	}
	w.Languages = filtered
}
