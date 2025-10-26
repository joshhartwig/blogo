package main

import (
	"testing"
	"testing/fstest"
)

func TestTemplateCache(t *testing.T) {
	testFs := fstest.MapFS{
		"templates/pages/home.html":      {Data: []byte(`<h1>Home Page</h1>`)},
		"templates/pages/about.html":     {Data: []byte(`<h1>About Page</h1>`)},
		"templates/base.html":            {Data: []byte(`<!DOCTYPE html><html><body>{{template "content" .}}</body></html>`)},
		"templates/partials/nav.html":    {Data: []byte(`<nav>Navigation</nav>`)},
		"templates/partials/footer.html": {Data: []byte(`<footer>Footer</footer>`)},
	}

	cache, err := LoadTemplatesAsMap(testFs, "templates/pages/*.html", "templates/partials/*.html", "templates/base.html")
	if err != nil {
		t.Errorf("caught error testing template cache reading %v", err)
	}

	if len(cache) < 2 {
		t.Errorf("expected at least 2 templates in cache, got %d", len(cache))
	}

	if _, exists := cache["home"]; !exists {
		t.Error("Expected 'home' template in cache")
	}
}
