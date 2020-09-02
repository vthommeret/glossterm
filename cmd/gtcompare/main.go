package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/vthommeret/glossterm/lib/gt"
)

const defaultInput = "data/words.gob"
const defaultPreviousInput = "data/previous/words.gob"

var input string
var previousInput string

func init() {
	flag.StringVar(&input, "i", defaultInput, "Input file (gob format)")
	flag.StringVar(&previousInput, "pi", defaultPreviousInput, "Previous input file (gob format)")
	flag.Parse()
}

func main() {
	// Get new words
	newWords, err := gt.GetWords(input)
	if err != nil {
		log.Fatalf("Unable to get %q words: %s", input, err)
	}

	// Get previous words
	previousWords, err := gt.GetWords(previousInput)
	if err != nil {
		log.Fatalf("Unable to get %q words: %s", previousInput, err)
	}

	// Remove words
	for w := range previousWords {
		if _, ok := newWords[w]; !ok {
			fmt.Printf("remove %s\n", w)
		}
	}

	ignoreUnexported := cmpopts.IgnoreUnexported(gt.Language{})

	for w, newWord := range newWords {
		if previousWord, ok := previousWords[w]; !ok {
			fmt.Printf("add %s\n", w)
		} else {
			previousWord.Indexed = nil
			//if diff := cmp.Diff(previousWord, newWord, ignoreUnexported); diff != "" {
			//fmt.Printf("update %s\n%s", w, diff)
			if !cmp.Equal(previousWord, newWord, ignoreUnexported) {
				fmt.Printf("update %s\n", w)
			}
		}
	}
}
