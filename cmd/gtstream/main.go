package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/vthommeret/glossterm/lib/gt"
)

const defaultInputFile = "cmd/gtsplit/pages.xml"
const defaultOutputFile = "data/words.gob"

const total = 200000 // approximate
const step = total / 100

var inputFile string
var outputFile string

func init() {
	flag.StringVar(&inputFile, "i", defaultInputFile, "Input file (xml format)")
	flag.StringVar(&outputFile, "o", defaultOutputFile, "Output file (gob format)")
	flag.Parse()
}

func main() {
	files, err := gt.GetSplitFiles(inputFile)
	if err != nil {
		log.Fatalf("Unable to get split files: %s", err)
	}
	nBuckets := len(files)

	if nBuckets == 0 {
		log.Fatalf("No split files found for %q.", inputFile)
	}

	// Whether stderr is redirected to a file.
	stat, err := os.Stderr.Stat()
	if err != nil {
		log.Fatalf("Unable to stat stderr.")
	}
	errFile := (stat.Mode() & os.ModeCharDevice) == 0

	wordsCh := make(chan gt.Word, 10)
	errorsCh := make(chan gt.Error, 10)
	doneCh := make(chan io.ReadCloser)

	count := 0
	completed := 0

	for _, f := range files {
		go gt.ParseXMLWords(f, wordsCh, errorsCh, doneCh)
	}

	words := make(map[string]*gt.Word)

Loop:
	for {
		select {
		case e := <-errorsCh:
			if e.Fatal {
				log.Fatalf("\nError parsing words: %s", e.Message)
			} else {
				var prefix string
				if !errFile {
					prefix = "\n"
				}
				fmt.Fprintf(os.Stderr, "%sError parsing words: %s\n", prefix, e.Message)
			}
		case f := <-doneCh:
			f.Close()
			completed++
			if completed == nBuckets {
				break Loop
			}
		case w := <-wordsCh:
			words[w.Name] = &w
			count++
			if count == 1 || count%step == 0 {
				fmt.Printf("\r%.1f%% (%d)", 100*float32(count)/total, count)
			}
		}
	}

	fmt.Printf("\n%d total words\n", count)

	err = gt.WriteGob(outputFile, words)
	if err != nil {
		log.Fatalf("Unable to write and compress %s: %s", outputFile, err)
	}
}
