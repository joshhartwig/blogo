package models

import (
	"html/template"
	"time"
)

type TemplateData struct {
	Title string
	Posts []Post
}

type PaginationData struct {
	CurrentPage  int
	TotalPages   int
	HasNext      bool
	HasPrev      bool
	NextPage     int
	PrevPage     int
	TotalPosts   int
	PostsPerPage int
}

type HomePageData struct {
	Title          string
	Posts          []Post
	PaginationData PaginationData
}

type Post struct {
	Metadata PostMetadata
	Content  template.HTML
}

type PostMetadata struct {
	Title           string    `yaml:"title"`
	Date            time.Time `yaml:"date"`
	Tags            []string  `yaml:"tags"`
	Summary         string    `yaml:"summary"`
	Slug            string
	ReadingDuration int
}

type RSS struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"`
	PubDate     time.Time `xml:"pubdate"`
	Category    string    `xml:"category"`
	Item        []Item    `xml:"items"`
}

type Item struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Category    string    `xml:"category"`
	GUID        string    `xml:"guid"`
	PubDate     time.Time `xml:"pubDate"`
}
