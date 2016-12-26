package gt

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// GetWords returns words either from path or compressed path.
func GetWords(path string) (map[string]*Word, error) {
	var f io.ReadCloser
	if exists(path) {
		wf, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		f = wf
	} else {
		ext := filepath.Ext(path)
		base := strings.TrimSuffix(path, ext)
		compressed := fmt.Sprintf("%s.gob.gz", base)
		log.Printf("Uncompressing %q.", compressed)
		cf, err := os.Open(compressed)
		if err != nil {
			return nil, err
		}
		gr, err := gzip.NewReader(cf)
		if err != nil {
			return nil, err
		}
		f = gr
	}
	var words map[string]*Word
	err := gob.NewDecoder(f).Decode(&words)
	f.Close()
	if err != nil {
		return nil, err
	}
	return words, nil
}

func exists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}
