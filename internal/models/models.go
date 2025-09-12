package models

import "html/template"

type TemplateData struct {
	Title    string
	Articles []Article
}

type Article struct {
	Title   string
	Content []byte
}

type Post struct {
	Title   string
	Content template.HTML
}
