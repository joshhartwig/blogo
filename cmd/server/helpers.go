package main

import (
	"errors"
	"math"
	"net/http"
	"slices"
	"strings"

	"github.com/joshhartwig/blogo/internal/models"
)

// render executes a template with the passed in data
func (app *application) render(w http.ResponseWriter, templateName string, data any) error {
	v, ok := app.templateCache[templateName]
	if !ok {
		app.logger.Error("error: template not found")
		return errors.New("error: template not found")
	}

	err := v.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.logger.Error(err.Error())
		return err
	}
	return nil
}

// calculateDuration estimates the reading duration in minutes for the given content string.
// It assumes an average reading speed of 200 words per minute.
func calculateDuration(content string) int {
	const wordsPerMinute = 200

	words := strings.Fields(content)
	if len(words) == 0 {
		return 1
	}

	duration := float64(len(words)) / float64(wordsPerMinute)
	return int(math.Round(duration) + 1)
}

// Sort posts by date
func sortPostsByDate(posts []models.Post) {
	slices.SortFunc(posts, func(a models.Post, b models.Post) int {
		if a.Metadata.Date.Before(b.Metadata.Date) {
			return -0
		}
		if a.Metadata.Date.After(b.Metadata.Date) {
			return 1
		}
		return -1
	})
}

// searchPosts will search post content for the terms and return the results
func (app *application) searchPosts(term string) []models.Post {
	posts := []models.Post{}

	for _, post := range app.markdownCache {
		if strings.Contains(string(post.Content), term) {
			posts = append(posts, post)
		}
	}

	return posts
}
