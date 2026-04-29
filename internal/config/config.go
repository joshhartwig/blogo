package config

type SiteConfig struct {
	SiteTitle    string       `toml:"site_title" yaml:"site_title"`
	BaseURL      string       `toml:"base_url" yaml:"base_url"`
	Description  string       `toml:"description" yaml:"description"`
	Theme        string       `toml:"theme" yaml:"theme"`
	ContentDir   string       `toml:"content_dir" yaml:"content_dir"`
	PostsPerPage int          `toml:"posts_per_page" yaml:"posts_per_page"`
	Author       AuthorConfig `toml:"author" yaml:"author"`
}

type AuthorConfig struct {
	Name string `toml:"name" yaml:"name"`
	Bio  string `toml:"bio" yaml:"bio"`
}
