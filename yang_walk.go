package main

import (
	"fmt"
	"strings"

	"github.com/openconfig/goyang/pkg/yang"
)

func yangWalk(entry *yang.Entry, paths [][]string) error {
	for p, path := range paths {
		var ne *yang.Entry
		for s, step := range path {
			if ne == nil {
				ne = entry
			}

			var ok bool
			if ne, ok = ne.Dir[step]; !ok {
				return fmt.Errorf("failed traversing yang path %q at step %d (%s)", strings.Join(path, pathSep), s, step)
			}
		}

		fmt.Printf("path %3d (%s) okay\n", p, strings.Join(path, pathSep))
	}

	return nil
}
