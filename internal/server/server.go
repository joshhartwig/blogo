package server

import (
	"log/slog"
	"net/http"
	"os"
	"text/template"

	"github.com/joshhartwig/blogo/internal/config"
	"github.com/joshhartwig/blogo/internal/markdown"
	"github.com/joshhartwig/blogo/internal/models"
	"github.com/joshhartwig/blogo/internal/posts"
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
	markdownCache := markdown.NewCache(os.DirFS(siteConfig.ContentDir))

	s := &Server{
		cfg:           siteConfig,
		mux:           http.NewServeMux(),
		markdownCache: markdownCache,
	}

	s.Routes()

	return s
}
