package config

import (
	"fmt"
	"io/fs"

	"github.com/BurntSushi/toml"
)

type ThemeConfig struct {
	Name        string      `toml:"name"`
	Description string      `toml:"description"`
	Fonts       ThemeFonts  `toml:"fonts"`
	Colors      ThemeColors `toml:"colors"`
	Tokens      ThemeTokens `toml:"tokens"`
}

type ThemeFonts struct {
	Sans           string `toml:"sans"`
	Mono           string `toml:"mono"`
	GoogleFontsURL string `toml:"google_fonts_url"`
}

type ThemeColors struct {
	Bg        string `toml:"bg"`
	BgElev    string `toml:"bg_elev"`
	Fg        string `toml:"fg"`
	Muted     string `toml:"muted"`
	Border    string `toml:"border"`
	Link      string `toml:"link"`
	LinkHover string `toml:"link_hover"`
}

type ThemeTokens struct {
	Radius string `toml:"radius"`
	Shadow string `toml:"shadow"`
}

// LoadTheme reads themes/<name>/theme.toml from the provided FS.
func LoadTheme(fsys fs.FS, name string) (ThemeConfig, error) {
	f, err := fsys.Open("themes/" + name + "/theme.toml")
	if err != nil {
		return ThemeConfig{}, fmt.Errorf("opening theme %q: %w", name, err)
	}
	defer f.Close()
	var cfg ThemeConfig
	if _, err = toml.NewDecoder(f).Decode(&cfg); err != nil {
		return ThemeConfig{}, fmt.Errorf("decoding theme %q: %w", name, err)
	}
	return cfg, nil
}
