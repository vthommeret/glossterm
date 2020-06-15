package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"vthommeret/glossterm/lib/gt"
)

const defaultInput = "cmd/gtsplit/pages.xml"

var inputFile string

func init() {
	flag.StringVar(&inputFile, "i", defaultInput, "Input file (xml format)")
	flag.Parse()
}

func main() {
	args := flag.Args()
	if len(args) == 0 {
		log.Fatalf("Must specify page title.")
	}
	t := args[0]

	files, err := gt.GetSplitFiles(inputFile)
	if err != nil {
		log.Fatalf("Unable to get split files: %s", err)
	}
	nBuckets := len(files)

	if nBuckets == 0 {
		log.Fatalf("No split files found for %q.", inputFile)
	}

	pageCh := make(chan gt.Page)
	errorsCh := make(chan gt.Error)
	doneCh := make(chan io.ReadCloser)

	completed := 0

	for _, f := range files {
		go gt.ParseXMLPage(f, t, pageCh, errorsCh, doneCh)
	}

	var page *gt.Page

Loop:
	for {
		select {
		case e := <-errorsCh:
			log.Fatalf("\nUnable to parse XML: %s", e.Message)
		case f := <-doneCh:
			f.Close()
			completed++
			if completed == nBuckets {
				break Loop
			}
		case p := <-pageCh:
			page = &p
			break Loop
		}
	}

	if page == nil {
		fmt.Println("Unable to find word.")
		os.Exit(1)
	}

	e := xml.NewEncoder(os.Stdout)
	e.Indent("", "  ")
	err = e.Encode(page)
	if err != nil {
		log.Fatalf("Unable to XML encode word: %s", err)
	}
}
