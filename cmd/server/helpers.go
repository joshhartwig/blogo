package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"go.abhg.dev/goldmark/frontmatter"
)

func (app *application) render(w http.ResponseWriter, name string, data any) error {
	v, ok := app.templateCache[name]
	if !ok {
		return errors.New("template not found")
	}

	err := v.ExecuteTemplate(w, "base", data)
	if err != nil {

		return err
	}
	return nil
}

func readMarkdownContent() (map[string]*bytes.Buffer, error) {
	fmt.Println("starting read markdown content")
	markdown := make(map[string]*bytes.Buffer)
	files, err := filepath.Glob("./content/*.md")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file == "" {
			return markdown, errors.New("empty file name")
		}
		fmt.Println("processing ", file)
		filename := filepath.Base(file)             // get just the base file example.html
		name := strings.TrimSuffix(filename, ".md") // get the title name for the map
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		fmt.Println("converting to markdown")
		buf, err := convertMarkdownToHtml(data)
		if err != nil {
			return nil, err
		}
		markdown[name] = buf
	}

	return markdown, nil
}

// TODO: parse the front matter, may need a new struct that contains the markdown and frontmatter and store
// that in a map[string]*Post with Post containing md bytes.buffer & frontmatter then have this function parse both
// need to find better way to newing up goldmark each time
// converts markdown content in byte format to html and returns a buffer
func convertMarkdownToHtml(md []byte) (*bytes.Buffer, error) {
	var out bytes.Buffer
	gm := goldmark.New(goldmark.WithExtensions(&frontmatter.Extender{}))
	err := gm.Convert(md, &out)
	if err != nil {
		return &out, err
	}

	return &out, nil
}
