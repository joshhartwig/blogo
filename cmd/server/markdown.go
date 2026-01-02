package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/joshhartwig/blogo/internal/models"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

// convertMarkdownToHtml takes the markdown data in bytes and converts it to html
// it then creates a new Post{} and adds both the HTML and PostMetadata and returns both
func convertMarkdownToHtml(slug string, data []byte) (models.Post, error) {
	post := models.Post{}
	meta := models.PostMetadata{}

	out := bytes.Buffer{}
	gm := goldmark.New(goldmark.WithExtensions(
		&frontmatter.Extender{}, highlighting.NewHighlighting(
			highlighting.WithStyle("monokai"),
			highlighting.WithFormatOptions(chromahtml.WithLineNumbers(true)),
		)))
	ctx := parser.NewContext()

	err := gm.Convert(data, &out, parser.WithContext(ctx))
	if err != nil {
		return post, err
	}

	d := frontmatter.Get(ctx)
	if d != nil {
		if err := d.Decode(&meta); err != nil {
			return post, err
		}
	}

	post.Content = template.HTML(out.String())                      // add content to the post
	post.Metadata = meta                                            // add the frontmatter
	post.Metadata.ReadingDuration = calculateDuration(out.String()) // calculate reading duration
	post.Metadata.Slug = slug                                       // add the slug

	return post, nil
}

func readMarkdownReturnPostsInOrder(fileSystem fs.FS) ([]models.Post, error) {
	var posts []models.Post

	files, err := fs.Glob(fileSystem, "*.md")
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no matching files")
	}

	for _, path := range files {
		data, err := fs.ReadFile(fileSystem, path)
		if err != nil {
			return nil, err
		}

		slug := strings.TrimSuffix(filepath.Base(path), ".md") // get the title name for the map
		post, err := convertMarkdownToHtml(slug, data)         // convert the markdown file into a post struct
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	slices.SortFunc(posts, func(a, b models.Post) int {
		return a.Metadata.Date.Compare(b.Metadata.Date)
	})

	return posts, nil
}

// readMarkdownContent reads the local content directory for .md files, converts them to a Post struct and
// adds them to a template cache
func readMarkdownContent(fileSystem fs.FS) (map[string]models.Post, error) {
	posts := make(map[string]models.Post)

	files, err := fs.Glob(fileSystem, "*.md") // fetch all .md files in the content dir
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no matching files")
	}

	for _, path := range files {
		data, err := fs.ReadFile(fileSystem, path)
		if err != nil {
			return nil, err
		}

		slug := strings.TrimSuffix(filepath.Base(path), ".md") // get the title name for the map
		post, err := convertMarkdownToHtml(slug, data)         // convert the markdown file into a post struct
		if err != nil {
			return nil, err
		}
		if post.Metadata.Draft {
			// skip draft posts
			continue
		}
		posts[slug] = post // add markdown to cache
	}

	return posts, nil
}
