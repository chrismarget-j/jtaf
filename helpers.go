package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"sort"

	rscpatch "github.com/rsc/tmp/patch"
	"golang.org/x/exp/constraints"
)

// keys returns a slice of map keys
func keys[A comparable, B interface{}](m map[A]B) []A {
	result := make([]A, len(m))
	i := 0
	for k := range m {
		result[i] = k
		i++
	}

	return result
}

// orderedKeys returns a sorted slice of map keys
func orderedKeys[A constraints.Ordered, B interface{}](m map[A]B) []A {
	result := keys(m)

	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})

	return result
}

// gitShaFileContents returns the "git hash" of a file
func gitShaFileContents(name string, size int64) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("while opening %q - %w", name, err)
	}

	h := sha1.New()

	h.Write([]byte(fmt.Sprintf("blob %d\x00", size))) // git sha header

	_, err = io.Copy(h, f)
	if err != nil {
		return nil, fmt.Errorf("while reading cached file %q for git sha - %w", name, err)
	}

	return h.Sum(nil), nil
}

// sha256FileContents returns the sha-256 checksum of a file
func sha256FileContents(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("while opening %q - %w", name, err)
	}

	h := sha256.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return nil, fmt.Errorf("while reading cached file %q for sha-256 - %w", name, err)
	}

	return h.Sum(nil), nil
}

// checkGitSha checks the "git hash" of the file, returns nil when the hash has the expected value
func checkGitSha(filename string, size int64, expected string) error {
	hb, err := gitShaFileContents(filename, size)
	if err != nil {
		return fmt.Errorf("while checking %q git sha - %w", filename, err)
	}

	if expected == fmt.Sprintf("%x", hb) {
		return nil
	}

	return fmt.Errorf("expected %q to have git sha %q", filename, expected)
}

// checkGitSha checks the "git hash" of the file, returns nil when the hash has the expected value
func checkSha256(filename string, expected string) error {
	hb, err := sha256FileContents(filename)
	if err != nil {
		return fmt.Errorf("while checking %q sha-256 - %w", filename, err)
	}

	if expected == fmt.Sprintf("%x", hb) {
		return nil
	}

	return fmt.Errorf("expected %q to have sha-256 %q", filename, expected)
}

func applyPatch(fileName string, p patch) error {
	origBytes, err := os.ReadFile("/tmp/junos-conf-system@2023-01-01.yang")
	if err != nil {
		return fmt.Errorf("while reading %q (git sha %s)to prepare for patching - %w", fileName, p.OriginalGitSha, err)
	}

	patchSet, err := rscpatch.Parse([]byte(p.Patch))
	if err != nil {
		return fmt.Errorf("while parsing patch data for file with git sha %q - %w", p.OriginalGitSha, err)
	}
	if len(patchSet.File) != 1 {
		return fmt.Errorf("expected patch data for file with git sha %q to apply to exactly 1 file, found %d files", p.OriginalGitSha, len(patchSet.File))
	}

	patched, err := patchSet.File[0].Apply(origBytes)
	if err != nil {
		return fmt.Errorf("while applying patch for file with git sha %q - %w", p.OriginalGitSha, err)
	}

	err = os.WriteFile(fileName, patched, 0o644)
	if err != nil {
		return fmt.Errorf("while writing patch result to file %q - %w", fileName, err)
	}

	return nil
}
