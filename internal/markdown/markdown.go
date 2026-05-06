package markdown

import (
	"bytes"
	"html/template"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
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

// readMarkdownContent reads the local content directory for .md files, converts them to a Post struct and
// adds them to a template cache. It supports both folder-based posts (post-name/index.md) and
// flat .md files for backwards compatibility.
func GetMarkdownFromFS(fileSystem fs.FS) []models.Post {
	posts := []models.Post{}

	filenames, err := fs.Glob(fileSystem, "*/index.md")
	if err != nil {
		return posts
	}

	// iterate through each filename and read the file
	for _, path := range filenames {
		data, err := fs.ReadFile(fileSystem, path)
		if err != nil {
			return posts
		}

		// For folder-based posts, use the folder name as the slug
		slug := slugify(filepath.Dir(path))

		post, err := convertMarkdownToHtml(slug, data)
		if err != nil {
			return posts
		}
		if post.Metadata.Draft {
			continue
		}

		post.Metadata.Slug = slug
		posts = append(posts, post)
	}

	// Also check for flat .md files (backwards compatibility)
	flatFiles, err := fs.Glob(fileSystem, "*.md")
	if err != nil {
		return posts
	}

	for _, path := range flatFiles {
		data, err := fs.ReadFile(fileSystem, path)
		if err != nil {
			return posts
		}

		slug := slugify(filepath.Base(path))

		post, err := convertMarkdownToHtml(slug, data)
		if err != nil {
			return posts
		}
		if post.Metadata.Draft {
			continue
		}

		post.Metadata.Slug = slug
		posts = append(posts, post)
	}

	if len(posts) == 0 {
		return posts
	}

	return posts
}

// calculateDuration estimates the reading duration in minutes for the given content string.
// It assumes an average reading speed of 200 words per minute.
func calculateDuration(content string) int {
	const wordsPerMinute = 200
	words := strings.Fields(content)
	if len(words) == 0 {
		return 1
	}
	duration := (len(words) / wordsPerMinute) + 1
	return duration
}

// slugify generates a URL-friendly slug from a string
func slugify(s string) string {
	// Remove file extension if present
	if dot := strings.LastIndex(s, "."); dot != -1 {
		s = s[:dot]
	}
	s = strings.ToLower(s)                   // lower
	s = strings.ReplaceAll(s, " ", "-")      // replace spaces with dash
	s = strings.ReplaceAll(s, "_", "-")      // underscores with -
	re := regexp.MustCompile(`[^a-z0-9\-]+`) // replace lower case, digits, hypens
	s = re.ReplaceAllString(s, "")
	// Collapse multiple consecutive dashes into single dash
	reDash := regexp.MustCompile(`-+`)
	s = reDash.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}
