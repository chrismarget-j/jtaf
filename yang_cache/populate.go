package yangcache

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	jtafCfg "github.com/chrismarget-j/jtaf/config"
)

// Populate populates directories with yang files appropriate for the supplied configuration.
// The returned slice indicates yang file directories relevant to the caller.
func Populate(ctx context.Context, cfg jtafCfg.Cfg, httpClient *http.Client) ([]string, error) {
	githubDirs, err := populateYangCacheFromGithub(ctx, cfg, httpClient)
	if err != nil {
		return nil, fmt.Errorf("while populating yang cache from github - %w", err)
	}

	bakedInDirs, err := populateBakedIn(cfg)
	if err != nil {
		return nil, fmt.Errorf("while populating yang cache with baked-in files - %w", err)
	}

	result := append(githubDirs, bakedInDirs...)
	sort.Strings(result)

	return result, nil
}
