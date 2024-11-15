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
	var githubDirs []string
	var err error

	if cfg.CacheIsFresh() {
		log.Println("YANG cache is fresh - skipping update.")
		var result []string
		for _, localDir := range cfg.LocalYangDirs() {
			ypwm, err := yang.PathsWithModules(localDir)
			if err != nil {
				return nil, fmt.Errorf("while discovering cached yang dirs in %q - %w", localDir, err)
			}
			result = append(result, ypwm...)
		}

		return result, nil
	}

	githubDirs, err = populateYangCacheFromGithub(ctx, cfg, httpClient)
	if err != nil {
		return nil, fmt.Errorf("while populating yang cache from github - %w", err)
	}

	cfg.UpdateCacheFreshnessMarker()

	bakedInDirs, err := populateBakedIn(cfg)
	if err != nil {
		return nil, fmt.Errorf("while populating yang cache with baked-in files - %w", err)
	}

	result := append(githubDirs, bakedInDirs...)
	sort.Strings(result)

	return result, nil
}
