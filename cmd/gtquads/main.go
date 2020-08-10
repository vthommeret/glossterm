package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/vthommeret/glossterm/lib/gt"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/quad"
	"github.com/walle/targz"

	_ "github.com/cayleygraph/cayley/graph/kv/bolt"
)

const inputLang = "fr"
const parentLang = "la"

const defaultInput = "data/words.gob"
const defaultOutput = "data/graph.db"

var input string
var output string

func init() {
	flag.StringVar(&input, "i", defaultInput, "Input file (gob format)")
	flag.StringVar(&output, "o", defaultOutput, "Output folder (boltdb format)")
	flag.Parse()
}

func main() {
	outputCompressed := fmt.Sprintf("%s.tar.gz", output)

	// Get words.
	words, err := gt.GetWords(input)
	if err != nil {
		log.Fatalf("Unable to get %q words: %s", input, err)
	}

	tmpDir, err := ioutil.TempDir("", "words-graph")
	if err != nil {
		log.Fatalf("Unable to create temp directory: %s", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize the database
	err = graph.InitQuadStore("bolt", tmpDir, nil)
	if err != nil {
		log.Fatalf("Unable to init quad store: %s", err)
	}
	graph.IgnoreDuplicates = true

	// Open and use the database
	store, err := cayley.NewGraph("bolt", tmpDir, nil)
	if err != nil {
		log.Fatalf("Unable to open %q output: %s", tmpDir, err)
	}

	// Prepare quads

	count := 0
	quadCount := 0

	var quads []quad.Quad

	for _, w := range words {
		for _, l := range w.Languages {
			if l.Code == inputLang {
				for _, b := range l.Etymology.Borrows {
					if b.FromLang == parentLang {
						quads = append(quads, quad.Make(
							fmt.Sprintf("%s/%s", l.Code, w.Name),
							"borrowing-from",
							fmt.Sprintf("%s/%s", b.FromLang, b.FromWord),
							nil,
						))
						quadCount++
					}
				}
				for _, d := range l.Etymology.Derived {
					if d.FromLang == parentLang {
						quads = append(quads, quad.Make(
							fmt.Sprintf("%s/%s", l.Code, w.Name),
							"derived-from",
							fmt.Sprintf("%s/%s", d.FromLang, d.FromWord),
							nil,
						))
						quadCount++
					}
				}
				for _, i := range l.Etymology.Inherited {
					if i.FromLang == parentLang {
						quads = append(quads, quad.Make(
							fmt.Sprintf("%s/%s", l.Code, w.Name),
							"inherited-from",
							fmt.Sprintf("%s/%s", i.FromLang, i.FromWord),
							nil,
						))
						quadCount++
					}
				}
				for _, m := range l.Etymology.Mentions {
					if m.Lang == parentLang {
						quads = append(quads, quad.Make(
							fmt.Sprintf("%s/%s", l.Code, w.Name),
							"mentions",
							fmt.Sprintf("%s/%s", m.Lang, m.Word),
							nil,
						))
						quadCount++
					}
				}
			} else if l.Code == parentLang {
				for _, s := range l.Etymology.Suffixes {
					if s.Lang != inputLang {
						quads = append(quads, quad.Make(
							fmt.Sprintf("%s/%s", l.Code, w.Name),
							"suffix",
							fmt.Sprintf("%s/%s", s.Lang, s.Root),
							nil,
						))
						quadCount++
					}
				}

				// Map both Links and Descendants to "descendant" for graph search

				for _, ln := range l.Links {
					if ln.Lang != inputLang {
						quads = append(quads, quad.Make(
							fmt.Sprintf("%s/%s", l.Code, w.Name),
							"descendant",
							fmt.Sprintf("%s/%s", ln.Lang, ln.Word),
							nil,
						))
						quadCount++
					}
				}
				for _, d := range l.Descendants {
					if d.Lang != inputLang {
						quads = append(quads, quad.Make(
							fmt.Sprintf("%s/%s", l.Code, w.Name),
							"descendant",
							fmt.Sprintf("%s/%s", d.Lang, d.Word),
							nil,
						))
						quadCount++
					}
				}

			}
		}
		count++
	}

	// Add quads
	log.Printf("Storing quads.")
	err = store.AddQuadSet(quads)
	if err != nil {
		log.Fatalf("Unable to add quads: %s", err)
	}

	// Move temp db to output.
	log.Printf("Moving tmp database to data dir.")
	err = os.Rename(tmpDir, output)
	if err != nil {
		log.Fatalf("Unable to move tmp database to output: %s", err)
	}

	err = targz.Compress(output, outputCompressed)
	if err != nil {
		log.Fatalf("Unable to tar and gzip db: %s", err)
	}

	log.Printf("Read %d words, wrote %d quads.", count, quadCount)
}
