package main

import (
	"errors"
	"fmt"

	"github.com/openconfig/goyang/pkg/yang"
	"github.com/openconfig/goyang/pkg/yangentry"
)

const (
	yangSuffix     = ".yang"
	yangConfigRoot = "junos-conf-root"
)

func getYangEntryConfigRoot(dirNames []string) (*yang.Entry, error) {
	var yangIncludePaths []string
	var yangFiles []string

	for _, dirName := range dirNames {
		ypwm, err := yang.PathsWithModules(dirName)
		if err != nil {
			return nil, fmt.Errorf("while parsing %q for .yang files - %w", dirName, err)
		}
		yangIncludePaths = append(yangIncludePaths, ypwm...)

		lfws, err := listFilesWithSuffix(dirName, yangSuffix)
		if err != nil {
			return nil, fmt.Errorf("while listing yang files in %q - %w", dirName, err)
		}
		yangFiles = append(yangFiles, lfws...)
	}

	entries, errs := yangentry.Parse(yangFiles, yangIncludePaths)
	err := errors.Join(errs...)
	if err != nil {
		return nil, fmt.Errorf("encountered %d error(s) while parsing entries from yang files - %w", len(errs), err)
	}

	configRootEntry, ok := entries[yangConfigRoot]
	if !ok {
		return nil, fmt.Errorf("failed to find an entry for %s after parsing %d %s files with %d include directories", yangConfigRoot, len(yangFiles), yangSuffix, len(yangIncludePaths))
	}

	return configRootEntry, nil
}
