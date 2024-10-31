package jtafCfg

import (
	"fmt"
	"github.com/chrismarget-j/jtaf/helpers"
	rscpatch "github.com/rsc/tmp/patch"
	"os"
)

type Patch struct {
	OriginalGitSha string `yaml:"original_git_sha"`
	RequiredSha256 string `yaml:"required_sha_256"`
	Patch          string `yaml:"diff"`
}

func (p Patch) applyToFile(fn string) error {
	origBytes, err := os.ReadFile(fn)
	if err != nil {
		return fmt.Errorf("while reading %q (git sha %s)to prepare for patching - %w", fn, p.OriginalGitSha, err)
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

	err = os.WriteFile(fn, patched, 0o644)
	if err != nil {
		return fmt.Errorf("while writing patched data to back to %q - %w", fn, err)
	}

	err = helpers.CheckSha256(fn, p.RequiredSha256)
	if err != nil {
		return err
	}

	return nil
}
