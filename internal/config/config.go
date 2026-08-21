package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const DefaultMaxFileBytes = 1 << 20

type Config struct {
	Allowlist    Allowlist `yaml:"allowlist"`
	Exclude      Exclude   `yaml:"exclude"`
	MaxFileBytes int64     `yaml:"max_file_bytes"`
}

type Allowlist struct {
	Paths    []string `yaml:"paths"`
	Rules    []string `yaml:"rules"`
	Patterns []string `yaml:"patterns"`
}

type Exclude struct {
	Paths []string `yaml:"paths"`
}

func Default() Config {
	return Config{MaxFileBytes: DefaultMaxFileBytes}
}

func Load(path string) (Config, error) {
	cfg := Default()
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("설정 파일 읽기 %s : %w", path, err)
	}
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return Config{}, fmt.Errorf("설정 파일 파싱 %s : %w", path, err)
	}
	if cfg.MaxFileBytes <= 0 {
		cfg.MaxFileBytes = DefaultMaxFileBytes
	}
	return cfg, nil
}

func LoadForRoot(root, explicit string) (Config, error) {
	if explicit != "" {
		return Load(explicit)
	}
	candidate := filepath.Join(root, ".ksecret.yml")
	if _, err := os.Stat(candidate); err != nil {
		if os.IsNotExist(err) {
			return Default(), nil
		}
		return Config{}, fmt.Errorf("설정 파일 확인 : %w", err)
	}
	return Load(candidate)
}

func (c Config) PathExcluded(rel string) bool {
	for _, p := range c.Exclude.Paths {
		if MatchPath(p, rel) {
			return true
		}
	}
	return false
}

func (c Config) Allowed(rel, rule, line string) bool {
	for _, p := range c.Allowlist.Paths {
		if MatchPath(p, rel) {
			return true
		}
	}
	for _, r := range c.Allowlist.Rules {
		if r == rule {
			return true
		}
	}
	for _, pat := range c.Allowlist.Patterns {
		if pat != "" && strings.Contains(line, pat) {
			return true
		}
	}
	return false
}

func MatchPath(pattern, rel string) bool {
	rel = filepath.ToSlash(rel)
	pattern = filepath.ToSlash(pattern)
	if pattern == rel {
		return true
	}
	if strings.HasSuffix(pattern, "/**") {
		base := strings.TrimSuffix(pattern, "/**")
		return rel == base || strings.HasPrefix(rel, base+"/")
	}
	if strings.HasPrefix(pattern, "**/") {
		rest := strings.TrimPrefix(pattern, "**/")
		ok, _ := filepath.Match(rest, filepath.Base(rel))
		return ok
	}
	ok, _ := filepath.Match(pattern, rel)
	if ok {
		return true
	}
	ok, _ = filepath.Match(pattern, filepath.Base(rel))
	return ok
}
