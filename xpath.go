package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const xPathSep = "/"

// deviceConfigXpaths reads the specified XML file and returns []string
// representing the deepest paths encountered while parsing the file.
func deviceConfigXpaths(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("while opening device config file %q - %w", filename, err)
	}
	defer func(closer io.Closer) { _ = closer.Close() }(f) // ignoring the error on read seems reasonable

	xmlDec := xml.NewDecoder(f)

	var path []string
	xPathMap := make(map[string]struct{})

	for {
		tok, err := xmlDec.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("while getting xml token - %w", err)
		}

		switch tok := tok.(type) {
		case xml.StartElement:
			// keep track of where we are by appending to the current path
			path = append(path, tok.Name.Local)
			// record this position in the tree
			xPathMap[strings.Join(path, xPathSep)] = struct{}{}
		case xml.EndElement:
			// keep track of where we are by trimming the current path
			path = path[:len(path)-1]
		}
	}

	// trim the map so that only leaf entries remain
	for k := range xPathMap {
		pathElems := strings.Split(k, xPathSep)
		for len(pathElems) > 0 {
			pathElems = pathElems[:len(pathElems)-1]
			delete(xPathMap, strings.Join(pathElems, xPathSep))
		}
	}

	return orderedKeys(xPathMap), nil
}
