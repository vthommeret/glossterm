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

func findRoots(rootMap map[string]string, word string, allDefns [][]gt.Definition) {
	for _, defns := range allDefns {
		for _, defn := range defns {
			if defn.Root != nil {
				rootMap[word] = defn.Root.Name
			}
		}
	}
}

func createQuads(rootMap map[string]string, typ, lang, word, fromLang, fromWord string) []quad.Quad {
	var quads []quad.Quad
	if root, ok := rootMap[fromWord]; ok {
		quads = append(quads, createQuad(typ, lang, word, fromLang, root))
	}
	quads = append(quads, createQuad(typ, lang, word, fromLang, fromWord))
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

func main() {
	// Get words.
	words, err := gt.GetWords(input)
	if err != nil {
		log.Fatalf("Unable to get %q words: %s", input, err)
	}

	rootMap := map[string]string{}

	for _, w := range words {
		if lang, ok := w.Languages[parentLang]; ok {
			if lang.Definitions != nil {
				findRoots(rootMap, w.Name, lang.AllDefinitions())
			}
		}
	}

	// Prepare quads

	count := 0
	quadCount := 0

	var quads []quad.Quad

	for _, w := range words {
		for _, l := range w.Languages {

			// Latin ancestors
			if _, ok := gt.SourceLangs[l.Code]; ok {
				if l.Etymology != nil {
					for _, b := range l.Etymology.Borrows {
						if b.FromLang == parentLang {
							allQuads := createQuads(rootMap, "borrowing-from", l.Code, w.Name, b.FromLang, b.FromWord)
							quads = append(quads, allQuads...)
							quadCount += len(allQuads)
						}
					}
					for _, d := range l.Etymology.Derived {
						if d.FromLang == parentLang {
							allQuads := createQuads(rootMap, "derived-from", l.Code, w.Name, d.FromLang, d.FromWord)
							quads = append(quads, allQuads...)
							quadCount += len(allQuads)
						}
					}
					for _, i := range l.Etymology.Inherited {
						if i.FromLang == parentLang {
							allQuads := createQuads(rootMap, "inherited-from", l.Code, w.Name, i.FromLang, i.FromWord)
							quads = append(quads, allQuads...)
							quadCount += len(allQuads)
						}
					}
					for _, m := range l.Etymology.Mentions {
						if m.Lang == parentLang {
							allQuads := createQuads(rootMap, "mentions", l.Code, w.Name, m.Lang, m.Word)
							quads = append(quads, allQuads...)
							quadCount += len(allQuads)
						}
					}
				}

				// Latin descendants
			} else if l.Code == parentLang {

				if l.Etymology != nil {
					for _, c := range l.Etymology.Cognates {
						if _, ok := gt.SourceLangs[c.Lang]; ok {
							quads = append(quads, quad.Make(
								fmt.Sprintf("%s/%s", l.Code, w.Name),
								"cognate",
								fmt.Sprintf("%s/%s", c.Lang, c.Word),
								nil,
							))
							if verbose {
								fmt.Printf("%s/%s cognate %s/%s\n", l.Code, w.Name, c.Lang, c.Word)
							}
							quadCount++
						}
					}

					for _, s := range l.Etymology.Suffixes {
						if _, ok := gt.SourceLangs[s.Lang]; ok {
							quads = append(quads, quad.Make(
								fmt.Sprintf("%s/%s", l.Code, w.Name),
								"suffix",
								fmt.Sprintf("%s/%s", s.Lang, s.Root),
								nil,
							))
							if verbose {
								fmt.Printf("%s/%s suffix %s/%s\n", l.Code, w.Name, s.Lang, s.Root)
							}
							quadCount++
						}
					}
				}

				// Map both Links and Descendants to "descendant" for graph search

				for _, ln := range l.Links {
					if _, ok := gt.SourceLangs[ln.Lang]; ok {
						quads = append(quads, quad.Make(
							fmt.Sprintf("%s/%s", l.Code, w.Name),
							"descendant",
							fmt.Sprintf("%s/%s", ln.Lang, ln.Word),
							nil,
						))
						if verbose {
							fmt.Printf("%s/%s descendant (link) %s/%s\n", l.Code, w.Name, ln.Lang, ln.Word)
						}
						quadCount++
					}
				}
				for _, d := range l.Descendants {
					if _, ok := gt.SourceLangs[d.Lang]; ok {
						quads = append(quads, quad.Make(
							fmt.Sprintf("%s/%s", l.Code, w.Name),
							"descendant",
							fmt.Sprintf("%s/%s", d.Lang, d.Word),
							nil,
						))
						if verbose {
							fmt.Printf("%s/%s descendant %s/%s\n", l.Code, w.Name, d.Lang, d.Word)
						}
						quadCount++
					}
				}

			}
		}
		count++
	}

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

	fmt.Printf("Read %d words, wrote %d quads.", count, quadCount)
}
