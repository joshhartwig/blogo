package templates

import (
	"errors"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
	"time"
)

// TemplateCache creates a cache of parsed HTML templates for the application.
// It scans the "templates/pages" directory for all page templates, and for each page:
//   - Parses the "templates/base.html" layout file.
//   - Parses all partial templates in "templates/partials".
//   - Parses the specific page template.
//
// The resulting *template.Template for each page is stored in a map, keyed by the page's filename.
// Returns the cache map and any error encountered during parsing.
func NewCache(fileSys fs.FS, pagesGlobPath, partialsGlobPath, basePath string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// get all the pages from the pages path
	pages, err := fs.Glob(fileSys, pagesGlobPath)
	if err != nil {
		return nil, err
	}

	// iterate through each page
	for _, page := range pages {
		fileName := filepath.Base(page) // strip the name from the filepath ex home.html
		name := strings.TrimSuffix(fileName, filepath.Ext(fileName))

		if name == "" {
			return nil, errors.New("template page name cannot be empty")
		}

		// parse the base template
		ts, err := template.New(name).Funcs(functions).ParseFS(fileSys, basePath)
		if err != nil {
			return nil, err
		}

		// parse the partials
		ts, err = ts.ParseFS(fileSys, partialsGlobPath)
		if err != nil {
			return nil, err
		}

		// parse the page using the provided filesystem so relative paths resolve
		ts, err = ts.ParseFS(fileSys, page)
		if err != nil {
			return nil, err
		}

		// add to cache
		cache[name] = ts
	}

	return cache, nil
}

// return a nice formatted date string
func humanDate(t time.Time) string {

	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006")
}

// init a funcmap obj and store it in a global variable
var functions = template.FuncMap{
	"humanDate": humanDate,
}
