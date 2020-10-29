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

const defaultInput = "data/words.gob"
const defaultCognatesInput = "data/cognates.jsonl"
const defaultOutput = "data/words.gob"

var input string
var cognatesInput string
var output string

func init() {
	flag.StringVar(&input, "i", defaultInput, "Input file (gob format)")
	flag.StringVar(&cognatesInput, "ci", defaultCognatesInput, "Cognates input file (json line format)")
	flag.StringVar(&output, "o", defaultOutput, "Output file (gob format)")
	flag.Parse()
}

func main() {
	// Read all cognates

	f, err := os.Open(cognatesInput)
	if err != nil {
		log.Fatalf("Unable to open cognates input file: %s, %s", cognatesInput, err)
	}
	defer f.Close()

	cognateWords := map[string]*gt.Word{}
	count := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var word *gt.Word
		if err = json.Unmarshal(scanner.Bytes(), &word); err != nil {
			log.Fatalf("Unable to read JSON: %s", err)
		}
		cognateWords[word.Name] = word
		count++
	}

	// Resolve cognates

	wordMap, err := gt.GetWords(input)
	if err != nil {
		log.Fatalf("Unable to get %q words: %s", input, err)
	}

	for name, word := range wordMap {
		if cognateWord, ok := cognateWords[name]; ok {
			for code, lang := range cognateWord.Languages {
				if _, ok = word.Languages[code]; ok {
					wordMap[name].Languages[code].Cognates = lang.Cognates
				}
			}
		}
	}

	fmt.Printf("Writing cognates for %d words.\n", count)

	err = gt.WriteGob(output, wordMap, true, false)
	if err != nil {
		log.Fatalf("Unable to write and compress %s: %s", output, err)
	}
}
