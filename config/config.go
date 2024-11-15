package jtafCfg

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v66/github"
	"gopkg.in/yaml.v2"
)

const (
	defaultCacheInterval   = 86400
	defaultGithubOwnerName = "Juniper"
	defaultGithubRepoName  = "yang"
	defaultConfigFile      = "./config.yaml"
	cacheFreshnessMarker   = ".cache_updated"
	configRoot             = "%s-conf-root"
)

var flagC = flag.String("c", defaultConfigFile, "YAML config file")

type githubInfo struct {
	Owner string
	Name  string
	Ref   string
}

type Cfg struct {
	YamlRepoInfo         githubInfo
	JunosConfigFile      string
	JunosVersion         string
	JunosFamily          string
	BaseCacheDir         string
	YangPatches          map[string]Patch
	repoYangDirs         []string
	cacheRefreshInterval time.Duration
}

func (o *Cfg) Patch(filename string, grc github.RepositoryContent) error {
	p, ok := o.YangPatches[*grc.SHA]
	if !ok {
		return nil // no patch; no problem
	}

	return p.applyToFile(filename)
}

// TargetFileName returns the filesystem location where the github.RepositoryContent should be stored.
func (o *Cfg) TargetFileName(grc github.RepositoryContent) string {
	return path.Clean(path.Join(o.YangCacheDir(), *grc.Path))
}

// mkYangCacheDir returns the yang cache dir and returns whether it needed to be created
func (o *Cfg) mkYangCacheDir() (string, bool) {
	yangCacheDir := o.yangCacheBaseDir()

	_, err := os.Stat(yangCacheDir)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(fmt.Errorf("while stat-ing cache dir %q - %w", yangCacheDir, err))
	}

	if os.IsNotExist(err) {
		err = os.MkdirAll(yangCacheDir, 0o755)
		if err != nil {
			log.Fatal(fmt.Errorf("while mkdir-ing cache dir %q - %w", yangCacheDir, err))
		}

		return yangCacheDir, true
	}

	return yangCacheDir, false
}

// RepoPath returns a string like "github.com/Juniper/yang" or "github.com/Juniper/yang@main"
func (o *Cfg) RepoPath() string {
	var gitRef string
	if o.YamlRepoInfo.Ref != "" {
		gitRef = "@" + o.YamlRepoInfo.Ref
	}
	return path.Join("github.com", o.YamlRepoInfo.Owner, o.YamlRepoInfo.Name) + gitRef
}

// yangCacheBaseDir returns the path to the top-level yang cache directory
func (o *Cfg) yangCacheBaseDir() string {
	return path.Clean(path.Join(o.BaseCacheDir, "yang"))
}

func (o *Cfg) YangCacheDir() string {
	ref := o.YamlRepoInfo.Ref
	if ref != "" && !strings.HasPrefix(ref, "@") {
		ref = "@" + ref
	}

	return path.Clean(path.Join(o.yangCacheBaseDir(), o.JunosFamily, o.JunosVersion, ref))
}

func Get() (Cfg, error) {
	flag.Parse()

	cfgBytes, err := os.ReadFile(*flagC)
	if err != nil {
		return Cfg{}, fmt.Errorf("while reading config file %q - %w", *flagC, err)
	}

	var yamlConfig struct {
		DeviceConfigFile     string                         `yaml:"junos_config_xml"`
		Family               string                         `yaml:"junos_family"`
		Version              string                         `yaml:"junos_version"`
		GitRef               string                         `yaml:"git_ref"`
		GithubOwnerName      string                         `yaml:"github_owner_name"`
		GithubRepoName       string                         `yaml:"github_repo_name"`
		CacheDir             string                         `yaml:"cache_dir"`
		YangPatches          []Patch                        `yaml:"yang_patches"`
		GitRepoYangPaths     map[string]map[string][]string `yaml:"git_repo_yang_paths"`
		CacheRefreshInterval *int                           `yaml:"cache_refresh_interval"`
	}

	err = yaml.Unmarshal(cfgBytes, &yamlConfig)
	if err != nil {
		return Cfg{}, fmt.Errorf("while parsing config file %q - %w", *flagC, err)
	}

	if yamlConfig.CacheRefreshInterval == nil {
		cri := defaultCacheInterval
		yamlConfig.CacheRefreshInterval = &cri
	}

	if yamlConfig.GithubOwnerName == "" {
		yamlConfig.GithubOwnerName = defaultGithubOwnerName
	}

	if yamlConfig.GithubRepoName == "" {
		yamlConfig.GithubRepoName = defaultGithubRepoName
	}

	deviceConfigFile, err := filepath.Abs(yamlConfig.DeviceConfigFile)
	if err != nil {
		return Cfg{}, fmt.Errorf("while expanding file path %q - %w", yamlConfig.DeviceConfigFile, err)
	}

	repoYangDirs, ok := yamlConfig.GitRepoYangPaths[yamlConfig.Version][yamlConfig.Family]
	if !ok {
		return Cfg{}, fmt.Errorf("unknown github location for %q version %q YAML files - check the configuration", yamlConfig.Family, yamlConfig.Version)
	}

	yangPatches := make(map[string]Patch, len(yamlConfig.YangPatches))
	for _, p := range yamlConfig.YangPatches {
		yangPatches[p.OriginalGitSha] = p
	}
	if len(yangPatches) != len(yamlConfig.YangPatches) {
		return Cfg{}, fmt.Errorf("unexpected Patch count - perhaps one of the 'original_git_sha' values appears twice")
	}

	result := Cfg{
		JunosVersion:         yamlConfig.Version,
		JunosFamily:          yamlConfig.Family,
		repoYangDirs:         repoYangDirs,
		BaseCacheDir:         yamlConfig.CacheDir,
		JunosConfigFile:      deviceConfigFile,
		YangPatches:          yangPatches,
		cacheRefreshInterval: time.Duration(*yamlConfig.CacheRefreshInterval) * time.Second,
		YamlRepoInfo: githubInfo{
			Owner: yamlConfig.GithubOwnerName,
			Name:  yamlConfig.GithubRepoName,
			Ref:   yamlConfig.GitRef,
		},
	}

	return result, nil
}

func (o *Cfg) RepoYangDirs() []string {
	return o.repoYangDirs
}

func (o *Cfg) cacheTouchFile() string {
	return path.Join(o.YangCacheDir(), cacheFreshnessMarker)
}

func (o *Cfg) UpdateCacheFreshnessMarker() {
	fn := o.cacheTouchFile()

	_, err := os.Stat(fn)
	if os.IsNotExist(err) {
		file, err := os.Create(fn)
		if err != nil {
			return
		}

		file.Close()
		return
	}

	now := time.Now()
	_ = os.Chtimes(fn, now, now)
}

func (o *Cfg) CacheIsFresh() bool {
	fi, err := os.Stat(o.cacheTouchFile())
	if err != nil {
		return false
	}

	if time.Now().Sub(fi.ModTime()) > o.cacheRefreshInterval {
		return false
	}

	return true
}

func (o *Cfg) ConfigRoot() string {
	return fmt.Sprintf(configRoot, strings.TrimSuffix(o.JunosFamily, "-evo"))
}
