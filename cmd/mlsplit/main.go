package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/vthommeret/memory.limited/lib/ml"
)

const total = 5500000 // approximate
const step = total / 100

var inputFile string
var outputFile string

func init() {
	flag.StringVar(&inputFile, "i", "", "Input file (xml format)")
	flag.StringVar(&outputFile, "o", "", "Output file (xml format)")
	flag.Parse()
}

func main() {
	if inputFile == "" {
		log.Fatalf("Must specify input file (-i)")
	}
	if outputFile == "" {
		log.Fatalf("Must specify out file (-o)")
	}

	in, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Unable to open %q input file: %s", inputFile, err)
	}
	defer in.Close()

	nBuckets := runtime.NumCPU()
	buckets := make([][]ml.Page, nBuckets)

	pages := make(chan ml.Page, 10)
	errors := make(chan ml.Error, 10)
	done := make(chan bool)

	go ml.ParseXML(in, pages, errors, done)

	i := 0
	count := 0

Loop:
	for {
		select {
		case e := <-errors:
			log.Fatalf("\nUnable to parse XML: %s", e.Message)
		case <-done:
			break Loop
		case p := <-pages:
			if i > nBuckets-1 {
				i = 0
			}
			buckets[i] = append(buckets[i], p)
			i++
			count++
			if count == 1 || count%step == 0 {
				fmt.Printf("\r%.1f%% (%d)", 100*float32(count)/total, count)
			}
		}
	}

	outputExt := filepath.Ext(outputFile)
	outputBase := strings.TrimSuffix(outputFile, outputExt)

	for i, bucket := range buckets {
		outN := fmt.Sprintf("%s-%d%s", outputBase, i+1, outputExt)
		outNFile, err := os.Create(outN)
		if err != nil {
			log.Fatalf("Unable to open %q file: %s", outN, err)
		}
		defer outNFile.Close()

		e := xml.NewEncoder(outNFile)
		e.Indent("", "  ")
		err = e.Encode(bucket)
		if err != nil {
			log.Fatalf("Unable to encode %q: %s", outN, err)
		}
	}
}
