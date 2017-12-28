package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	humanize "github.com/dustin/go-humanize"
)

const defaultInput = "https://dumps.wikimedia.org/enwiktionary/latest/enwiktionary-latest-pages-articles.xml.bz2"
const defaultOutput = "data/en.xml.bz2"

var input string
var output string

func init() {
	flag.StringVar(&input, "i", defaultInput, "Input file (Wiktionary dump, .xml.bz2)")
	flag.StringVar(&output, "o", defaultOutput, "Output file (xml)")
	flag.Parse()
}

func main() {
	start := time.Now()

	input = defaultInput
	output = defaultOutput

	// Get content length.
	hr, err := http.Head(input)
	if err != nil {
		log.Fatalf("Unable to get headers for %s: %s", input, err)
	}
	cl := hr.ContentLength

	// Create remote .xml.bz2 input.
	res, err := http.Get(input)
	defer res.Body.Close()

	// Create local .xml output file.
	out, err := os.Create(output)
	defer out.Close()

	// Tee output to counter to display progress.
	tee := io.TeeReader(res.Body, &WriteCounter{Total: uint64(cl)})

	// Create bzip2 decompresser.
	// br := bzip2.NewReader(res.Body)

	fmt.Printf("From: %s\nTo: %s\n\n", input, output)

	// Download, decompress, and write XML dump.
	_, err = io.Copy(out, tee)
	if err != nil {
		log.Fatalf("Error writing dump: %s", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("\n\nDownloaded in %s\n", elapsed)
}

type WriteCounter struct {
	Downloaded uint64
	Total      uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Downloaded += uint64(n)
	fmt.Printf("\r%s / %s (%.2f%%)", humanize.Bytes(wc.Downloaded), humanize.Bytes(wc.Total), 100*float64(wc.Downloaded)/float64(wc.Total))
	return n, nil
}
