package main

import (
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/vthommeret/memory.limited/lib/ml"
)

const defaultInput = "words.gob"

var input string

func init() {
	flag.StringVar(&input, "i", defaultInput, "Input file (gob format)")
	flag.Parse()
}

func main() {
	f, err := os.Open(input)
	if err != nil {
		log.Fatalf("Unable to open fp: %s", err)
	}

	dec := gob.NewDecoder(f)

	var words map[string]*ml.Word
	err = dec.Decode(&words)
	if err != nil {
		log.Fatalf("Unable to decode gob: %s", err)
	}

	b, err := json.MarshalIndent(words, "", "  ")
	if err != nil {
		log.Fatalf("Unable to marshal JSON: %s", err)
	}

	fmt.Println(string(b))
}
