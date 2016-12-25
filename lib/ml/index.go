package ml

import (
	"log"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/simple"

	_ "github.com/blevesearch/bleve/analysis/analyzers/simple_analyzer"
)

const batchSize = 1000

func GetIndex(indexPath string) (bleve.Index, error) {
	return bleve.Open(indexPath)
}

func CreateIndex(indexPath string) (bleve.Index, error) {
	log.Printf("Creating index...")

	indexMapping := bleve.NewIndexMapping()

	wordMapping := bleve.NewDocumentMapping()

	nameFieldMapping := bleve.NewTextFieldMapping()
	nameFieldMapping.Analyzer = simple.Name
	wordMapping.AddFieldMappingsAt("name", nameFieldMapping)

	normalFieldMapping := bleve.NewTextFieldMapping()
	normalFieldMapping.Analyzer = simple.Name
	wordMapping.AddFieldMappingsAt("normal", nameFieldMapping)

	indexMapping.AddDocumentMapping("word", wordMapping)

	index, err := bleve.New(indexPath, indexMapping)
	if err != nil {
		return nil, err
	}

	return index, nil
}

func Index(index bleve.Index, words map[string]*Word) error {
	log.Printf("Indexing words...")
	batch := index.NewBatch()
	batchCount := 0
	count := 0
	startTime := time.Now()
	for _, w := range words {
		d := NewDocument(w)
		batch.Index(d.Name, d)
		batchCount++
		if batchCount >= batchSize {
			err := index.Batch(batch)
			if err != nil {
				return err
			}
			batch = index.NewBatch()
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
