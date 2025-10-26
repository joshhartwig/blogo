package main

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/joshhartwig/blogo/internal/models"
)

func TestPing_ReturnsString(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rr := httptest.NewRecorder()

	ping(rr, req)
	if rr.Code != 200 {
		t.Errorf("wanted status code 200 got %v", rr.Code)
	}

	got := rr.Body.String()
	want := "pong"

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

}

func TestHomeHandler_Returns200(t *testing.T) {
	app, err := returnMockedApp()
	if err != nil {
		t.Errorf("error fetching mocked app")
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	app.homeHandler(rr, req)
	got := rr.Body.String()
	if !strings.Contains(got, "Home") {
		t.Errorf("got %s", got)
	}
}

func returnMockedApp() (application, error) {

	testFs := fstest.MapFS{
		"templates/pages/home.html":      {Data: []byte(`{{define "main"}}<h1>Home Page</h1>{{end}}`)},
		"templates/pages/about.html":     {Data: []byte(`<h1>About Page</h1>`)},
		"templates/base.html":            {Data: []byte(`{{define "base"}}<!DOCTYPE html><html><body>{{template "nav" .}}<main>{{template "main" .}}</main></body></html>{{end}}`)},
		"templates/partials/nav.html":    {Data: []byte(`{{define "nav"}}<nav>Nav</nav>{{end}}`)},
		"templates/partials/footer.html": {Data: []byte(`<footer>Footer</footer>`)},
	}

	templates, err := TemplateCache(
		testFs,
		"templates/pages/*.html",
		"templates/partials/*.html",
		"templates/base.html",
	)
	if err != nil {
		return application{}, err
	}

	posts := make(map[string]models.Post)
	posts["home"] = models.Post{
		Metadata: models.PostMetadata{
			Title: "home",
			Slug:  "home",
			Tags:  []string{"test1", "test2"},
		},
		Content: template.HTML("<h1>markdown</h1>"),
	}
	app := &application{
		markdownCache: posts,
		templateCache: templates,
	}

	return *app, nil
}
