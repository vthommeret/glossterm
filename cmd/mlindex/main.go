package main

import (
	"log"
	"os"

	"github.com/blevesearch/bleve"
	"github.com/vthommeret/memory.limited/lib/ml"
)

const defaultWordsPath = "data/words.gob"
const defaultIndexPath = "data/words.bleve"

func main() {
	var wordsPath string
	if len(os.Args) < 2 {
		wordsPath = defaultWordsPath
	} else {
		wordsPath = os.Args[1]
	}
	indexPath := defaultIndexPath

	words, err := ml.GetWords(wordsPath)
	if err != nil {
		log.Fatalf("Unable to get %q words: %s", wordsPath, err)
	}

	index, err := ml.GetIndex(indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		index, err = ml.CreateIndex(indexPath)
		if err != nil {
			log.Fatalf("Unable to create index: %s", err)
		}
	} else if err != nil {
		log.Fatalf("Unable to index %q words: %s", indexPath, err)
	}

	err = ml.Index(index, words)
	if err != nil {
		log.Fatalf("Unable to index words: %s", err)
	}
}
