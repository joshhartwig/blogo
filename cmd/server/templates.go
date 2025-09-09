package main

import (
	"errors"
	"html/template"
	"path/filepath"
	"strings"
)

// TemplateCache creates a cache of parsed HTML templates for the application.
// It scans the "templates/pages" directory for all page templates, and for each page:
//   - Parses the "templates/base.html" layout file.
//   - Parses all partial templates in "templates/partials".
//   - Parses the specific page template.
//
// The resulting *template.Template for each page is stored in a map, keyed by the page's filename.
// Returns the cache map and any error encountered during parsing.
func TemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// get all the pages from the pages path
	pages, err := filepath.Glob("./ui/templates/pages/*.html")
	if err != nil {
		return nil, err
	}

	// iterate through each page
	for _, page := range pages {
		fileName := filepath.Base(page) // strip the name from the filepath ex home.html
		names := strings.Split(fileName, ".")
		name := names[0]

		if len(name) == 0 {
			return nil, errors.New("the page file name does not contain a . when trying to split")
		}

		// parse the base template
		ts, err := template.New(name).ParseFiles("./ui/templates/base.html")
		if err != nil {
			return nil, err
		}

		// parse the partials
		ts, err = ts.ParseGlob("./ui/templates/partials/*.html")
		if err != nil {
			return nil, err
		}

		// parse the page
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// add to cache
		cache[name] = ts
	}

	return cache, nil
}
