package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/vthommeret/memory.limited/lib/ml"
)

var langs = []string{"English", "Spanish"}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Must specify file.")
	}

	fp := os.Args[1]
	f, err := os.Open(fp)
	if err != nil {
		log.Fatalf("Unable to open fp: %s", err)
	}

	pages := make(chan ml.Page, 10)
	errors := make(chan ml.Error, 10)
	done := make(chan bool)

	go ml.ParseXML(f, pages, errors, done)

Loop:
	for {
		select {
		case e := <-errors:
			log.Fatalf("Unable to parse XML: %s", e.Message)
		case <-done:
			break Loop
		case p := <-pages:
			w, err := ml.Parse(p)
			if err != nil {
				log.Fatalf("Unable to parse page: %s", err)
			}
			filterLangs(&w, langs)

			b, err := json.MarshalIndent(w, "", "  ")
			if err != nil {
				log.Fatalf("Unable to marshal JSON: %s", err)
			}
			fmt.Printf("%s\n", string(b))
		}
	}
}

func filterLangs(w *ml.Word, langs []string) {
	langMap := make(map[string]bool)
	for _, l := range langs {
		langMap[l] = true
	}
	var filtered []ml.Language
	for _, l := range w.Languages {
		if _, ok := langMap[l.Name]; ok {
			filtered = append(filtered, l)
		}
	}
	w.Languages = filtered
}
