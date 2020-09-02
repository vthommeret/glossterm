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

	"github.com/vthommeret/glossterm/lib/radix"
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

func ShouldIndex(word *Word) bool {
	if word.Indexed != nil {
		return false
	}

	// Not supported by Firestore and probably not something people
	// are searching for
	if strings.Contains(word.Name, "/") {
		return false
	}
	if word.Languages == nil {
		return false
	}

	// Require definitions
	hasDefinitions := false
	for _, l := range *word.Languages {
		if l.Definitions != nil {
			hasDefinitions = true
			break
		}
	}
	if !hasDefinitions {
		return false
	}
	return true
}
