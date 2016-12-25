package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/vthommeret/memory.limited/lib/ml"
)

const total = 200000 // approximate
const step = total / 100

var langs = []string{"en", "es", "fr", "la"}
var langMap map[string]bool

var outputFile string
var inputFile string

func init() {
	flag.StringVar(&inputFile, "i", "", "Input file (xml format)")
	flag.StringVar(&outputFile, "o", "", "Output file (gob format)")
	flag.Parse()
	langMap = ml.ToLangMap(langs)
}

func main() {
	if inputFile == "" {
		log.Fatalf("Must specify input file (-i)")
	}
	if outputFile == "" {
		log.Fatalf("Must specify output file (-o)")
	}

	in, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Unable to open %q input file: %s", inputFile, err)
	}
	defer in.Close()

	pages := make(chan ml.Page, 10)
	errors := make(chan ml.Error, 10)
	done := make(chan bool)

	count := 0

	go ml.ParseXML(in, pages, errors, done)

	var words []ml.Word

Loop:
	for {
		select {
		case e := <-errors:
			log.Fatalf("\nUnable to parse XML: %s", e.Message)
		case <-done:
			break Loop
		case p := <-pages:
			w, err := ml.Parse(p, langMap)
			if err != nil {
				fmt.Printf("\nUnable to parse %q page: %s\n", p.Title, err)
				continue
			}
			if w.IsEmpty() {
				continue
			}
			words = append(words, w)
			count++
			if count == 1 || count%step == 0 {
				fmt.Printf("\r%.1f%% (%d)", 100*float32(count)/total, count)
			}
		}
	}

	fmt.Printf("\n%d total words\n", count)

	out, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Unable to open %q file: %s", outputFile, out)
	}
	defer out.Close()

	enc := gob.NewEncoder(out)
	err = enc.Encode(words)
	if err != nil {
		log.Fatalf("Unable to encode words: %s", err)
	}
}
