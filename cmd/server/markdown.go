package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/joshhartwig/blogo/internal/models"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

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

	post.Content = template.HTML(out.String())                      // add content to the post
	post.Metadata = meta                                            // add the frontmatter
	post.Metadata.ReadingDuration = calculateDuration(out.String()) // calculate reading duration
	post.Metadata.Slug = slug                                       // add the slug

	return post, nil
}

// readMarkdownContent reads the local content directory for .md files, converts them to a Post struct and
// adds them to a template cache
func readMarkdownContent(path string) (map[string]models.Post, error) {
	fmt.Println("Reading markdown files:")
	markdown := make(map[string]models.Post)   // create a post cache
	files, err := filepath.Glob(path + "*.md") // fetch all .md files in the content dir
	if err != nil {
		return nil, err
	}

	for i, file := range files {
		fmt.Printf("\t- %d:%s\n", i, file)
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
		markdown[name] = post // add markdown to cache
	}

	return markdown, nil
}
