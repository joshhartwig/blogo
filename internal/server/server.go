package server

import (
	"html/template"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joshhartwig/blogo/internal/config"
	"github.com/joshhartwig/blogo/internal/markdown"
	"github.com/joshhartwig/blogo/internal/models"
	"github.com/joshhartwig/blogo/internal/posts"
	"github.com/joshhartwig/blogo/internal/templates"
	"github.com/joshhartwig/blogo/ui"
)

type Server struct {
	cfg           config.SiteConfig
	mux           *http.ServeMux
	templateCache map[string]*template.Template // use to search for template by name
	markdownCache map[string]models.Post        // used to search for markdown by slug name
	postRepo      posts.PostRepository
	logger        *slog.Logger
}

func New(siteConfig config.SiteConfig) *Server {
	templateCache, err := newTemplateCache()
	if err != nil {
		log.Fatalf("error building template cache: %v", err)
	}

	markdownCache := markdown.NewCache(os.DirFS(siteConfig.ContentDir))

	s := &Server{
		cfg:           siteConfig,
		mux:           http.NewServeMux(),
		markdownCache: markdownCache,
		templateCache: templateCache,
	}

	s.Routes()

	return s
}

func newTemplateCache() (map[string]*template.Template, error) {
	templateFS, err := fs.Sub(ui.Files, "templates")
	if err != nil {
		return nil, err
	}
	return templates.NewCache(
		templateFS,
		"pages/*.html",
		"partials/*.html",
		"base.html",
	)

}
