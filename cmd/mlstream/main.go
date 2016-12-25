package main

import (
	"compress/gzip"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/vthommeret/memory.limited/lib/ml"
)

const defaultInputFile = "cmd/mlsplit/words.xml"
const defaultOutputFile = "data/words.gob"

const total = 200000 // approximate
const step = total / 100

var langMap map[string]bool

var inputFile string
var outputFile string

func init() {
	flag.StringVar(&inputFile, "i", defaultInputFile, "Input file (xml format)")
	flag.StringVar(&outputFile, "o", defaultOutputFile, "Output file (gob format)")
	flag.Parse()
	langMap = ml.ToLangMap(ml.DefaultLangs)
}

func main() {
	outputExt := filepath.Ext(outputFile)
	outputBase := strings.TrimSuffix(outputFile, outputExt)
	outputCompressedFile := fmt.Sprintf("%s.gob.gz", outputBase)

	files, err := ml.GetSplitFiles(inputFile)
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

	pages := make(chan ml.Page, 10)
	errors := make(chan ml.Error, 10)
	done := make(chan io.ReadCloser)

	count := 0
	completed := 0

	for _, f := range files {
		go ml.ParseXML(f, pages, errors, done)
	}

	words := make(map[string]*ml.Word)

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
			w, err := ml.Parse(p, langMap)
			if err != nil {
				var prefix string
				if !errFile {
					prefix = "\n"
				}
				fmt.Fprintf(os.Stderr,
					"%sUnable to parse %q page: %s\n", prefix, p.Title, err)
				continue
			}
			if w.IsEmpty() {
				continue
			}
			words[w.Name] = &w
			count++
			if count == 1 || count%step == 0 {
				fmt.Printf("\r%.1f%% (%d)", 100*float32(count)/total, count)
			}
		}
	}

	fmt.Printf("\n%d total words\n", count)

	// Gob writer
	out, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Unable to create %q file: %s", outputFile, out)
	}
	defer out.Close()

	// Gzip writer
	outCompressed, err := os.Create(outputCompressedFile)
	if err != nil {
		log.Fatalf("Unable to create %q file: %s", outputCompressedFile, out)
	}
	defer outCompressed.Close()
	gw := gzip.NewWriter(outCompressed)
	defer gw.Close()

	// Multi writer
	w := io.MultiWriter(out, gw)

	// Write gob and gzip simultaneously.
	enc := gob.NewEncoder(w)
	err = enc.Encode(words)
	if err != nil {
		log.Fatalf("Unable to encode words: %s", err)
	}

	fmt.Printf("Wrote %q and %q\n", outputFile, outputCompressedFile)
}
