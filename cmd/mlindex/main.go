package main

import (
	"compress/gzip"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/blevesearch/segment"
	"github.com/vthommeret/memory.limited/lib/ml"
	"github.com/vthommeret/memory.limited/lib/radix"
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
	outputExt := filepath.Ext(output)
	outputBase := strings.TrimSuffix(output, outputExt)
	outputCompressed := fmt.Sprintf("%s.gob.gz", outputBase)

	// Get words.
	words, err := ml.GetWords(input)
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

	// Gob writer
	f, err := os.Create(output)
	if err != nil {
		log.Fatalf("Unable to create %q: %s", output, err)
	}
	defer f.Close()

	// Gzip writer
	g, err := os.Create(outputCompressed)
	if err != nil {
		log.Fatalf("Unable to create %q: %s", outputCompressed, err)
	}
	gw := gzip.NewWriter(g)
	defer gw.Close()

	// Multi writer
	w := io.MultiWriter(f, gw)

	// Write gob and gzip simultaneously.
	e := gob.NewEncoder(w)
	err = e.Encode(r)
	if err != nil {
		log.Fatalf("Unable to encode radix tree: %s", err)
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
			terms[ml.Normalize(t)] = true
		}
	}
	if err := segmenter.Err(); err != nil {
		return nil, err
	}
	return terms, nil
}
