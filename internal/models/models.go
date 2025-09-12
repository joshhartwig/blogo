package models

import "html/template"

type TemplateData struct {
	Title string
	Posts []Post
}

type Post struct {
	Metadata PostMetadata
	Content  template.HTML
}

type PostMetadata struct {
	Title string   `yaml:"title"`
	Tags  []string `yaml:"tags"`
	Desc  string   `yaml:"description"`
}
