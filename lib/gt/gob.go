package gt

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const previousDir = "previous"

// ReadGob reads gob either from path or compressed path.
func ReadGob(path string, data interface{}) error {
	var f io.ReadCloser
	if exists(path) {
		wf, err := os.Open(path)
		if err != nil {
			return err
		}
		f = wf
	} else {
		ext := filepath.Ext(path)
		base := strings.TrimSuffix(path, ext)
		compressed := fmt.Sprintf("%s.gob.gz", base)
		log.Printf("Uncompressing %q.", compressed)
		cf, err := os.Open(compressed)
		if err != nil {
			return err
		}
		gr, err := gzip.NewReader(cf)
		if err != nil {
			return err
		}
		f = gr
	}
	err := gob.NewDecoder(f).Decode(data)
	f.Close()
	if err != nil {
		return err
	}
	return nil
}

func backup(p string) error {
	if exists(p) {
		dir := path.Dir(p)
		base := path.Base(p)
		backup := fmt.Sprintf("%s/%s/%s", dir, previousDir, base)
		err := os.Rename(p, backup)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteGob writes and compresses gob.
func WriteGob(p string, data interface{}, verbose bool, shouldBackup bool) error {
	compressedPath := fmt.Sprintf("%s.gz", p)

	// Backup files for comparison
	if shouldBackup {
		backup(p)
		backup(compressedPath)
	}

	// Gob writer
	out, err := os.Create(p)
	if err != nil {
		return err
	}
	defer out.Close()

	// Gzip writer
	outCompressed, err := os.Create(compressedPath)
	if err != nil {
		return err
	}
	defer outCompressed.Close()
	gw := gzip.NewWriter(outCompressed)
	defer gw.Close()

	// Multi writer
	w := io.MultiWriter(out, gw)

	// Write gob and gzip simultaneously.
	enc := gob.NewEncoder(w)
	err = enc.Encode(data)
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Wrote %q and %q\n", p, compressedPath)
	}

	return nil
}
