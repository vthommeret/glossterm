package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"

	"github.com/vthommeret/memory.limited/lib/ml"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Must specify dictionary file and word.")
	}

	fp := os.Args[1]
	f, err := os.Open(fp)
	if err != nil {
		log.Fatalf("Unable to open fp: %s", err)
	}

	word := os.Args[2]

	w, err := ml.ParseXMLWord(f, word)
	if err != nil {
		log.Fatalf("Unable to parse XML: %s", err)
	}

	if w == nil {
		fmt.Println("Unable to find word.")
		os.Exit(1)
	}

	enc := gob.NewEncoder(os.Stdout)
	err = enc.Encode(w)
	if err != nil {
		log.Fatalf("Unable to encode word: %s", err)
	}
}
