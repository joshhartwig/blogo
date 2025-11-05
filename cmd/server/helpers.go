package main

import (
	"bytes"
	"fmt"
	"math"
	"net/http"
	"strings"
)

// render executes a template with the passed in data
func (app *application) render(w http.ResponseWriter, status int, templateName string, data any) {
	v, ok := app.templateCache[templateName]
	if !ok {
		err := fmt.Errorf("template not found: %s", templateName)
		app.logger.Error(err.Error())
	}

	buf := new(bytes.Buffer)

	err := v.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.logger.Error("template execution failed", "error", err.Error())
	}

	w.WriteHeader(status)

	// write contents of buffer to response writer
	buf.WriteTo(w)
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
