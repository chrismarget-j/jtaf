package jtafCfg

type Patch struct {
	OriginalGitSha string `yaml:"original_git_sha"`
	RequiredSha256 string `yaml:"required_sha_256"`
	Patch          string `yaml:"diff"`
}
