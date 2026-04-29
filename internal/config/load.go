package config

import (
	"embed"
	"os"

	"github.com/BurntSushi/toml"
)

//go:embed default.toml
var defaultConfigFS embed.FS

func Load(path string) (SiteConfig, error) {
	var sc SiteConfig

	if path == "" {
		return loadDefault()
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return sc, err
	}

	_, err = toml.Decode(string(bytes), &sc)
	if err != nil {
		return sc, err
	}
	return sc, nil
}

func LoadWithDefaults(path string) (SiteConfig, error) {
	cfg, err := loadDefault()
	if err != nil {
		return SiteConfig{}, err
	}
	if path == "" {
		return cfg, nil
	}
	userCfg, err := Load(path)
	if err != nil {
		return SiteConfig{}, err
	}
	return Merge(cfg, userCfg), nil
}

func Merge(def, user SiteConfig) SiteConfig {
	var combined SiteConfig

	if user.SiteTitle == "" {
		combined.SiteTitle = def.SiteTitle
	}

	if user.BaseURL == "" {
		combined.BaseURL = def.BaseURL
	}

	if user.Theme == "" {
		combined.Theme = def.Theme
	}

	if user.ContentDir == "" {
		combined.ContentDir = def.ContentDir
	}

	if user.Description == "" {
		combined.Description = def.Description
	}

	if user.PostsPerPage == 0 {
		combined.PostsPerPage = def.PostsPerPage
	}

	if user.Author.Bio == "" {
		combined.Author.Bio = def.Author.Bio
	}

	if user.Author.Name == "" {
		combined.Author.Name = def.Author.Name
	}

	return combined
}

// loads the embedded default toml
func loadDefault() (SiteConfig, error) {
	var sc SiteConfig
	bytes, err := defaultConfigFS.ReadFile("default.toml")
	if err != nil {
		return sc, err
	}

	_, err = toml.Decode(string(bytes), &sc)
	if err != nil {
		return sc, err
	}

	return sc, nil
}
