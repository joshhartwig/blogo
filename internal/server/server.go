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
	"github.com/joshhartwig/blogo/internal/models"
	"github.com/joshhartwig/blogo/internal/posts"
	"github.com/joshhartwig/blogo/internal/templates"
	"github.com/joshhartwig/blogo/ui"
)

type PostRepository interface {
	Count() int
	GetTopPosts(n int) []models.Post
	GetPostsBetweenRange(x, y int) []models.Post
	GetPostBySlug(slug string) (models.Post, error)
	SearchPosts(term string) []models.Post
	GetAllPostsInOrder() []models.Post
}

// PageData is the top-level data passed to every template. render() builds it
// automatically, so handlers only need to supply page-specific data via Page.
type PageData struct {
	Site  config.SiteConfig
	Theme config.ThemeConfig
	Page  any
}

type Server struct {
	cfg           config.SiteConfig
	themeCfg      config.ThemeConfig
	mux           *http.ServeMux
	templateCache map[string]*template.Template // use to search for template by name
	postRepo      PostRepository
	logger        *slog.Logger
}

func New(cfg config.SiteConfig, logger *slog.Logger) (*Server, error) {
	templateCache, err := newTemplateCache()
	if err != nil {
		return nil, fmt.Errorf("error creating new template cache: %w", err)
	}

	themeCfg, err := config.LoadTheme(ui.Files, cfg.Theme)
	if err != nil {
		return nil, fmt.Errorf("loading theme %q: %w", cfg.Theme, err)
	}

	postRepo, err := posts.NewPostRepository(os.DirFS(cfg.ContentDir))
	if err != nil {
		return &Server{}, err
	}

	s := &Server{
		cfg:           cfg,
		themeCfg:      themeCfg,
		mux:           http.NewServeMux(),
		templateCache: templateCache,
		postRepo:      postRepo,
		logger:        logger,
	}

	err = s.Routes()
	if err != nil {
		return s, err
	}

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
