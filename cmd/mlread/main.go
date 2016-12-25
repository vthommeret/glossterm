package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/vthommeret/memory.limited/lib/ml"
)

func main() {
	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatalf("Unable to stat stdin.")
	}

	var f io.Reader

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		f = os.Stdin
	} else {
		if len(os.Args) < 2 {
			log.Fatalf("Must specify file.")
		}
		fp := os.Args[1]
		f, err = os.Open(fp)
		if err != nil {
			log.Fatalf("Unable to open fp: %s", err)
		}
	}

	dec := gob.NewDecoder(f)

	var words []ml.Word
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