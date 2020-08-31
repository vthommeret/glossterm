package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/vthommeret/glossterm/lib/gt"

	"github.com/blevesearch/segment"

	"golang.org/x/net/context"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

const defaultInput = "data/words.gob"
const defaultOutput = "data/words.gob"

//const max = 25
const batch = 1000

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
	skipped := 0

	// Create wait group
	var wg sync.WaitGroup

	for _, w := range words {
		// Not supported by Firestore and probably not something people
		// are searching for
		if strings.Contains(w.Name, "/") {
			continue
		}
		if w.Languages == nil {
			continue
		}
		if w.Indexed != nil {
			skipped++
			continue
		}

		// Require definitions
		hasDefinitions := false
		for _, l := range *w.Languages {
			if l.Definitions != nil {
				hasDefinitions = true
				break
			}
		}
		if !hasDefinitions {
			continue
		}

		ts, err := getTerms(w.Name)
		if err != nil {
			log.Fatalf("Unable to get %q terms: %s", w.Name, err)
		}

		/*
			b, err := json.MarshalIndent(w, "", "  ")
			if err != nil {
				log.Fatalf("Unable to marshal JSON: %s", err)
			}
			fmt.Printf("%s\n", string(b))
		*/

		wg.Add(1)
		go addWord(ctx, store, words, w, ts, &wg)

		count++
		termCount += len(ts)

		if count%batch == 0 {
			commitWords(&wg, words, count, termCount, skipped)
		}

		/*
			if count == max {
				break
			}
		*/
	}

	if count%batch != 0 {
		commitWords(&wg, words, count, termCount, skipped)
	}
}

func commitWords(wg *sync.WaitGroup, words map[string]*gt.Word, count, termCount, skipped int) {
	wg.Wait()
	err := gt.WriteGob(output, words, false)
	if err != nil {
		log.Fatalf("Unable to write and compressed words %s: %s", output, err)
	}
	fmt.Printf("\rIndexed %d words (%d terms); %d already indexed", count, termCount, skipped)
}

func addWord(ctx context.Context, store *firestore.Client, words map[string]*gt.Word, w *gt.Word, ts map[string]bool, wg *sync.WaitGroup) {
	wordsRef := store.Collection("words")

	_, err := wordsRef.Doc(w.Name).Set(ctx, map[string]interface{}{
		"name":      w.Name,
		"terms":     ts,
		"languages": w.Languages,
	})
	if err != nil {
		log.Fatalf("Failed indexing word: %v", err)
	}

	now := time.Now()
	w.Indexed = &now

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
