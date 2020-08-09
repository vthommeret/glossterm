package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"vthommeret/glossterm/lib/gt"
)

const defaultIndexPath = "data/index.gob"
const max = 10

var indexPath string

func init() {
	flag.StringVar(&indexPath, "i", defaultIndexPath, "Index path (gob format)")
	flag.Parse()
}

func main() {
	args := flag.Args()
	if len(args) == 0 {
		log.Fatalf("Must specify query.")
	}
	q := args[0]

	t, err := gt.GetIndex(indexPath)
	if err != nil {
		log.Fatalf("Unable to get radix tree: %s", err)
	}

	rs := t.FindWordsWithPrefix(strings.ToLower(q), max)
	if len(rs) > max {
		rs = rs[:max]
	}

	for i, r := range rs {
		fmt.Printf("%2d. %s\n", i+1, r)
	}

	if len(rs) == 0 {
		fmt.Printf("No results found.\n")
	}
}
