package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	defaultGithubOwnerName = "Juniper"
	defaultGithubRepoName  = "yang"
	defaultConfigFile      = "./config.yaml"
	yangCommonDir          = "common"
	yangConfDir            = "conf"
	junosVersionRegexp     = `^((\d\d)\.\d+)R\d+.*`
)

var flagC = flag.String("c", defaultConfigFile, "YAML config file")

type patch struct {
	OriginalGitSha string `yaml:"original_git_sha"`
	RequiredSha256 string `yaml:"required_sha_256"`
	Patch          string `yaml:"diff"`
}

type githubInfo struct {
	Owner string
	Name  string
	Ref   string
}

type jtafConfig struct {
	yamlRepoInfo     githubInfo
	DeviceConfigFile string
	JunosVersion     string
	JunosFamily      string
	BaseCacheDir     string
	YangPatches      map[string]patch
	repoYangDir      string
}

// mkYangCacheDir returns the yang cache dir and returns whether it needed to be created
func (o *jtafConfig) mkYangCacheDir() (string, bool) {
	yangCacheDir := o.yangCacheDir()

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

// repoPath returns a string like "github.com/Juniper/yang" or "github.com/Juniper/yang@main"
func (o *jtafConfig) repoPath() string {
	var gitRef string
	if o.yamlRepoInfo.Ref != "" {
		gitRef = "@" + o.yamlRepoInfo.Ref
	}
	return path.Join("github.com", o.yamlRepoInfo.Owner, o.yamlRepoInfo.Name) + gitRef
}

// yangCacheDir returns the path to the top-level yang cache directory
func (o *jtafConfig) yangCacheDir() string {
	return path.Clean(path.Join(o.BaseCacheDir, "yang"))
}

func (o *jtafConfig) junosYangCacheDir() string {
	ref := o.yamlRepoInfo.Ref
	if ref != "" && !strings.HasPrefix(ref, "@") {
		ref = "@" + ref
	}
	return path.Clean(path.Join(o.yangCacheDir(), "github.com", o.yamlRepoInfo.Owner, o.yamlRepoInfo.Name, ref))
}

func (o *jtafConfig) otherYangCacheDir() string {
	return path.Clean(path.Join(o.yangCacheDir(), "ietf"))
}

func getConfig() (jtafConfig, error) {
	flag.Parse()

	cfgBytes, err := os.ReadFile(*flagC)
	if err != nil {
		return jtafConfig{}, fmt.Errorf("while reading config file %q - %w", *flagC, err)
	}

	var yamlConfig struct {
		DeviceConfigFile string  `yaml:"junos_config_xml"`
		Family           string  `yaml:"junos_family"`
		Version          string  `yaml:"junos_version"`
		GitRef           string  `yaml:"git_ref"`
		GithubOwnerName  string  `yaml:"github_owner_name"`
		GithubRepoName   string  `yaml:"github_repo_name"`
		CacheDir         string  `yaml:"cache_dir"`
		YangPatches      []patch `yaml:"yang_patches"`
	}

	err = yaml.Unmarshal(cfgBytes, &yamlConfig)
	if err != nil {
		return jtafConfig{}, fmt.Errorf("while parsing config file %q - %w", *flagC, err)
	}

	if yamlConfig.GithubOwnerName == "" {
		yamlConfig.GithubOwnerName = defaultGithubOwnerName
	}

	if yamlConfig.GithubRepoName == "" {
		yamlConfig.GithubRepoName = defaultGithubRepoName
	}

	deviceConfigFile, err := filepath.Abs(yamlConfig.DeviceConfigFile)
	if err != nil {
		return jtafConfig{}, fmt.Errorf("while expanding file path %q - %w", yamlConfig.DeviceConfigFile, err)
	}

	family := osFamilies.Parse(yamlConfig.Family)
	if family == nil {
		return jtafConfig{}, fmt.Errorf("family must be one of %s, got %q", osFamilies.Members(), yamlConfig.Family)
	}

	s := regexp.MustCompile(junosVersionRegexp).FindStringSubmatch(yamlConfig.Version)
	if len(s) != 3 {
		return jtafConfig{}, fmt.Errorf("failed to parse junos version %q - version should look like 23.1R1 or 23.4R1.10", yamlConfig.Version)
	}

	yangPatches := make(map[string]patch, len(yamlConfig.YangPatches))
	for _, p := range yamlConfig.YangPatches {
		yangPatches[p.OriginalGitSha] = p
	}
	if len(yangPatches) != len(yamlConfig.YangPatches) {
		return jtafConfig{}, fmt.Errorf("unexpected patch count - perhaps one of the 'original_git_sha' values appears twice")
	}

	result := jtafConfig{
		JunosVersion:     yamlConfig.Version,
		JunosFamily:      family.Value,
		repoYangDir:      path.Join(s[1], yamlConfig.Version),
		BaseCacheDir:     yamlConfig.CacheDir,
		DeviceConfigFile: deviceConfigFile,
		YangPatches:      yangPatches,
		yamlRepoInfo: githubInfo{
			Owner: yamlConfig.GithubOwnerName,
			Name:  yamlConfig.GithubRepoName,
			Ref:   yamlConfig.GitRef,
		},
	}

	return result, nil
}

func (o *jtafConfig) RepoDirYangCommon() string {
	return path.Join(o.repoYangDir, yangCommonDir)
}

func (o *jtafConfig) RepoDirYangFamily() string {
	return path.Join(o.repoYangDir, o.JunosFamily, yangConfDir)
}
