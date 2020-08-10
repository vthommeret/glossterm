package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/vthommeret/glossterm/lib/gt"
)

const sourceLang = "fr"

const defaultInput = "data/graph.db"

var input string

func init() {
	flag.StringVar(&input, "i", defaultInput, "Input file (boltdb format)")
	flag.Parse()
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Must specify word.")
	}
	w := os.Args[1]

	g, err := gt.GetGraph(input)
	if err != nil {
		log.Fatalf("Unable to get %q graph: %s", input, err)
	}

	ds := gt.GetDescendants(g, sourceLang, w)

	for _, d := range ds {
		fmt.Printf("%s (%s)\n", d.Word, d.From)
	}
}
