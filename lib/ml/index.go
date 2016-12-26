package ml

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/vthommeret/memory.limited/lib/radix"
)

// GetIndex returns radix tree either from path or compressed path.
func GetIndex(path string) (*radix.Tree, error) {
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
	var t *radix.Tree
	err := gob.NewDecoder(f).Decode(&t)
	f.Close()
	if err != nil {
		return nil, err
	}
	return t, nil
}
