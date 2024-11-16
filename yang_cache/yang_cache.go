// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package yangcache

import (
	"fmt"
	"os"

	jtafCfg "github.com/chrismarget-j/jtaf/config"
	"github.com/chrismarget-j/jtaf/helpers"
	"github.com/google/go-github/v66/github"
)

// isCached checks whether the file described by the github.RepositoryContent is
// cached at the supplied filename.
func isCached(fileName string, cfg jtafCfg.Cfg, content github.RepositoryContent) (bool, error) {
	if content.SHA == nil {
		return false, fmt.Errorf("cannot validate cache entry because content shasum is nil")
	}

	fi, err := os.Stat(fileName)
	if err != nil && !os.IsNotExist(err) {
		return false, fmt.Errorf("while stat-ing file %q - %w", fileName, err)
	}

	if os.IsNotExist(err) {
		return false, nil // cache miss - this is fine
	}

	if fi.IsDir() {
		return false, fmt.Errorf("%q is a directory, expected file", fileName)
	}

	if yp, ok := cfg.YangPatches[*content.SHA]; ok {
		err = helpers.CheckSha256(fileName, yp.RequiredSha256) // check the expected post-patch hash
	} else {
		err = helpers.CheckGitSha(fileName, fi.Size(), *content.SHA) // check the expected git hash
	}
	if err != nil {
		return false, fmt.Errorf("while checking cached file %q - %w", fileName, err)
	}

	return true, nil
}
