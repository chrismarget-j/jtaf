package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"

	"github.com/google/go-github/v66/github"
)

const (
	commonYangFileRegex = `\/junos-common-types@\d\d\d\d-\d\d-\d\d.yang`
)

// commonYangFilesRepositoryContent returns map[string]github.RepositoryContent keyed by path within the repository.
// The result describes yang files shared by all platforms described by the jtafConfig
func commonYangFilesRepositoryContent(ctx context.Context, cfg jtafConfig, client *github.Client) (map[string]github.RepositoryContent, error) {
	_, directoryContent, _, err := client.Repositories.GetContents(ctx, cfg.GithubOwnerName, cfg.GithubRepoName, cfg.YangDirCommon, nil)
	if err != nil {
		return nil, fmt.Errorf("while getting repository content %q - %w", path.Join(cfg.GithubOwnerName, cfg.GithubRepoName, cfg.YangDirCommon), err)
	}

	result := make(map[string]github.RepositoryContent, len(directoryContent))
	re := regexp.MustCompile(commonYangFileRegex)
	for i, content := range directoryContent {
		if content.Path == nil {
			return nil, fmt.Errorf("content %d from %q has nil Path element", i, path.Join(cfg.repoPath(), cfg.YangDirCommon))
		}

		if re.MatchString(*content.Path) {
			result[*content.Path] = *content
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no common yang files found in github repo %q path %q using %q", path.Join(cfg.GithubOwnerName, cfg.GithubRepoName), cfg.YangDirCommon, commonYangFileRegex)
	}

	return result, nil
}

// platformYangFilesRepositoryContent returns map[string]github.RepositoryContent keyed by path within the repository.
// The result describes yang files specific to the version and junos described by the jtafConfig
func platformYangFilesRepositoryContent(ctx context.Context, cfg jtafConfig, client *github.Client) (map[string]github.RepositoryContent, error) {
	_, directoryContent, _, err := client.Repositories.GetContents(ctx, cfg.GithubOwnerName, cfg.GithubRepoName, cfg.YangDirPlatform, &github.RepositoryContentGetOptions{Ref: cfg.GitRef})
	if err != nil {
		return nil, fmt.Errorf("while getting repository content %q - %w", path.Join(cfg.repoPath(), cfg.YangDirPlatform), err)
	}

	result := make(map[string]github.RepositoryContent, len(directoryContent))
	for i, content := range directoryContent {
		if content.Path == nil {
			return nil, fmt.Errorf("content %d from %q has nil Path element", i, path.Join(cfg.repoPath(), cfg.YangDirPlatform))
		}

		result[*content.Path] = *content
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no platform yang files found in github repo %q path %q using %q", cfg.repoPath(), cfg.YangDirPlatform, commonYangFileRegex)
	}

	return result, nil
}

// yangFilesRepositoryContent returns map[string]github.RepositoryContent keyed by path within the repository
// The result describes yang files available
func yangFilesRepositoryContent(ctx context.Context, cfg jtafConfig, client *github.Client) (map[string]github.RepositoryContent, error) {
	commonYangFiles, err := commonYangFilesRepositoryContent(ctx, cfg, client)
	if err != nil {
		return nil, fmt.Errorf("while getting common yang file URLs - %w", err)
	}

	platformYangFiles, err := platformYangFilesRepositoryContent(ctx, cfg, client)
	if err != nil {
		return nil, fmt.Errorf("while getting platform yang file URLs - %w", err)
	}

	// Stick common files into the platform file map to unify the maps
	for k, v := range commonYangFiles {
		if _, ok := platformYangFiles[k]; ok {
			// this can never happen because the keys from each map use different directory paths
			return nil, fmt.Errorf("file %q found in both common and platform yang file maps", k)
		}

		platformYangFiles[k] = v
	}

	// At ths point, platformYangFiles contains all yang files (common and platform) keyed by full path
	// within the repository. We intend to cache these files in a single directory within the filesystem.
	// Merging files from multiple paths into a single directory risks name collisions, but I think it's
	// okay...

	// ensure no collisions exist between platform and common file basenames
	baseNameMap := make(map[string]struct{}, len(platformYangFiles))
	for fullPath := range platformYangFiles {
		baseNameMap[path.Base(fullPath)] = struct{}{}
	}
	if len(baseNameMap) != len(platformYangFiles) {
		return nil, fmt.Errorf("cannot cache yang files into %q with current directory strategy "+
			"due to a collision between platform and common file basenames", cfg.yangCacheDir())
	}

	return platformYangFiles, nil
}

// isCached checks whether the file described by the github.RepositoryContent is
// cached at the supplied filename.
func isCached(fileName string, cfg jtafConfig, content github.RepositoryContent) (bool, error) {
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
		err = checkSha256(fileName, yp.RequiredSha256) // check the expected post-patch hash
	} else {
		err = checkGitSha(fileName, fi.Size(), *content.SHA) // check the expected git hash
	}
	if err != nil {
		return false, fmt.Errorf("while checking cached file %q - %w", fileName, err)
	}

	return true, nil
}

// targetFileName returns the filesystem location where the github.RepositoryContent should be stored,
// based on the jtafConfig.
func targetFileName(cfg jtafConfig, content github.RepositoryContent) string {
	return path.Join(cfg.yangCacheDir(), path.Base(*content.Path))
}

// cacheRepositoryContent downloads a github.RepositoryContent into the cache dir
// selected by the jtafConfig. The returned string is the path to the cached file.
func cacheRepositoryContent(ctx context.Context, cfg jtafConfig, client *http.Client, content github.RepositoryContent) (string, error) {
	if content.DownloadURL == nil {
		return "", fmt.Errorf("requested content has nil DownloadURL element")
	}

	if content.Path == nil {
		return "", fmt.Errorf("requested content has nil Path element")
	}

	if content.SHA == nil {
		return "", fmt.Errorf("requested content has nil SHA element")
	}

	targetName := targetFileName(cfg, content)

	ok, err := isCached(targetName, cfg, content)
	if err != nil {
		return "", fmt.Errorf("while checking cache for %q - %w", path.Join(cfg.repoPath(), *content.Path), err)
	}
	if ok {
		return targetName, nil
	}

	err = os.MkdirAll(path.Dir(targetName), 0o755)
	if err != nil {
		return "", fmt.Errorf("while creating cache dir %q - %w", path.Dir(targetName), err)
	}

	req, err := http.NewRequest(http.MethodGet, *content.DownloadURL, nil)
	if err != nil {
		return "", fmt.Errorf("while preparing http request for %q - %w", *content.DownloadURL, err)
	}
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("while fetching %q - %w", *content.DownloadURL, err)
	}
	defer func(c io.Closer) { _ = c.Close() }(resp.Body) // ignoring the error on read seems reasonable

	tmpFile, err := os.CreateTemp(cfg.BaseCacheDir, "."+path.Base(*content.Path))
	if err != nil {
		return "", fmt.Errorf("while making temp download file in %q - %w", cfg.BaseCacheDir, err)
	}
	// do not defer close - we need explicit closeure to apply patches

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		_ = tmpFile.Close()
		return "", fmt.Errorf("while copying data from %q into temp file %q - %w", req.URL, tmpFile.Name(), err)
	}
	err = tmpFile.Close()
	if err != nil {
		return "", fmt.Errorf("while closing tempfile %q - %w", path.Join(cfg.BaseCacheDir, tmpFile.Name()), err)
	}

	if yp, ok := cfg.YangPatches[*content.SHA]; ok {
		err = applyPatch(tmpFile.Name(), yp)
		if err != nil {
			return "", fmt.Errorf("while applying patch to %q - %w", tmpFile.Name(), err)
		}

		err := checkSha256(tmpFile.Name(), yp.RequiredSha256)
		if err != nil {
			return "", fmt.Errorf("failed validating expected checksum (%s) of %q - %w", yp.RequiredSha256, tmpFile.Name(), err)
		}
	}

	err = os.Rename(tmpFile.Name(), targetName)
	if err != nil {
		return "", fmt.Errorf("while renaming tempfile %q to %q - %w", tmpFile.Name(), targetName, err)
	}

	return targetName, nil
}

func populateYangCache(ctx context.Context, cfg jtafConfig, httpClient *http.Client) error {
	githubClient := github.NewClient(httpClient)

	repoPathToRepositoryContent, err := yangFilesRepositoryContent(ctx, cfg, githubClient)
	if err != nil {
		return fmt.Errorf("while getting common yang file URLs - %w", err)
	}

	for repoPath, repositoryContent := range repoPathToRepositoryContent {
		_, err = cacheRepositoryContent(ctx, cfg, httpClient, repositoryContent)
		if err != nil {
			return fmt.Errorf("while caching yang file %q - %w", path.Join(cfg.yangCacheDir(), repoPath), err)
		}
	}

	return nil
}
