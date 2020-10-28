package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/vthommeret/glossterm/lib/gt"
)

const defaultInput = "data/cognates.jsonl"
const defaultOutput = "data/words.gob"

var input string
var output string

func init() {
	flag.StringVar(&input, "i", defaultInput, "Input file (json line format)")
	flag.StringVar(&output, "o", defaultOutput, "Output file (gob format)")
	flag.Parse()
}

func main() {
	f, err := os.Open(input)
	if err != nil {
		log.Fatalf("Unable to open input file: %s, %s", input, err)
	}
	defer f.Close()

	words := map[string]*gt.Word{}
	count := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var word *gt.Word
		if err = json.Unmarshal(scanner.Bytes(), &word); err != nil {
			log.Fatalf("Unable to read JSON: %s", err)
		}
		words[word.Name] = word
		count++
	}

	fmt.Printf("Writing %d words.\n", count)

	err = gt.WriteGob(output, words, true, false)
	if err != nil {
		log.Fatalf("Unable to write and compress %s: %s", output, err)
	}
}
