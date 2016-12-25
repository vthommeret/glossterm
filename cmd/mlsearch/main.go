package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/blevesearch/bleve"
	_ "github.com/blevesearch/bleve/analysis/analyzers/simple_analyzer"
)

var indexPath string

func init() {
	flag.StringVar(&indexPath, "i", "", "Index path (bleve format)")
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

	index, err := bleve.Open(indexPath)
	if err != nil {
		log.Fatalf("Unable to open %q index: %s", indexPath, err)
	}

	query := bleve.NewPrefixQuery(q)
	search := bleve.NewSearchRequest(query)
	searchResults, err := index.Search(search)
	if err != nil {
		log.Fatalf("Unable to search: %s", err)
	}
	fmt.Println(searchResults)
}