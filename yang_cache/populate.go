// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package yangcache

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"

	jtafCfg "github.com/chrismarget-j/jtaf/config"
	"github.com/openconfig/goyang/pkg/yang"
)

// Populate populates directories with yang files appropriate for the supplied configuration.
// The returned slice indicates yang file directories relevant to the caller.
func Populate(ctx context.Context, cfg jtafCfg.Cfg, httpClient *http.Client) ([]string, error) {
	var localYangDirs []string
	var err error

	if cfg.CacheIsFresh() {
		log.Println("YANG cache is fresh - skipping update.")

		result, err := yang.PathsWithModules(cfg.YangCacheDir())
		if err != nil {
			return nil, fmt.Errorf("while discovering cached yang dirs in %q - %w", cfg.YangCacheDir(), err)
		}

		return result, nil
	}

	localYangDirs, err = populateYangCacheFromGithub(ctx, cfg, httpClient)
	if err != nil {
		return nil, fmt.Errorf("while populating yang cache from github - %w", err)
	}

	cfg.UpdateCacheFreshnessMarker()

	sort.Strings(localYangDirs)
	return localYangDirs, nil
}
