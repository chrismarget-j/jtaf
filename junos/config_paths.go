package junos

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const PathSep = "/"

func GetConfigBreadcrumbTrails(cfgFile string) ([][]string, error) {
	log.Printf("Parsing Junos configuration file %q", cfgFile)
	f, err := os.Open(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("while opening device config file %q - %w", cfgFile, err)
	}
	defer func(closer io.Closer) { _ = closer.Close() }(f) // ignoring the error on read seems reasonable

	xmlDec := xml.NewDecoder(f)

	var breadcrumbs []string
	breadcrumbTrails := make(map[string]struct{})

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
			// keep track of where we are by dropping a breadcrumb
			breadcrumbs = append(breadcrumbs, tok.Name.Local)
			// remember this spot
			breadcrumbTrails[strings.Join(breadcrumbs, PathSep)] = struct{}{}
		case xml.EndElement:
			// keep track of where we are by cleaning up a breadcrumb
			breadcrumbs = breadcrumbs[:len(breadcrumbs)-1]
		}
	}

	// trim the save points so that only leaf nodes remain
	for k := range breadcrumbTrails {
		pathElems := strings.Split(k, PathSep)
		for len(pathElems) > 0 {
			pathElems = pathElems[:len(pathElems)-1]
			delete(breadcrumbTrails, strings.Join(pathElems, PathSep))
		}
	}

	result := make([][]string, len(breadcrumbTrails))
	var i int
	for breadcrumbTrail := range breadcrumbTrails {
		result[i] = strings.Split(breadcrumbTrail, PathSep)
		i++
	}

	log.Printf("Junos configuration contains %d unique configuration paths.", len(result))
	return result, nil
}
