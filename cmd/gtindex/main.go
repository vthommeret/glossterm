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

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/blevesearch/segment"

	"golang.org/x/net/context"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

const defaultInput = "data/words.gob"
const defaultPreviousInput = "data/previous/words.gob"
const defaultOutput = "data/words.gob"

const batch = 1000

var input string
var previousInput string
var output string

func init() {
	flag.StringVar(&input, "i", defaultInput, "Input file (gob format)")
	flag.StringVar(&previousInput, "pi", defaultPreviousInput, "Previous input file (gob format)")
	flag.StringVar(&output, "o", defaultOutput, "Output file (gob format)")
	flag.Parse()
}

type IndexAction struct {
	Type IndexActionType
	Word *gt.Word
}

type IndexActionType int

const (
	ActionAdd IndexActionType = iota
	ActionRemove
	ActionUpdate
)

func main() {
	ctx := context.Background()
	opt := option.WithCredentialsFile("./cognate-service-account.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("Unable to initialize Firebase app: %v", err)
	}

	// Get new words.
	newWords, err := gt.GetWords(input)
	if err != nil {
		log.Fatalf("Unable to get %q words: %s", input, err)
	}

	// Get preview words
	previousWords, err := gt.GetWords(previousInput)
	if err != nil {
		log.Fatalf("Unable to get %q words: %s", previousInput, err)
	}

	// Update index

	actions := []IndexAction{}

	// Remove words
	for w, previousWord := range previousWords {
		if previousWord.Indexed == nil {
			continue
		}
		if _, ok := newWords[w]; !ok {
			actions = append(actions, IndexAction{
				Type: ActionRemove,
				Word: previousWord,
			})
		}
	}

	ignoreUnexported := cmpopts.IgnoreUnexported(gt.Language{})

	for w, newWord := range newWords {
		if newWord.Indexed != nil {
			continue
		}

		if !gt.ShouldIndex(newWord) {
			if previousWord, ok := previousWords[w]; ok && previousWord.Indexed != nil {
				actions = append(actions, IndexAction{
					Type: ActionRemove,
					Word: newWord,
				})
			}
			continue
		}

		/*
			b, err := json.MarshalIndent(newWord, "", "  ")
			if err != nil {
				log.Fatalf("Unable to marshal JSON: %s", err)
			}
			fmt.Printf("%s\n", string(b))
		*/

		previousWord, isPrevious := previousWords[w]

		var isUpdated = false
		if previousWord != nil {
			previousWord.Indexed = nil
			isUpdated = !cmp.Equal(previousWord, newWord, ignoreUnexported)
		}

		if !isPrevious || isUpdated {
			var actionType IndexActionType
			if !isPrevious {
				actionType = ActionAdd
			} else if isUpdated {
				actionType = ActionUpdate
			}
			actions = append(actions, IndexAction{
				Type: actionType,
				Word: newWord,
			})
		}
	}

	store, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Unable to initialize Firestore: %v", err)
	}
	defer store.Close()

	var wg sync.WaitGroup

	added := 0
	removed := 0
	updated := 0
	total := 0

	for _, action := range actions {
		wg.Add(1)

		word := action.Word

		ts, err := getTerms(word.Name)
		if err != nil {
			log.Fatalf("Unable to get %q terms: %s", word.Name, err)
		}

		switch action.Type {
		case ActionAdd:
			go updateWord(ctx, store, word, ts, &wg)
			added++
		case ActionUpdate:
			go updateWord(ctx, store, word, ts, &wg)
			updated++
		case ActionRemove:
			go removeWord(ctx, store, word, &wg)
			removed++
		}
		total = added + removed + updated

		if total%batch == 0 {
			commitWords(&wg, newWords, added, updated, removed)
		}
	}

	if total%batch != 0 {
		commitWords(&wg, newWords, added, updated, removed)
	}
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

func removeWord(ctx context.Context, store *firestore.Client, w *gt.Word, wg *sync.WaitGroup) {
	wordsRef := store.Collection("words")

	_, err := wordsRef.Doc(w.Name).Delete(ctx)
	if err != nil {
		log.Fatalf("Failed deleting word: %v", err)
	}

	wg.Done()
}

func updateWord(ctx context.Context, store *firestore.Client, w *gt.Word, ts map[string]bool, wg *sync.WaitGroup) {
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

func commitWords(wg *sync.WaitGroup, words map[string]*gt.Word, added, updated, removed int) {
	wg.Wait()
	err := gt.WriteGob(output, words, false, false)
	if err != nil {
		log.Fatalf("Unable to write and compressed words %s: %s", output, err)
	}
	fmt.Printf("\rAdded %d words; updated %d words; removed %d words", added, updated, removed)
}
