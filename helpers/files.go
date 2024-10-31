package helpers

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func ListFilesWithSuffix(dir, suf string) ([]string, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var result []string
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}

		if !strings.HasSuffix(dirEntry.Name(), suf) {
			continue
		}

		filePath := path.Join(dir, dirEntry.Name())
		abs, err := filepath.Abs(filePath)
		if err != nil {
			return nil, fmt.Errorf("while expanding absolute path to %q - %w", filePath, err)
		}

		result = append(result, abs)
	}

	return result, nil
}
