// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package helpers

import (
	"fmt"
	"io"
	"os"
)

var _ io.WriteCloser = (*tmpFileWithRenameOnClose)(nil)

type tmpFileWithRenameOnClose struct {
	newpath string
	tmpFile *os.File
}

func (t tmpFileWithRenameOnClose) Write(p []byte) (n int, err error) {
	return t.tmpFile.Write(p)
}

func (t tmpFileWithRenameOnClose) Close() error {
	err := t.tmpFile.Close()
	if err != nil {
		return fmt.Errorf("while closing temp file %q - %w", t.tmpFile.Name(), err)
	}

	err = os.Rename(t.tmpFile.Name(), t.newpath)
	if err != nil {
		return fmt.Errorf("while renaming temp file %q to %q - %w", t.tmpFile.Name(), t.newpath, err)
	}

	return nil
}

func NewTmpFileWithRenameOnClose(dir, pattern, rename string) (io.WriteCloser, error) {
	f, err := os.CreateTemp(dir, pattern)

	return &tmpFileWithRenameOnClose{
		newpath: rename,
		tmpFile: f,
	}, err
}
