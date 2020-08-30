package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"

	"cloud.google.com/go/firestore"
	"github.com/vthommeret/glossterm/lib/gt"

	"github.com/blevesearch/segment"

	"golang.org/x/net/context"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

const defaultInput = "data/words.gob"
const defaultOutput = "data/index.gob"

const progress = 100

var input string
var output string

func init() {
	flag.StringVar(&input, "i", defaultInput, "Input file (gob format)")
	flag.StringVar(&output, "o", defaultOutput, "Output file (gob format)")
	flag.Parse()
}

func main() {
	ctx := context.Background()
	opt := option.WithCredentialsFile("./cognate-service-account.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("Unable to initialize Firebase app: %v", err)
	}

	store, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Unable to initialize Firestore: %v", err)
	}
	defer store.Close()

	// Get words.
	words, err := gt.GetWords(input)
	if err != nil {
		log.Fatalf("Unable to get %q words: %s", input, err)
	}

	count := 0
	termCount := 0
	//r := radix.NewTree()

	// Create wait group
	var wg sync.WaitGroup

	max := 5

	//completed := 0

	for _, w := range words {
		ts, err := getTerms(w.Name)
		if err != nil {
			log.Fatalf("Unable to get %q terms: %s", w.Name, err)
		}

		fmt.Printf("%s\n", w.Name)

		/*
			wg.Add(1)
			go addWord(ctx, store, w, ts, &wg, &completed)
		*/

		count++
		termCount += len(ts)

		if count == max {
			break
		}
	}

	wg.Wait()

	fmt.Printf("Wrote %d words (%d terms)\n", count, termCount)
}

func addWord(ctx context.Context, store *firestore.Client, w *gt.Word, ts map[string]bool, wg *sync.WaitGroup, completed *int) {
	wordsRef := store.Collection("words")

	_, _, err := wordsRef.Add(ctx, map[string]interface{}{
		"word":      w.Name,
		"terms":     ts,
		"languages": w.Languages,
	})
	if err != nil {
		log.Fatalf("Failed adding word: %v", err)
	}

	*completed++

	if *completed%progress == 0 {
		fmt.Printf("\rAdded %d words", *completed)
	}

	wg.Done()
}

// Returns list of unique and normalized terms for a given word.
func getTerms(w string) (terms map[string]bool, err error) {
	terms = make(map[string]bool)
	segmenter := segment.NewWordSegmenterDirect([]byte(w))
	for segmenter.Segment() {
		if segmenter.Type() != segment.None {
			t := strings.ToLower(string(segmenter.Bytes()))
			terms[t] = true
			terms[gt.Normalize(t)] = true
		}
	}
	if err := segmenter.Err(); err != nil {
		return nil, err
	}
	return terms, nil
}
