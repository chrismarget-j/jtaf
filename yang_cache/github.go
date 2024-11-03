package yangcache

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	jtafcfg "github.com/chrismarget-j/jtaf/config"
	ourgh "github.com/chrismarget-j/jtaf/github"
	"github.com/chrismarget-j/jtaf/helpers"
	"github.com/google/go-github/v66/github"
)

// repoContentByDir returns map[string]github.RepositoryContent keyed by path within the repository.
func repoContentByDir(ctx context.Context, dir string, cfg jtafcfg.Cfg, client *github.Client) (map[string]github.RepositoryContent, error) {
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
func yangFilesRepoContent(ctx context.Context, cfg jtafcfg.Cfg, client *github.Client) (map[string]github.RepositoryContent, error) {
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

// cacheRepositoryContent downloads a github.RepositoryContent into the cache dir
// selected by the jtafConfig. The returned string is the path to the cached file.
func cacheRepositoryContent(ctx context.Context, cfg jtafcfg.Cfg, client *http.Client, content github.RepositoryContent) (string, error) {
	err := ourgh.ValidateRepositoryContent(content)
	if err != nil {
		return "", fmt.Errorf("while validating repository content - %w", err)
	}

	targetName := cfg.TargetFileName(content)

	ok, err := isCached(targetName, cfg, content)
	if err != nil {
		return "", fmt.Errorf("while checking cache for %q - %w", path.Join(cfg.RepoPath(), *content.Path), err)
	}
	if ok {
		log.Printf("...%s - cache OK", *content.DownloadURL)
		return targetName, nil
	}

	log.Printf("Downloading %s", *content.DownloadURL)

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
	// do not defer close - we need explicit closure to apply patches

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		_ = tmpFile.Close()
		return "", fmt.Errorf("while copying data from %q into temp file %q - %w", req.URL, tmpFile.Name(), err)
	}
	err = tmpFile.Close()
	if err != nil {
		return "", fmt.Errorf("while closing tempfile %q - %w", path.Join(cfg.BaseCacheDir, tmpFile.Name()), err)
	}

	err = cfg.Patch(tmpFile.Name(), content)
	if err != nil {
		return "", fmt.Errorf("while patching yang file - %w", err)
	}

	err = os.Rename(tmpFile.Name(), targetName)
	if err != nil {
		return "", fmt.Errorf("while renaming tempfile %q to %q - %w", tmpFile.Name(), targetName, err)
	}

	return targetName, nil
}

// populateYangCacheFromGithub populates directories with yang files appropriate for the supplied
// configuration. The returned slice indicates yang file directories relevant to the caller.
func populateYangCacheFromGithub(ctx context.Context, cfg jtafcfg.Cfg, httpClient *http.Client) ([]string, error) {
	client := ourgh.GithubClient(httpClient)

	urlPath := strings.TrimSuffix(path.Join(cfg.YamlRepoInfo.Owner, cfg.YamlRepoInfo.Name, cfg.YamlRepoInfo.Ref), "/")
	if client.BaseURL != nil {
		urlPath = path.Join(client.BaseURL.String() + urlPath)
	}

	log.Printf("Populating cache from %s...\n", urlPath)

	repoPathToRepositoryContent, err := yangFilesRepoContent(ctx, cfg, client)
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

	return helpers.Keys(yangDirs), nil
}
