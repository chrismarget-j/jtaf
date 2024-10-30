package main

import (
	"context"
	"fmt"
	"github.com/chrismarget-j/jtaf/data/yang"
	"net/http"
	"os"
	"path"
	"sort"
)

// populateYangCacheFromBakedIn drops the baked-in yang files
// (data/yang/publisher/*.yang) into the yang cache dir. The
// returned strings represent yang directories the caller might
// use as a module source.
func populateYangCacheFromBakedIn(cfg jtafConfig) ([]string, error) {
	resultMap := make(map[string]struct{})
	for k, v := range yang.Files {
		fn := path.Join(cfg.BaseCacheDir, k)
		dn := path.Dir(fn)

		err := os.MkdirAll(dn, 0o755)
		if err != nil {
			return nil, fmt.Errorf("while mkdir-ing %q - %w", dn, err)
		}

		tf, err := os.CreateTemp(dn, "."+path.Base(fn))
		if err != nil {
			return nil, fmt.Errorf("while creating temporary file - %w", err)
		}
		//tfn := path.Join(dn, tf.Name())
		tfn := tf.Name()

		_, err = tf.Write([]byte(v))
		if err != nil {
			return nil, fmt.Errorf("while writing to temporary file %q - %w", tfn, err)
		}

		err = tf.Close()
		if err != nil {
			return nil, fmt.Errorf("while closing temporary file %q - %w", tfn, err)
		}

		err = os.Rename(tfn, fn)
		if err != nil {
			return nil, fmt.Errorf("while renaming temporary file %q to %q - %w", tfn, fn, err)
		}

		resultMap[dn] = struct{}{}
	}

	return keys(resultMap), nil
}

// populateYangCache populates directories with yang files appropriate for the supplied configuration.
// The returned slice indicates yang file directories relevant to the caller.
func populateYangCache(ctx context.Context, cfg jtafConfig, httpClient *http.Client) ([]string, error) {
	githubDirs, err := populateYangCacheFromGithub(ctx, cfg, httpClient)
	if err != nil {
		return nil, fmt.Errorf("while populating yang cache from github - %w", err)
	}

	bakedInDirs, err := populateYangCacheFromBakedIn(cfg)
	if err != nil {
		return nil, fmt.Errorf("while populating yang cache with baked-in files - %w", err)
	}

	result := append(githubDirs, bakedInDirs...)
	sort.Strings(result)

	return result, nil
}
