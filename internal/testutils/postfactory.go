package testutils

import (
	"fmt"
	"html/template"
	"math/rand"
	"time"

	"github.com/joshhartwig/blogo/internal/models"
)

type PostOption func(*models.Post)

func WithTitle(title string) PostOption {
	return func(p *models.Post) {
		p.Metadata.Title = title
	}
}

func WithSlug(slug string) PostOption {
	return func(p *models.Post) {
		p.Metadata.Slug = slug
	}
}

func WithSummary(summary string) PostOption {
	return func(p *models.Post) {
		p.Metadata.Summary = summary
	}
}

func WithContent(content string) PostOption {
	return func(p *models.Post) {
		p.Content = template.HTML(content)
	}
}

func WithDate(t time.Time) PostOption {
	return func(p *models.Post) {
		p.Metadata.Date = t
	}
}

func WithTags(tags ...string) PostOption {
	return func(p *models.Post) {
		p.Metadata.Tags = append([]string(nil), tags...)
	}
}

func WithReadingDuration(mins int) PostOption {
	return func(p *models.Post) {
		p.Metadata.ReadingDuration = mins
	}
}

// NewPost creates a new models.Post instance with default metadata and content.
// Optional PostOption functions can be provided to customize the post fields.
// If the Slug field is empty after applying options, it will be set to a unique value
// based on the current timestamp.
//
// Example usage:
//
//	post := NewPost(WithTitle("Custom Title"), WithTags([]string{"go", "example"}))
func NewPost(opts ...PostOption) models.Post {
	post := models.Post{
		Metadata: models.PostMetadata{
			Title:           "Test Post",
			Slug:            "test-post",
			Summary:         "Test Summary",
			Tags:            []string{"#testing"},
			Date:            time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			ReadingDuration: 1,
		},
		Content: template.HTML("<p>test content</p>"),
	}

	for _, opt := range opts {
		opt(&post)
	}

	if post.Metadata.Slug == "" {
		post.Metadata.Slug = fmt.Sprintf("post-%d", time.Now().UnixNano())
	}

	return post
}

// WithRandomDate sets a random date within the last year (adjust range as needed).
func WithRandomDate() PostOption {
	return func(p *models.Post) {
		max := time.Now()
		min := max.AddDate(-1, 0, 0) // one year ago
		p.Metadata.Date = randomTime(min, max)
	}
}

func randomTime(min, max time.Time) time.Time {
	if max.Before(min) {
		return min
	}
	delta := max.Unix() - min.Unix()
	sec := rand.Int63n(delta + 1)
	return time.Unix(min.Unix()+sec, 0).UTC()
}
