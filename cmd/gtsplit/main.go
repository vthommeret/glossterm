package main

import (
	"compress/bzip2"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/vthommeret/glossterm/lib/gt"
)

const total = 5500000 // approximate
const step = total / 100

const defaultInputFile = "data/enwiktionary-latest-pages-articles.xml.bz2"
const defaultOutputFile = "cmd/gtsplit/pages.xml"

var inputFile string
var outputFile string

func init() {
	flag.StringVar(&inputFile, "i", defaultInputFile, "Input file (xml format)")
	flag.StringVar(&outputFile, "o", defaultOutputFile, "Output file (xml format)")
	flag.Parse()
}

type bzipXml struct {
	io.Reader
	io.Closer
}

func main() {
	start := time.Now()

	in, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Unable to open %q input file: %s", inputFile, err)
	}
	bx := bzipXml{bzip2.NewReader(in), in}

	nBuckets := runtime.NumCPU()
	buckets := make([][]gt.Page, nBuckets)

	pagesCh := make(chan gt.Page, 10)
	errorsCh := make(chan gt.Error, 10)
	doneCh := make(chan io.ReadCloser)

	go gt.ParseXMLPages(bx, pagesCh, errorsCh, doneCh)

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

	elapsed := time.Since(start)

	fmt.Printf("\nWrote %d pages to %d files in %s.\n", count, nBuckets, elapsed)
}
