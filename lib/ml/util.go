package ml

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetSplitFiles(pathTemplate string) (files []*os.File, err error) {
	ext := filepath.Ext(pathTemplate)
	base := strings.TrimSuffix(pathTemplate, ext)

	paths, err := filepath.Glob(fmt.Sprintf("%s-*%s", base, ext))
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}
