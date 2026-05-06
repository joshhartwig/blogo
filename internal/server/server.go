package server

import (
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joshhartwig/blogo/internal/config"
	"github.com/joshhartwig/blogo/internal/posts"
	"github.com/joshhartwig/blogo/internal/templates"
	"github.com/joshhartwig/blogo/ui"
)

type Server struct {
	cfg           config.SiteConfig
	mux           *http.ServeMux
	templateCache map[string]*template.Template // use to search for template by name
	postRepo      posts.PostRepository
	logger        *slog.Logger
}

func New(cfg config.SiteConfig, logger *slog.Logger) (*Server, error) {
	templateCache, err := newTemplateCache()
	if err != nil {
		return nil, fmt.Errorf("error creating new template cache: %w", err)
	}

	s := &Server{
		cfg:           cfg,
		mux:           http.NewServeMux(),
		templateCache: templateCache,
		postRepo:      posts.NewPostRepository(os.DirFS(cfg.ContentDir)),
		logger:        logger,
	}

	s.Routes()
	return s, nil
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

func (s *Server) ListenAndServe(addr string) error {
	server := &http.Server{
		Addr:              addr,
		Handler:           s.mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return server.ListenAndServe()
}
