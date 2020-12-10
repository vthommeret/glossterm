package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/vthommeret/glossterm/lib/gt"

	"github.com/cayleygraph/cayley/quad"
	"github.com/cayleygraph/quad/nquads"
)

const parentLang = "la"

const defaultInput = "data/words.gob"
const defaultOutput = "data/words.nq"
const defaultVerbose = false

var input string
var output string
var verbose bool

func init() {
	flag.StringVar(&input, "i", defaultInput, "Input file (gob format)")
	flag.StringVar(&output, "o", defaultOutput, "Output file (nquads format)")
	flag.BoolVar(&verbose, "v", defaultVerbose, "Verbose")
	flag.Parse()
}

func findRoots(rootMap map[string][]string, word string, allDefns [][]gt.Definition) {
	for _, defns := range allDefns {
		for _, defn := range defns {
			if defn.Root != nil {
				rootMap[word] = append(rootMap[word], defn.Root.Name)
			}
		}
	}
}

func createAncestorQuads(rootMap map[string][]string, typ, lang, word, fromLang, fromWord string) []quad.Quad {
	var quads []quad.Quad
	if roots, ok := rootMap[fromWord]; ok {
		for _, root := range roots {
			quads = append(quads, createQuad(typ, lang, word, fromLang, root))
		}
	}
	quads = append(quads, createQuad(typ, lang, word, fromLang, fromWord))
	quads = append(quads, reverseQuads(quads)...)
	return quads
}

func createQuad(typ, lang, word, fromLang, fromWord string) quad.Quad {
	q := quad.Make(
		fmt.Sprintf("%s/%s", lang, word),
		typ,
		fmt.Sprintf("%s/%s", fromLang, fromWord),
		nil,
	)
	if verbose {
		fmt.Printf("%s/%s %s %s/%s\n", lang, word, typ, fromLang, fromWord)
	}
	return q
}

// Most Latin roots don't explicitly list every descendant, so create descendants explicitly.
func reverseQuads(qs []quad.Quad) []quad.Quad {
	var reversed []quad.Quad
	for _, q := range qs {
		reversed = append(reversed, quad.Make(q.Object, "descendant", q.Subject, nil))
	}
	return reversed
}

func uniqueQuads(qs []quad.Quad) []quad.Quad {
	unique := []quad.Quad{}
	seen := map[string]bool{}
	for _, q := range qs {
		k := fmt.Sprintf("%s,%s,%s", q.Subject, q.Predicate, q.Object)
		if _, ok := seen[k]; !ok {
			unique = append(unique, q)
		}
		seen[k] = true
	}
	return unique
}

func main() {
	// Get words.
	words, err := gt.GetWords(input)
	if err != nil {
		log.Fatalf("Unable to get %q words: %s", input, err)
	}

	rootMap := map[string][]string{}

	for _, w := range words {
		if lang, ok := w.Languages[parentLang]; ok {
			if lang.Definitions != nil {
				findRoots(rootMap, w.Name, lang.AllDefinitions())
			}
		}
	}

	// Prepare quads

	count := 0

	var quads []quad.Quad

	for _, w := range words {
		for _, l := range w.Languages {

			// Latin ancestors
			if _, ok := gt.SourceLangs[l.Code]; ok {
				if l.Etymology != nil {
					for _, c := range l.Etymology.Cognates {
						if c.Lang == parentLang {
							allQuads := createAncestorQuads(rootMap, "cognate", l.Code, w.Name, c.Lang, c.Word)
							quads = append(quads, allQuads...)
						}
					}
					for _, s := range l.Etymology.Suffixes {
						if s.Lang == parentLang {
							allQuads := createAncestorQuads(rootMap, "suffix", l.Code, w.Name, s.Lang, s.Root)
							quads = append(quads, allQuads...)
						}
					}
					for _, b := range l.Etymology.Borrows {
						if b.FromLang == parentLang {
							allQuads := createAncestorQuads(rootMap, "borrowing-from", l.Code, w.Name, b.FromLang, b.FromWord)
							quads = append(quads, allQuads...)
						}
					}
					for _, d := range l.Etymology.Derived {
						if d.FromLang == parentLang {
							allQuads := createAncestorQuads(rootMap, "derived-from", l.Code, w.Name, d.FromLang, d.FromWord)
							quads = append(quads, allQuads...)
						}
					}
					for _, i := range l.Etymology.Inherited {
						if i.FromLang == parentLang {
							allQuads := createAncestorQuads(rootMap, "inherited-from", l.Code, w.Name, i.FromLang, i.FromWord)
							quads = append(quads, allQuads...)
						}
					}
					for _, m := range l.Etymology.Mentions {
						if m.Lang == parentLang {
							allQuads := createAncestorQuads(rootMap, "mentions", l.Code, w.Name, m.Lang, m.Word)
							quads = append(quads, allQuads...)
						}
					}
					for _, e := range l.Etymology.Links {
						if e.Lang == parentLang {
							allQuads := createAncestorQuads(rootMap, "etyl", l.Code, w.Name, e.Lang, e.Word)
							quads = append(quads, allQuads...)
						}
					}
				}

				// Latin descendants
			} else if l.Code == parentLang {

				// Map both Links and Descendants to "descendant" for graph search
				for _, ln := range l.Links {
					if _, ok := gt.SourceLangs[ln.Lang]; ok {
						quads = append(quads, createQuad("descendant", l.Code, w.Name, ln.Lang, ln.Word))
					}
				}
				for _, d := range l.Descendants {
					if _, ok := gt.SourceLangs[d.Lang]; ok {
						quads = append(quads, createQuad("descendant", l.Code, w.Name, d.Lang, d.Word))
					}
				}

			}
		}
		count++
	}

	quads = uniqueQuads(quads)

	// Create nquads file

	f, err := os.Create(output)
	if err != nil {
		log.Fatalf("Unable to create nquads file: %s", err)
	}

	// Write nquads file

	w := bufio.NewWriter(f)
	nw := nquads.NewWriter(w)
	_, err = nw.WriteQuads(quads)
	if err != nil {
		log.Fatalf("Error writing quads: %s", err)
	}
	w.Flush()
	f.Close()

	fmt.Printf("Read %d words, wrote %d quads.", count, len(quads))
}
