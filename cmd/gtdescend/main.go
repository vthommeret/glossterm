package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/quad/nquads"
	"github.com/vthommeret/glossterm/lib/gt"
)

const defaultInput = "data/words.nq"

var input string

func init() {
	flag.StringVar(&input, "i", defaultInput, "Input file (nquads format)")
	flag.Parse()
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Must specify word, e.g. es/helado.")
	}
	w := os.Args[1]

	parts := strings.Split(w, "/")
	if len(parts) < 2 {
		log.Fatalf("Word must be in <lang>/<word> format, e.g. es/helado")
	}

	lang := parts[0]
	word := parts[1]

	fmt.Printf("Word: %s (%s)\n", word, lang)

	f, err := os.Open(input)
	if err != nil {
		log.Fatalf("Unable to open %s input: %s", input, err)
	}

	r := bufio.NewReader(f)
	nr := nquads.NewReader(r, false)

	store, err := cayley.NewMemoryGraph()
	if err != nil {
		log.Fatalf("Unable to create memory store: %s", err)
	}

	for {
		q, err := nr.ReadQuad()
		if err != nil {
			if err != io.EOF {
				log.Fatalf("Unable to read quad: %s", err)
			}
			break
		}
		if !q.IsValid() {
			fmt.Printf("Got invalid quad: %+v\n", q)
		}
		store.AddQuad(q)
	}

	ds := gt.GetCognates(store, lang, word)
	for _, d := range ds {
		fmt.Printf("%s (%s)\n", d.Word, d.From)
	}
}
