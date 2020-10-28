package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"

	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/io/textio"
	"github.com/apache/beam/sdks/go/pkg/beam/x/beamx"
	"github.com/cayleygraph/cayley"
	"github.com/vthommeret/glossterm/lib/gt"
)

const defaultInput = "data/words.gob"
const defaultGraphInput = "data/words.nq"
const defaultOutput = "data/cognates.jsonl"

var input string
var graphInput string
var output string

var graph *cayley.Handle

func init() {
	flag.StringVar(&input, "i", defaultInput, "Input file")
	flag.StringVar(&graphInput, "gi", defaultGraphInput, "Graph file")
	flag.StringVar(&output, "o", defaultOutput, "Output file")
	flag.Parse()

	beam.Init()
}

func cognateFn(word gt.Word) string {
	for lang := range word.Languages {
		if _, ok := gt.SourceLangs[lang]; !ok {
			continue
		}
		word.Languages[lang].Cognates = gt.GetCognates(graph, lang, word.Name)
	}

	b, err := json.Marshal(word)
	if err != nil {
		log.Fatalf("Unable to marshal JSON: %s", err)
	}

	return string(b)
}

func main() {
	// Get words
	wordMap, err := gt.GetWords(input)
	if err != nil {
		log.Fatalf("Unable to get %q words: %s", input, err)
	}

	// Get graph
	g, err := gt.GetGraphNquads(graphInput)
	if err != nil {
		log.Fatalf("Unable to get %s graph: %s", graphInput, err)
	}
	graph = g

	// Create collection

	words := []gt.Word{}
	for _, w := range wordMap {
		words = append(words, *w)
	}

	// Create pipeline

	p := beam.NewPipeline()
	s := p.Root()

	wordList := beam.CreateList(s, words)
	cognates := beam.ParDo(s, cognateFn, wordList)
	textio.Write(s, output, cognates)

	if err := beamx.Run(context.Background(), p); err != nil {
		log.Fatalf("Unable to run pipeline: %s", err)
	}
}
