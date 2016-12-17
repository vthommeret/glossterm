package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/vthommeret/memory.limited/lib/ml"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Must specify file.")
	}

	fp := os.Args[1]
	f, err := os.Open(fp)
	if err != nil {
		log.Fatalf("Unable to open fp: %s", err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("Unable to read file: %s", err)
	}

	l := ml.NewLexer(string(b))
	l.PrettyPrint()
}
