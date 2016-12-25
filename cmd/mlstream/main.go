package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/vthommeret/memory.limited/lib/ml"
)

const total = 200000 // approximate
const step = total / 100

var langs = []string{"en", "es", "fr", "la"}
var langMap map[string]bool

var inputFile string
var outputFile string

func init() {
	flag.StringVar(&inputFile, "i", "", "Input file (xml format)")
	flag.StringVar(&outputFile, "o", "", "Output file (gob format)")
	flag.Parse()
	langMap = ml.ToLangMap(langs)
}

func main() {
	if inputFile == "" {
		log.Fatalf("Must specify input file (-i)")
	}
	if outputFile == "" {
		log.Fatalf("Must specify output file (-o)")
	}

	ext := filepath.Ext(inputFile)
	base := strings.TrimSuffix(inputFile, ext)

	filePaths, err := filepath.Glob(fmt.Sprintf("%s-*%s", base, ext))
	if err != nil {
		log.Fatalf("Invalid glob: %s", err)
	}

	var files []*os.File

	for _, filePath := range filePaths {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatalf("Unable to open %q input file: %s", filePath, err)
		}
		defer file.Close()
		files = append(files, file)
	}

	nBuckets := len(files)

	// Whether stderr is redirected to a file.
	stat, err := os.Stderr.Stat()
	if err != nil {
		log.Fatalf("Unable to stat stderr.")
	}
	errFile := (stat.Mode() & os.ModeCharDevice) == 0

	pages := make(chan ml.Page, 10)
	errors := make(chan ml.Error, 10)
	done := make(chan bool)

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
		case <-done:
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

	out, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Unable to open %q file: %s", outputFile, out)
	}
	defer out.Close()

	enc := gob.NewEncoder(out)
	err = enc.Encode(words)
	if err != nil {
		log.Fatalf("Unable to encode words: %s", err)
	}
}
