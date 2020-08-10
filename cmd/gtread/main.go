package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/vthommeret/glossterm/lib/gt"
)

const defaultInput = "data/words.gob"

var input string

func init() {
	flag.StringVar(&input, "i", defaultInput, "Input file (gob format)")
	flag.Parse()
}

func main() {
	words, err := gt.GetWords(input)
	if err != nil {
		log.Fatal("Unable to get %q words: %s", input, err)
	}

	b, err := json.MarshalIndent(words, "", "  ")
	if err != nil {
		log.Fatalf("Unable to marshal JSON: %s", err)
	}

	fmt.Println(string(b))
}
