package main

import (
	"fmt"
	"math"
	"net/http"
	"strings"
)

// render executes a template with the passed in data
func (app *application) render(w http.ResponseWriter, templateName string, data any) error {
	v, ok := app.templateCache[templateName]
	if !ok {
		err := fmt.Errorf("template not found: %s", templateName)
		app.logger.Error(err.Error())
		return err
	}

	err := v.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.logger.Error("template execution failed", "error", err.Error())
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
