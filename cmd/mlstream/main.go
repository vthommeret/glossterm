package main

import (
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

			var etyms []string
			for _, l := range w.Languages {
				if l.Etymology != "" {
					etyms = append(etyms, fmt.Sprintf("  %s - %s", l.Name, l.Etymology))
				}
			}
			if len(etyms) > 0 {
				fmt.Printf("%s\n", w.Value)
				for _, e := range etyms {
					fmt.Println(e)
				}
			}
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
