package server

import (
	"fmt"
	"html/template"
	"io/fs"
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

func New(cfg config.SiteConfig, logger *slog.Logger) (*Server, error) {
	templateCache, err := newTemplateCache()
	if err != nil {
		return nil, fmt.Errorf("error creating new template cache: %v", err)
	}

	markdownCache := markdown.NewCache(os.DirFS(cfg.ContentDir))

	// load all posts
	var blogPosts []models.Post
	for _, p := range markdownCache {
		blogPosts = append(blogPosts, p)
	}

	s := &Server{
		cfg:           cfg,
		mux:           http.NewServeMux(),
		markdownCache: markdownCache,
		templateCache: templateCache,
		postRepo:      posts.NewPostRepository(blogPosts),
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
	return http.ListenAndServe(addr, s.mux)
}
