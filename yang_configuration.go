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

func getYangEntryConfigRoot(dirName string) (*yang.Entry, error) {
	yangIncludePaths, err := yang.PathsWithModules(dirName)
	if err != nil {
		return nil, fmt.Errorf("while parsing %q for .yang files - %w", dirName, err)
	}

	yangFiles, err := listFilesWithSuffix(dirName, yangSuffix)
	if err != nil {
		return nil, fmt.Errorf("while listing yang files in %q - %w", dirName, err)
	}

	entries, errs := yangentry.Parse(yangFiles, yangIncludePaths)
	err = errors.Join(errs...)
	if err != nil {
		return nil, fmt.Errorf("encountered %d error(s) while parsing entries from yang files - %w", len(errs), err)
	}

	configRootEntry, ok := entries[yangConfigRoot]
	if !ok {
		return nil, fmt.Errorf("parsed %d %s files with %d include directories, but failed to find an entry for %s", len(yangFiles), yangSuffix, len(yangIncludePaths), yangConfigRoot)
	}

	return configRootEntry, nil
}
