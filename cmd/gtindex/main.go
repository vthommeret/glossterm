package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/blevesearch/segment"
	"github.com/vthommeret/glossterm/lib/gt"
	"github.com/vthommeret/glossterm/lib/radix"
)

const defaultInput = "data/words.gob"
const defaultOutput = "data/index.gob"

var input string
var output string

func init() {
	flag.StringVar(&input, "i", defaultInput, "Input file (gob format)")
	flag.StringVar(&output, "o", defaultOutput, "Output file (gob format)")
	flag.Parse()
}

func main() {
	// Get words.
	words, err := gt.GetWords(input)
	if err != nil {
		log.Fatalf("Unable to get %q words: %s", input, err)
	}

	// Populate radix tree.
	count := 0
	termCount := 0
	r := radix.NewTree()
	for _, w := range words {
		ts, err := getTerms(w.Name)
		if err != nil {
			log.Fatalf("Unable to get %q terms: %s", w.Name, err)
		}
		id := radix.EntryID(w.Name)
		for t := range ts {
			r.Insert(t, id)
			termCount++
		}
		count++
	}

	err = gt.WriteGob(output, r)
	if err != nil {
		log.Fatalf("Unable to write and compress %s: %s", output, err)
	}

	fmt.Printf("Wrote %d words (%d terms)\n", count, termCount)
}

// Returns list of unique and normalized terms for a given word.
func getTerms(w string) (terms map[string]bool, err error) {
	terms = make(map[string]bool)
	segmenter := segment.NewWordSegmenterDirect([]byte(w))
	for segmenter.Segment() {
		if segmenter.Type() != segment.None {
			t := strings.ToLower(string(segmenter.Bytes()))
			terms[t] = true
			terms[gt.Normalize(t)] = true
		}
	}
	if err := segmenter.Err(); err != nil {
		return nil, err
	}
	return terms, nil
}
