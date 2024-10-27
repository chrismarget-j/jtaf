package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

const (
	junosVersionRegexp = `^((\d\d)\.\d+)R\d+.*`
)

var (
	flagC = flag.String("c", "config.yaml", "YAML config file")
)

type jtafConfig struct {
	DeviceConfigFile string
	Platform         junosPlatform
	JunosVersion     string
	GitRef           string
	YangDirPlatform  string
	YangDirCommon    string
	GithubOwnerName  string
	GithubRepoName   string
	BaseCacheDir     string
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
	if o.GitRef != "" {
		gitRef = "@" + o.GitRef
	}
	return path.Join("github.com", o.GithubOwnerName, o.GithubRepoName) + gitRef
}

// yangCacheDir returns the path to the yang cache directory for configured git ref, junosPlatform and junos version
func (o *jtafConfig) yangCacheDir() string {
	return path.Join(o.BaseCacheDir, "yang", o.GitRef, o.Platform.Value, o.JunosVersion)
}

func getConfig() (jtafConfig, error) {
	flag.Parse()

	cfgBytes, err := os.ReadFile(*flagC)
	if err != nil {
		return jtafConfig{}, fmt.Errorf("while reading config file %q - %w", *flagC, err)
	}

	var yamlConfig struct {
		DeviceConfigFile string `yaml:"device_config_file"`
		Platform         string `yaml:"platform"`
		JunosVersion     string `yaml:"junos_version"`
		GitRef           string `yaml:"git_ref"`
		GithubOwnerName  string `yaml:"github_owner_name"`
		GithubRepoName   string `yaml:"github_repo_name"`
		CacheDir         string `yaml:"cache_dir"`
	}

	err = yaml.Unmarshal(cfgBytes, &yamlConfig)
	if err != nil {
		return jtafConfig{}, fmt.Errorf("while parsing config file %q - %w", *flagC, err)
	}

	deviceConfigFile, err := filepath.Abs(yamlConfig.DeviceConfigFile)
	if err != nil {
		return jtafConfig{}, fmt.Errorf("while expanding file path %q - %w", yamlConfig.DeviceConfigFile, err)
	}

	p := platforms.Parse(yamlConfig.Platform)
	if p == nil {
		return jtafConfig{}, fmt.Errorf("junosPlatform must be one of %s, got %q", platforms.Members(), yamlConfig.Platform)
	}

	s := regexp.MustCompile(junosVersionRegexp).FindStringSubmatch(yamlConfig.JunosVersion)
	if len(s) != 3 {
		return jtafConfig{}, fmt.Errorf("failed to parse junos version %q - version should look like 23.1R1 or 23.4R1.10", yamlConfig.JunosVersion)
	}

	osFam, ok := platformToOsFamily[*p]
	if !ok {
		return jtafConfig{}, fmt.Errorf("unknown os for junosPlatform %q", p.Value)
	}

	result := jtafConfig{
		Platform:         *p,
		JunosVersion:     yamlConfig.JunosVersion,
		GitRef:           yamlConfig.GitRef,
		YangDirPlatform:  path.Join(s[1], yamlConfig.JunosVersion, osFam.Value, "conf"),
		YangDirCommon:    path.Join(s[1], yamlConfig.JunosVersion, "common"),
		BaseCacheDir:     yamlConfig.CacheDir,
		DeviceConfigFile: deviceConfigFile,
		GithubOwnerName:  yamlConfig.GithubOwnerName,
		GithubRepoName:   yamlConfig.GithubRepoName,
	}

	return result, nil
}
