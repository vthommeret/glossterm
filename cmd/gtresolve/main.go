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

var descendantsIndex map[string]map[string]*gt.Language

func main() {
	// Get words.
	words, err := gt.GetWords(input)
	if err != nil {
		log.Fatalf("Unable to get %q words: %s", input, err)
	}

	// Index descendants
	descendantsIndex = map[string]map[string]*gt.Language{}
	for _, w := range words {
		for _, l := range w.Languages {
			if len(l.Descendants) > 0 {
				if _, ok := descendantsIndex[w.Name]; !ok {
					descendantsIndex[w.Name] = map[string]*gt.Language{}
				}
				descendantsIndex[w.Name][l.Code] = l
			}
		}
	}

	// Get (legacy) etymtree descendants
	var etymTreeDescendants map[string]gt.Descendants
	err = gt.ReadGob(descendantsInput, &etymTreeDescendants)
	if err != nil {
		log.Fatalf("Unable to get %q descendants: %s", descendantsInput, err)
	}
	fmt.Printf("%d legacy descendants\n", len(etymTreeDescendants))

	count := 0
	resolved := 0
	resolvedLegacy := 0
	for _, w := range words {
		langs := w.Languages
		for _, l := range langs {
			// Resolve desctrees
			var descendants []tpl.Descendant
			resolveDescTrees(l, &descendants, &resolved)
			l.Descendants = append(l.Descendants, descendants...)

			// Resolve (legacy) etymtrees
			if l.DescendantTrees != nil {
				for _, t := range l.DescendantTrees {
					n := t.ToEntryName()
					if ds, ok := etymTreeDescendants[n]; ok {
						l.Links = append(l.Links, ds.Links...)
						l.Descendants = append(l.Descendants, ds.Descendants...)
						resolvedLegacy++
					}
				}
			}
		}
		count++
	}

	// Dedupe descendants
	for _, w := range words {
		langs := w.Languages
		for _, l := range langs {
			if len(l.Descendants) > 0 {
				l.Descendants = uniqueDescendants(l.Descendants)
			}
		}
	}

	fmt.Printf("Read %d words, resolved %d, resolved (legacy) %d.\n", count, resolved, resolvedLegacy)

	err = gt.WriteGob(output, words, true, false)
	if err != nil {
		log.Fatalf("Unable to write and compress %s: %s", output, err)
	}
}

func resolveDescTrees(l *gt.Language, descendants *[]tpl.Descendant, resolved *int) {
	var levels int
	resolveDescTreesHelper(l, descendants, resolved, &levels)
}

const maxResolveDescTrees = 10

func resolveDescTreesHelper(l *gt.Language, descendants *[]tpl.Descendant, resolved, levels *int) {
	if l.DescTrees != nil {
		for _, d := range l.DescTrees {
			if _, ok := descendantsIndex[d.Word]; ok {
				if descLanguage, ok := descendantsIndex[d.Word][d.Lang]; ok {
					*descendants = append(*descendants, descLanguage.Descendants...)
					(*resolved)++
					(*levels)++
					if *levels <= maxResolveDescTrees {
						resolveDescTreesHelper(descLanguage, descendants, resolved, levels)
					}
				}
			}
		}
	}
}

func uniqueDescendants(ds []tpl.Descendant) []tpl.Descendant {
	var unique []tpl.Descendant
	seen := map[string]bool{}
	for _, d := range ds {
		key := fmt.Sprintf("%s/%s", d.Lang, d.Word)
		if _, ok := seen[key]; ok {
			continue
		}
		unique = append(unique, d)
		seen[key] = true
	}
	return unique
}
