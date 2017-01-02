package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/quad"
	"github.com/vthommeret/glossterm/lib/gt"

	_ "github.com/cayleygraph/cayley/graph/bolt"
)

const defaultInput = "data/words.gob"
const defaultOutput = "data/graph.db"

var input string
var output string

func init() {
	flag.StringVar(&input, "i", defaultInput, "Input file (gob format)")
	flag.StringVar(&output, "o", defaultOutput, "Output file (boltdb format)")
	flag.Parse()
}

func main() {
	outputCompressed := fmt.Sprintf("%s.gz", output)

	// Get words.
	words, err := gt.GetWords(input)
	if err != nil {
		log.Fatalf("Unable to get %q words: %s", input, err)
	}

	tf, err := ioutil.TempFile("", "words-graph")
	if err != nil {
		log.Fatalf("Unable to create temp file: %s", err)
	}
	tfName := tf.Name()
	defer os.Remove(tfName)

	// Initialize the database
	graph.InitQuadStore("bolt", tfName, nil)
	graph.IgnoreDuplicates = true

	// Open and use the database
	store, err := cayley.NewGraph("bolt", tfName, nil)
	if err != nil {
		log.Fatalf("Unable to open %q output: %s", tfName, err)
	}

	// Prepare quads

	count := 0
	quadCount := 0

	var quads []quad.Quad

	for _, w := range words {
		for _, l := range w.Languages {
			if l.Code == "es" {
				for _, m := range l.Etymology.Mentions {
					if m.Lang == "la" {
						quads = append(quads, quad.Make(
							fmt.Sprintf("%s/%s", l.Code, w.Name),
							"mentions",
							fmt.Sprintf("%s/%s", m.Lang, m.Word),
							nil,
						))
						quadCount++
					}
				}
				for _, d := range l.Etymology.Derived {
					if d.FromLang == "la" {
						quads = append(quads, quad.Make(
							fmt.Sprintf("%s/%s", l.Code, w.Name),
							"derived-from",
							fmt.Sprintf("%s/%s", d.FromLang, d.FromWord),
							nil,
						))
						quadCount++
					}
				}
			} else if l.Code == "la" {
				for _, d := range l.Descendants {
					if d.Lang != "es" {
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
	err = os.Rename(tfName, output)
	if err != nil {
		log.Fatalf("Unable to move tmp database to output: %s", err)
	}

	// Open db file
	db, err := os.Open(output)
	if err != nil {
		log.Fatalf("Unable to open db file: %s", err)
	}
	defer db.Close()

	// Gzip writer
	log.Printf("Gzipping database.")
	g, err := os.Create(outputCompressed)
	if err != nil {
		log.Fatalf("Unable to create %q: %s", outputCompressed, err)
	}
	gw := gzip.NewWriter(g)
	defer gw.Close()

	// Gzip db
	_, err = io.Copy(gw, db)
	if err != nil {
		log.Fatalf("Unable to gzip db: %s", err)
	}

	log.Printf("Read %d words, wrote %d quads.", count, quadCount)
}
