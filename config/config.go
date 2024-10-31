package jtafCfg

import (
	"flag"
	"fmt"
	"github.com/google/go-github/v66/github"
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

type githubInfo struct {
	Owner string
	Name  string
	Ref   string
}

type Cfg struct {
	YamlRepoInfo    githubInfo
	JunosConfigFile string
	JunosVersion    string
	JunosFamily     string
	BaseCacheDir    string
	YangPatches     map[string]Patch
	repoYangDir     string
}

func (o Cfg) Patch(filename string, grc github.RepositoryContent) error {
	p, ok := o.YangPatches[*grc.SHA]
	if !ok {
		return nil // no patch; no problem
	}

	return p.applyToFile(filename)
}

// TargetFileName returns the filesystem location where the github.RepositoryContent should be stored.
func (o Cfg) TargetFileName(grc github.RepositoryContent) string {
	return path.Clean(path.Join(o.JunosYangCacheDir(), *grc.Path))
}

// mkYangCacheDir returns the yang cache dir and returns whether it needed to be created
func (o *Cfg) mkYangCacheDir() (string, bool) {
	yangCacheDir := o.YangCacheDir()

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

// YangCacheDir returns the path to the top-level yang cache directory
func (o *Cfg) YangCacheDir() string {
	return path.Clean(path.Join(o.BaseCacheDir, "yang"))
}

func (o *Cfg) JunosYangCacheDir() string {
	ref := o.YamlRepoInfo.Ref
	if ref != "" && !strings.HasPrefix(ref, "@") {
		ref = "@" + ref
	}
	return path.Clean(path.Join(o.YangCacheDir(), "github.com", o.YamlRepoInfo.Owner, o.YamlRepoInfo.Name, ref))
}

func (o *Cfg) otherYangCacheDir() string {
	return path.Clean(path.Join(o.YangCacheDir(), "ietf"))
}

func Get() (Cfg, error) {
	flag.Parse()

	cfgBytes, err := os.ReadFile(*flagC)
	if err != nil {
		return Cfg{}, fmt.Errorf("while reading config file %q - %w", *flagC, err)
	}

	var yamlConfig struct {
		DeviceConfigFile string  `yaml:"junos_config_xml"`
		Family           string  `yaml:"junos_family"`
		Version          string  `yaml:"junos_version"`
		GitRef           string  `yaml:"git_ref"`
		GithubOwnerName  string  `yaml:"github_owner_name"`
		GithubRepoName   string  `yaml:"github_repo_name"`
		CacheDir         string  `yaml:"cache_dir"`
		YangPatches      []Patch `yaml:"yang_patches"`
	}

	err = yaml.Unmarshal(cfgBytes, &yamlConfig)
	if err != nil {
		return Cfg{}, fmt.Errorf("while parsing config file %q - %w", *flagC, err)
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

	family := osFamilies.Parse(yamlConfig.Family)
	if family == nil {
		return Cfg{}, fmt.Errorf("family must be one of %s, got %q", osFamilies.Members(), yamlConfig.Family)
	}

	s := regexp.MustCompile(junosVersionRegexp).FindStringSubmatch(yamlConfig.Version)
	if len(s) != 3 {
		return Cfg{}, fmt.Errorf("failed to parse junos version %q - version should look like 23.1R1 or 23.4R1.10", yamlConfig.Version)
	}

	yangPatches := make(map[string]Patch, len(yamlConfig.YangPatches))
	for _, p := range yamlConfig.YangPatches {
		yangPatches[p.OriginalGitSha] = p
	}
	if len(yangPatches) != len(yamlConfig.YangPatches) {
		return Cfg{}, fmt.Errorf("unexpected Patch count - perhaps one of the 'original_git_sha' values appears twice")
	}

	result := Cfg{
		JunosVersion:    yamlConfig.Version,
		JunosFamily:     family.Value,
		repoYangDir:     path.Join(s[1], yamlConfig.Version),
		BaseCacheDir:    yamlConfig.CacheDir,
		JunosConfigFile: deviceConfigFile,
		YangPatches:     yangPatches,
		YamlRepoInfo: githubInfo{
			Owner: yamlConfig.GithubOwnerName,
			Name:  yamlConfig.GithubRepoName,
			Ref:   yamlConfig.GitRef,
		},
	}

	return result, nil
}

func (o *Cfg) RepoDirYangCommon() string {
	return path.Join(o.repoYangDir, yangCommonDir)
}

func (o *Cfg) RepoDirYangFamily() string {
	return path.Join(o.repoYangDir, o.JunosFamily, yangConfDir)
}
