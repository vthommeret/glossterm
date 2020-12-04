package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/vthommeret/glossterm/lib/gt"
)

const defaultInput = "data/words.gob"

var input string

func init() {
	flag.StringVar(&input, "i", defaultInput, "Input file (gob format)")
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

	words, err := gt.GetWords(input)
	if err != nil {
		log.Fatal("Unable to get %q words: %s", input, err)
	}

	if _, ok := words[word]; !ok {
		log.Fatalf("Unable to find word: %s", word)
	}
	if _, ok := words[word].Languages[lang]; !ok {
		log.Fatalf("Unable to find language %s for word: %s", lang, word)
	}

	b, err := json.MarshalIndent(words[word].Languages[lang], "", "  ")
	if err != nil {
		log.Fatalf("Unable to marshal JSON: %s", err)
	}

	fmt.Println(string(b))
}
