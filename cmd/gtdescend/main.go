package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/quad"
	"github.com/vthommeret/glossterm/lib/gt"
)

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
	w := fmt.Sprintf("es/%s", os.Args[1])

	graph, err := gt.GetGraph(input)
	if err != nil {
		log.Fatalf("Unable to get %q graph: %s", input, err)
	}

	p := cayley.StartPath(graph, quad.String(w)).
		Out(quad.String("mentions")).
		Out(quad.String("descendant"))

	rs, err := gt.QueryGraph(graph, p)
	if err != nil {
		log.Fatalf("Unable to execute query: %s", err)
	}

	for _, r := range rs {
		fmt.Printf("%s\n", r)
	}
}
