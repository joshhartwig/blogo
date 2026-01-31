package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"slices"
	"time"

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

	// we had no frontmatter see issue https://github.com/joshhartwig/blogo/issues/9
	if meta.Title == "" {
		meta.Title = "post missing frontmatter"
		meta.Date = time.Now()
		meta.Draft = false
		meta.Slug = "post missing frontmatter"
		meta.Summary = "post missing frontmatter"
	}

	post.Content = template.HTML(out.String())                      // add content to the post
	post.Metadata = meta                                            // add the frontmatter
	post.Metadata.ReadingDuration = calculateDuration(out.String()) // calculate reading duration
	post.Metadata.Slug = slug                                       // add the slug

	return post, nil
}

func readMarkdownReturnPostsInOrder(fileSystem fs.FS) ([]models.Post, error) {
	var posts []models.Post

	// First, look for folder-based posts (*/index.md)
	folderFiles, err := fs.Glob(fileSystem, "*/index.md")
	if err != nil {
		return nil, err
	}

	for _, path := range folderFiles {
		data, err := fs.ReadFile(fileSystem, path)
		if err != nil {
			return nil, err
		}

		slug := slugify(filepath.Dir(path))
		post, err := convertMarkdownToHtml(slug, data)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	// Also check for flat .md files (backwards compatibility)
	flatFiles, err := fs.Glob(fileSystem, "*.md")
	if err != nil {
		return nil, err
	}

	for _, path := range flatFiles {
		data, err := fs.ReadFile(fileSystem, path)
		if err != nil {
			return nil, err
		}

		slug := slugify(filepath.Base(path))
		post, err := convertMarkdownToHtml(slug, data)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if len(posts) == 0 {
		return nil, fmt.Errorf("no matching files")
	}

	slices.SortFunc(posts, func(a, b models.Post) int {
		return a.Metadata.Date.Compare(b.Metadata.Date)
	})

	return posts, nil
}

// readMarkdownContent reads the local content directory for .md files, converts them to a Post struct and
// adds them to a template cache. It supports both folder-based posts (post-name/index.md) and
// flat .md files for backwards compatibility.
func readMarkdownContent(fileSystem fs.FS) (map[string]models.Post, error) {
	posts := make(map[string]models.Post)

	// First, look for folder-based posts (*/index.md)
	folderFiles, err := fs.Glob(fileSystem, "*/index.md")
	if err != nil {
		return nil, err
	}

	for _, path := range folderFiles {
		data, err := fs.ReadFile(fileSystem, path)
		if err != nil {
			return nil, err
		}

		// For folder-based posts, use the folder name as the slug
		slug := slugify(filepath.Dir(path))

		post, err := convertMarkdownToHtml(slug, data)
		if err != nil {
			return nil, err
		}
		if post.Metadata.Draft {
			continue
		}
		posts[slug] = post
	}

	// Also check for flat .md files (backwards compatibility)
	flatFiles, err := fs.Glob(fileSystem, "*.md")
	if err != nil {
		return nil, err
	}

	for _, path := range flatFiles {
		data, err := fs.ReadFile(fileSystem, path)
		if err != nil {
			return nil, err
		}

		slug := slugify(filepath.Base(path))

		post, err := convertMarkdownToHtml(slug, data)
		if err != nil {
			return nil, err
		}
		if post.Metadata.Draft {
			continue
		}
		posts[slug] = post
	}

	if len(posts) == 0 {
		return nil, fmt.Errorf("no matching files")
	}

	return posts, nil
}
