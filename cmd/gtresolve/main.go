package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/vthommeret/glossterm/lib/gt"
	"github.com/vthommeret/glossterm/lib/tpl"
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

	// Index descendants
	descendants := map[string]map[string][]tpl.Descendant{}
	for _, w := range words {
		for _, l := range w.Languages {
			if len(l.Descendants) > 0 {
				if _, ok := descendants[w.Name]; !ok {
					descendants[w.Name] = map[string][]tpl.Descendant{}
				}
				descendants[w.Name][l.Code] = l.Descendants
			}
		}
	}

	// Get (legacy) etymtree descendants
	var etymTreeDescendants map[string]gt.Descendants
	err = gt.ReadGob(descendantsInput, &etymTreeDescendants)
	if err != nil {
		log.Fatalf("Unable to get %q descendants: %s", descendantsInput, err)
	}
	fmt.Printf("%d descendants\n", len(etymTreeDescendants))

	count := 0
	resolved := 0
	resolvedLegacy := 0
	for _, w := range words {
		langs := w.Languages
		for _, l := range langs {
			// Resolve desctrees
			if l.DescTrees != nil {
				for _, d := range l.DescTrees {
					if _, ok := descendants[d.Word]; ok {
						if _, ok := descendants[d.Word][d.Lang]; ok {
							l.Descendants = append(l.Descendants, descendants[d.Word][d.Lang]...)
							resolved++
						}
					}
				}
			}

			// Resolve (legacy) etymtrees
			if l.DescendantTrees != nil {
				for _, t := range l.DescendantTrees {
					n := t.ToEntryName()
					if ds, ok := etymTreeDescendants[n]; ok {
						l.Links =
							append(l.Links, ds.Links...)
						resolvedLegacy++
					}
				}
			}
		}
		count++
	}

	fmt.Printf("Read %d words, resolved %d, resolved (legacy) %d.\n", count, resolved, resolvedLegacy)

	err = gt.WriteGob(output, words, true, false)
	if err != nil {
		log.Fatalf("Unable to write and compress %s: %s", output, err)
	}
}
