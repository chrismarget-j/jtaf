package helpers

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// GitShaFileContents returns the "git hash" of a file
func GitShaFileContents(name string, size int64) ([]byte, error) {
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

// CheckGitSha checks the "git hash" of the file, returns nil when the hash has the expected value
func CheckGitSha(filename string, size int64, expected string) error {
	hb, err := GitShaFileContents(filename, size)
	if err != nil {
		return fmt.Errorf("while checking %q git sha - %w", filename, err)
	}

	if expected == fmt.Sprintf("%x", hb) {
		return nil
	}

	return fmt.Errorf("expected %q to have git sha %q", filename, expected)
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

// CheckSha256 checks the "git hash" of the file, returns nil when the hash has the expected value
func CheckSha256(filename string, expected string) error {
	hb, err := sha256FileContents(filename)
	if err != nil {
		return fmt.Errorf("while checking %q sha-256 - %w", filename, err)
	}

	if expected == fmt.Sprintf("%x", hb) {
		return nil
	}

	return fmt.Errorf("expected %q to have sha-256 %q", filename, expected)
}
