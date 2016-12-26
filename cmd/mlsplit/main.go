package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/vthommeret/memory.limited/lib/ml"
)

const total = 5500000 // approximate
const step = total / 100

const defaultOutputFile = "cmd/mlsplit/pages.xml"

var inputFile string
var outputFile string

func init() {
	flag.StringVar(&inputFile, "i", "", "Input file (xml format)")
	flag.StringVar(&outputFile, "o", defaultOutputFile, "Output file (xml format)")
	flag.Parse()
}

func main() {
	if inputFile == "" {
		log.Fatalf("Must specify input file (-i)")
	}

	in, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Unable to open %q input file: %s", inputFile, err)
	}

	nBuckets := runtime.NumCPU()
	buckets := make([][]ml.Page, nBuckets)

	pagesCh := make(chan ml.Page, 10)
	errorsCh := make(chan ml.Error, 10)
	doneCh := make(chan io.ReadCloser)

	go ml.ParseXMLPages(in, pagesCh, errorsCh, doneCh)

	i := 0
	count := 0

Loop:
	for {
		select {
		case e := <-errorsCh:
			log.Fatalf("\nUnable to parse XML: %s", e.Message)
		case f := <-doneCh:
			f.Close()
			break Loop
		case p := <-pagesCh:
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
