package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/blevesearch/bleve"
	"github.com/vthommeret/memory.limited/lib/ml"
)

const defaultIndexPath = "data/words.bleve"

var indexPath string

func init() {
	flag.StringVar(&indexPath, "i", defaultIndexPath, "Index path (bleve format)")
	flag.Parse()
}

func main() {
	if indexPath == "" {
		log.Fatalf("Must specify index (-i).")
	}

	args := flag.Args()
	if len(args) == 0 {
		log.Fatalf("Must specify query.")
	}
	q := args[0]

	index, err := ml.GetIndex(indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		log.Fatalf("Unable to get %q index: %s", indexPath, err)
	}
	defer index.Close()

	query := bleve.NewPrefixQuery(q)
	search := bleve.NewSearchRequest(query)
	searchResults, err := index.Search(search)
	if err != nil {
		log.Fatalf("Unable to search: %s", err)
	}
	fmt.Println(searchResults)
}
