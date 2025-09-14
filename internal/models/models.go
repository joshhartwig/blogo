package models

import (
	"html/template"
	"time"
)

type TemplateData struct {
	Title string
	Posts []Post
}

type Post struct {
	Metadata PostMetadata
	Content  template.HTML
}

type PostMetadata struct {
	Title string    `yaml:"title"`
	Date  time.Time `yaml:"date"`
	Tags  []string  `yaml:"tags"`
	Desc  string    `yaml:"description"`
	Slug  string
}
