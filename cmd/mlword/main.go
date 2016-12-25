package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/vthommeret/memory.limited/lib/ml"
)

const defaultInput = "cmd/mlsplit/words.xml"

var inputFile string

func init() {
	flag.StringVar(&inputFile, "i", defaultInput, "Input file (xml format)")
	flag.Parse()
}

func main() {
	args := flag.Args()
	if len(args) == 0 {
		log.Fatalf("Must specify word.")
	}
	w := args[0]

	files, err := ml.GetSplitFiles(inputFile)
	if err != nil {
		log.Fatalf("Unable to get split files: %s", err)
	}
	nBuckets := len(files)

	if nBuckets == 0 {
		log.Fatalf("No split files found for %q.", inputFile)
	}

	pages := make(chan ml.Page)
	errors := make(chan ml.Error)
	done := make(chan io.ReadCloser)

	completed := 0

	for _, f := range files {
		go ml.ParseXMLWord(f, w, pages, errors, done)
	}

	var word *ml.Page

Loop:
	for {
		select {
		case e := <-errors:
			log.Fatalf("\nUnable to parse XML: %s", e.Message)
		case f := <-done:
			f.Close()
			completed++
			if completed == nBuckets {
				break Loop
			}
		case p := <-pages:
			word = &p
			break Loop
		}
	}

	if word == nil {
		fmt.Println("Unable to find word.")
		os.Exit(1)
	}

	enc := gob.NewEncoder(os.Stdout)
	err = enc.Encode(word)
	if err != nil {
		log.Fatalf("Unable to encode word: %s", err)
	}
}
