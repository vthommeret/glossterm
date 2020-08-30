package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/vthommeret/glossterm/lib/gt"
)

const defaultInput = "data/words.gob"
const defaultDescendantsInput = "data/descendants.gob"
const defaultOutput = "data/words.gob"

var input string
var descendantsInput string
var output string

func init() {
	flag.StringVar(&input, "i", defaultInput, "Input file (gob format)")
	flag.StringVar(&descendantsInput, "di", defaultDescendantsInput, "Descendants input file (gob format)")
	flag.StringVar(&output, "o", defaultOutput, "Output file (gob format)")
	flag.Parse()
}

func main() {
	// Get words.
	words, err := gt.GetWords(input)
	if err != nil {
		log.Fatalf("Unable to get %q words: %s", input, err)
	}

	// Get descendants.
	var descendants map[string]gt.Descendants
	err = gt.ReadGob(descendantsInput, &descendants)
	if err != nil {
		log.Fatalf("Unable to get %q descendants: %s", descendantsInput, err)
	}
	fmt.Printf("%d descendants\n", len(descendants))

	// Resolve descendant trees.
	count := 0
	resolved := 0
	for _, w := range words {
		langs := w.Languages
		for _, l := range *langs {
			if l.DescendantTrees != nil {
				for _, t := range l.DescendantTrees {
					n := t.ToEntryName()
					if ds, ok := descendants[n]; ok {
						l.Links =
							append(l.Links, ds.Links...)
						resolved++
					}
				}
			}
		}
		count++
	}

	fmt.Printf("Read %d words, resolved %d.\n", count, resolved)

	err = gt.WriteGob(output, words)
	if err != nil {
		log.Fatalf("Unable to write and compress %s: %s", output, err)
	}
}
