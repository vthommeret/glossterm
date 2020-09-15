package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/vthommeret/glossterm/lib/gt"
	"github.com/walle/targz"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/quad"

	_ "github.com/cayleygraph/cayley/graph/kv/bolt"
)

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
			if _, ok := gt.SourceLangs[l.Code]; ok {
				if l.Etymology != nil {
					for _, b := range l.Etymology.Borrows {
						if b.FromLang == parentLang {
							quads = append(quads, quad.Make(
								fmt.Sprintf("%s/%s", l.Code, w.Name),
								"borrowing-from",
								fmt.Sprintf("%s/%s", b.FromLang, b.FromWord),
								nil,
							))
							fmt.Printf("%s/%s borrowing-from %s/%s\n", l.Code, w.Name, b.FromLang, b.FromWord)
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
							fmt.Printf("%s/%s derived-from %s/%s\n", l.Code, w.Name, d.FromLang, d.FromWord)
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
							fmt.Printf("%s/%s inherited-from %s/%s\n", l.Code, w.Name, i.FromLang, i.FromWord)
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
							fmt.Printf("%s/%s mentions %s/%s\n", l.Code, w.Name, m.Lang, m.Word)
							quadCount++
						}
					}
				}
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
							fmt.Printf("%s/%s cognate %s/%s\n", l.Code, w.Name, c.Lang, c.Word)
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
							fmt.Printf("%s/%s suffix %s/%s\n", l.Code, w.Name, s.Lang, s.Root)
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
						fmt.Printf("%s/%s descendant (link) %s/%s\n", l.Code, w.Name, ln.Lang, ln.Word)
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
						fmt.Printf("%s/%s descendant %s/%s\n", l.Code, w.Name, d.Lang, d.Word)
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
