package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/joshhartwig/blogo/internal/models"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
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

func readMarkdownContent() (map[string]models.Post, error) {
	fmt.Println("starting read markdown content")
	markdown := make(map[string]models.Post)
	files, err := filepath.Glob("./content/*.md")
	if err != nil {
		return nil, err
	}

	for i, file := range files {
		fmt.Printf("processing markdown file %d: %s \n", i, file)
		if file == "" {
			fmt.Printf("file name is empty returning")
			return markdown, errors.New("empty file name")
		}

		filename := filepath.Base(file)             // get just the base file example.html
		name := strings.TrimSuffix(filename, ".md") // get the title name for the map
		data, err := os.ReadFile(file)              // read the file into a []byte
		if err != nil {
			return nil, err
		}

		post, err := convertMarkdownToHtml(name, data) // convert the markdown file into a post struct
		if err != nil {
			return nil, err
		}
		markdown[name] = post
	}

	return markdown, nil
}

func convertMarkdownToHtml(slug string, data []byte) (models.Post, error) {
	post := models.Post{}
	meta := models.PostMetadata{}

	out := bytes.Buffer{}
	gm := goldmark.New(goldmark.WithExtensions(&frontmatter.Extender{}))
	ctx := parser.NewContext()

	err := gm.Convert(data, &out, parser.WithContext(ctx))
	if err != nil {
		return post, err
	}

	d := frontmatter.Get(ctx)

	if err := d.Decode(&meta); err != nil {
		return post, err
	}

	post.Content = template.HTML(out.String())
	post.Metadata = meta
	post.Metadata.ReadingDuration = calculateDuration(out.String()) // calculate reading duration
	post.Metadata.Slug = slug                                       // add the slug

	return post, nil
}

// calculateDuration estimates the reading duration in minutes for the given content string.
// It assumes an average reading speed of 200 words per minute.
// The function returns the rounded duration as an integer.
// If the content is empty, it returns 0.
func calculateDuration(content string) int {
	const wordsPerMinute = 200

	words := strings.Fields(content)
	if len(words) == 0 {
		return 0
	}

	duration := float64(len(words)) / float64(wordsPerMinute)
	return int(math.Round(duration))
}
