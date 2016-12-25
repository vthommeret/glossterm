package main

import (
	"encoding/gob"
	"log"
	"os"
	"time"

	"github.com/blevesearch/bleve"
	simple "github.com/blevesearch/bleve/analysis/analyzers/simple_analyzer"
	"github.com/vthommeret/memory.limited/lib/ml"
)

const defaultInput = "words.gob"
const indexPath = "words.bleve"
const batchSize = 1000

func main() {
	var input string
	if len(os.Args) < 2 {
		input = defaultInput
	} else {
		input = os.Args[1]
	}
	f, err := os.Open(input)
	if err != nil {
		log.Fatalf("Unable to open %q input: %s", input, err)
	}

	dec := gob.NewDecoder(f)

	var words map[string]*ml.Word
	err = dec.Decode(&words)
	if err != nil {
		log.Fatalf("Unable to decode gob: %s", err)
	}

	index, err := bleve.Open(indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		indexMapping := bleve.NewIndexMapping()

		wordMapping := bleve.NewDocumentMapping()

		nameFieldMapping := bleve.NewTextFieldMapping()
		nameFieldMapping.Analyzer = simple.Name
		wordMapping.AddFieldMappingsAt("name", nameFieldMapping)

		normalFieldMapping := bleve.NewTextFieldMapping()
		normalFieldMapping.Analyzer = simple.Name
		wordMapping.AddFieldMappingsAt("normal", nameFieldMapping)

		indexMapping.AddDocumentMapping("word", wordMapping)

		index, err = bleve.New(indexPath, indexMapping)
		if err != nil {
			log.Fatalf("Unable to create %q index: %s", indexPath, err)
		}
	} else if err != nil {
		log.Fatalf("Unable to open %q index: %s", indexPath, err)
	}

	err = indexWords(input, index, words)
	if err != nil {
		log.Fatalf("Unable to index words: %s", err)
	}
}

func indexWords(input string, i bleve.Index, ws map[string]*ml.Word) error {
	log.Printf("Indexing %s...", input)
	batch := i.NewBatch()
	batchCount := 0
	count := 0
	startTime := time.Now()
	for _, w := range ws {
		d := ml.NewDocument(w)
		batch.Index(d.Name, d)
		batchCount++
		if batchCount >= batchSize {
			err := i.Batch(batch)
			if err != nil {
				return err
			}
			batch = i.NewBatch()
			batchCount = 0
		}
		count++
		if count%10000 == 0 {
			indexDuration := time.Since(startTime)
			indexDurationSeconds := float64(indexDuration) / float64(time.Second)
			timePerDoc := float64(indexDuration) / float64(count)
			log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
		}
	}
	return nil
}
