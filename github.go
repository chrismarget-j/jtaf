package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	jtafCfg "github.com/chrismarget-j/jtaf/config"
	"github.com/google/go-github/v66/github"
)

const envGithubToken = "GITHUB_PUB_API_TOKEN"

// repoContentByDir returns map[string]github.RepositoryContent keyed by path within the repository.
func repoContentByDir(ctx context.Context, dir string, cfg jtafCfg.Cfg, client *github.Client) (map[string]github.RepositoryContent, error) {
	_, directoryContent, _, err := client.Repositories.GetContents(ctx, cfg.YamlRepoInfo.Owner, cfg.YamlRepoInfo.Name, dir, &github.RepositoryContentGetOptions{Ref: cfg.YamlRepoInfo.Ref})
	if err != nil {
		return nil, fmt.Errorf("while getting repository content %q - %w", path.Join(cfg.RepoPath(), dir), err)
	}

	result := make(map[string]github.RepositoryContent, len(directoryContent))
	for i, content := range directoryContent {
		if content.Path == nil {
			return nil, fmt.Errorf("content %d from %q has nil Path element", i, path.Join(cfg.RepoPath(), dir))
		}

		result[*content.Path] = *content
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no platform yang files found in github repo %q path %q", cfg.RepoPath(), dir)
	}

	return result, nil
}

// yangFilesRepoContent returns map[string]github.RepositoryContent keyed by path within the repository
// The result describes yang files available
func yangFilesRepoContent(ctx context.Context, cfg jtafCfg.Cfg, client *github.Client) (map[string]github.RepositoryContent, error) {
	commonYangFiles, err := repoContentByDir(ctx, cfg.RepoDirYangCommon(), cfg, client)
	if err != nil {
		return nil, fmt.Errorf("while getting common yang file URLs - %w", err)
	}

	familyYangFiles, err := repoContentByDir(ctx, cfg.RepoDirYangFamily(), cfg, client)
	if err != nil {
		return nil, fmt.Errorf("while getting platform yang file URLs - %w", err)
	}

	// copy both maps into a single map
	allYangFiles := make(map[string]github.RepositoryContent, len(commonYangFiles)+len(familyYangFiles))
	for k, v := range familyYangFiles {
		allYangFiles[k] = v
	}
	for k, v := range commonYangFiles {
		allYangFiles[k] = v
	}

	return allYangFiles, nil
}

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
func targetFileName(cfg jtafCfg.Cfg, content github.RepositoryContent) string {
	return path.Clean(path.Join(cfg.JunosYangCacheDir(), *content.Path))
}

func validateRepositoryContent(content github.RepositoryContent) error {
	if content.DownloadURL == nil {
		return fmt.Errorf("requested content has nil DownloadURL element")
	}

	if content.HTMLURL == nil {
		return fmt.Errorf("requested content has nil HTMLURL element")
	}

	if content.Path == nil {
		return fmt.Errorf("requested content has nil Path element")
	}

	if content.SHA == nil {
		return fmt.Errorf("requested content has nil SHA element")
	}

	_, err := url.Parse(*content.HTMLURL)
	if err != nil {
		return fmt.Errorf("while validating content HTMLURL element - %w", err)
	}

	return nil
}

// cacheRepositoryContent downloads a github.RepositoryContent into the cache dir
// selected by the jtafConfig. The returned string is the path to the cached file.
func cacheRepositoryContent(ctx context.Context, cfg jtafCfg.Cfg, client *http.Client, content github.RepositoryContent) (string, error) {
	err := validateRepositoryContent(content)
	if err != nil {
		return "", fmt.Errorf("while validating repository content - %w", err)
	}

	targetName := targetFileName(cfg, content)

	ok, err := isCached(targetName, cfg, content)
	if err != nil {
		return "", fmt.Errorf("while checking cache for %q - %w", path.Join(cfg.RepoPath(), *content.Path), err)
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

func githubClient(httpClient *http.Client) *github.Client {
	ghc := github.NewClient(httpClient)

	if token, ok := os.LookupEnv(envGithubToken); ok {
		ghc = ghc.WithAuthToken(token)
	}

	return ghc
}

// populateYangCacheFromGithub populates directories with yang files appropriate for the supplied
// configuration. The returned slice indicates yang file directories relevant to the caller.
func populateYangCacheFromGithub(ctx context.Context, cfg jtafCfg.Cfg, httpClient *http.Client) ([]string, error) {
	repoPathToRepositoryContent, err := yangFilesRepoContent(ctx, cfg, githubClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("while getting common yang file URLs - %w", err)
	}

	yangDirs := make(map[string]struct{})
	for repoPath, repositoryContent := range repoPathToRepositoryContent {
		filePath, err := cacheRepositoryContent(ctx, cfg, httpClient, repositoryContent)
		if err != nil {
			return nil, fmt.Errorf("while caching yang file %q - %w", path.Join(cfg.JunosYangCacheDir(), repoPath), err)
		}

		yangDirs[path.Dir(filePath)] = struct{}{}
	}

	return keys(yangDirs), nil
}
