package models

type TemplateData struct {
	Title    string
	Articles []Article
}

type Article struct {
	Title   string
	Content []byte
}
