package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/cayleygraph/cayley"
	"github.com/montanaflynn/stats"
	"github.com/vthommeret/glossterm/lib/gt"
)

const batch = 1000

const defaultInput = "data/words.gob"
const defaultGraphInput = "data/graph.db"
const defaultOutput = "data/words.gob"

const timeLayout = "1/2 - 3:04pm"

var input string
var graphInput string
var output string

func init() {
	flag.StringVar(&input, "i", defaultInput, "Input file (gob format)")
	flag.StringVar(&graphInput, "gi", defaultGraphInput, "Graph input file (boltdb format)")
	flag.StringVar(&output, "o", defaultOutput, "Output file (gob format)")
	flag.Parse()
}

func main() {
	// Get words
	words, err := gt.GetWords(input)
	if err != nil {
		log.Fatalf("Unable to get %q words: %s", input, err)
	}

	// Get graph
	g, err := gt.GetGraph(graphInput)
	if err != nil {
		log.Fatalf("Unable to get %q graph: %s", graphInput, err)
	}

	count := 0
	updated := 0
	total := len(words)

	var wg sync.WaitGroup

	start := time.Now()
	var durations []float64

	for _, w := range words {
		for langName, lang := range w.Languages {
			if lang.FetchedCognates != nil {
				continue
			}
			if _, ok := gt.SourceLangs[langName]; !ok {
				continue
			}
			wg.Add(1)

			go getCognates(&wg, g, lang, langName, w.Name, &updated)
			count++

			if count%batch == 0 {
				commitWords(&wg, words, &count, &updated, total, &start, &durations)
			}
		}
	}

	if count%batch != 0 {
		commitWords(&wg, words, &count, &updated, total, &start, &durations)
	}
}

func getCognates(wg *sync.WaitGroup, g *cayley.Handle, lang *gt.Language, langName string, word string, updated *int) {
	ds := gt.GetCognates(g, langName, word)
	if len(ds) > 0 {
		lang.Cognates = ds
		*updated++
	}
	now := time.Now()
	lang.FetchedCognates = &now
	wg.Done()
}

func commitWords(wg *sync.WaitGroup, words map[string]*gt.Word, count, updated *int, total int, start *time.Time, durations *[]float64) {
	wg.Wait()
	err := gt.WriteGob(output, words, false, false)
	if err != nil {
		log.Fatalf("Unable to write and compress %s: %s", output, err)
	}

	*durations = append(*durations, float64(time.Since(*start)))
	*start = time.Now()

	median, _ := stats.Median(*durations)
	remaining := time.Duration(1e9 * int64(median*float64(total%batch-*count%batch)/1e9))
	eta := start.Add(remaining)

	fmt.Printf("\rLooked up %d/%d (%.1f%%) words; updated %d. ETA %s (%s)", *count, total, 100*float64(*count)/float64(total), *updated, eta.Format(timeLayout), remaining)
}
