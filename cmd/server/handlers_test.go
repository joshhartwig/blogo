package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
	fmt.Println(app.templateCache)
	if err != nil {
		t.Errorf("error fetching mocked app")
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	app.homeHandler(rr, req)
	got := rr.Body.String()
	if !strings.Contains(got, "Josh's") {
		t.Errorf("got %s", got)
	}
}

func returnMockedApp() (application, error) {

	templates, err := TemplateCache()
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
		Content: template.HTML("<h1>hi</h1>"),
	}
	app := &application{
		logger:        slog.New(slog.DiscardHandler),
		markdownCache: posts,
		templateCache: templates,
	}

	return *app, nil
}
