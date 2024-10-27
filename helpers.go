package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"sort"

	"golang.org/x/exp/constraints"
)

// orderedKeys returns a sorted slice of map keys
func orderedKeys[A constraints.Ordered, B interface{}](m map[A]B) []A {
	result := make([]A, len(m))
	i := 0
	for k := range m {
		result[i] = k
		i++
	}

	sort.Slice(result, func(i, j int) bool {
		return i < j
	})
	return result
}

// gitSha returns the "git hash" of a file
func gitSha(name string, size int64) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("while opening %q - %w", name, err)
	}

	h := sha1.New()

	h.Write([]byte(fmt.Sprintf("blob %d\x00", size))) // git sha header

	_, err = io.Copy(h, f)
	if err != nil {
		return nil, fmt.Errorf("while readhing cached file %q - %w", name, err)
	}

	return h.Sum(nil), nil
}
