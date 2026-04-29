package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joshhartwig/blogo/internal/logger"
	"github.com/joshhartwig/blogo/internal/models"
	"github.com/joshhartwig/blogo/internal/posts"
	"github.com/joshhartwig/blogo/ui"
)

type config struct {
	port        int
	env         string
	contentPath string
}

type application struct {
	templateCache map[string]*template.Template // use to search for template by name
	markdownCache map[string]models.Post        // used to search for markdown by slug name
	contentPath   string
	logger        *slog.Logger

	postsPerPage int
	postRepo     posts.PostRepository
	cfg          config
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 3999, "Server Port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|production)")
	flag.StringVar(&cfg.contentPath, "content", "./content/", "Path to the directory containing markdown content files (e.g., ./content/)")

	postsPerPage := flag.Int("PostsPerPage", 5, "sets the default count of posts per page")

	flag.Parse()
	logger := logger.New()

	logger.Info("starting server")

	templateFS, err := fs.Sub(ui.Files, "templates")
	if err != nil {
		logger.Error("failed to get embedded templates", "error", err)
		os.Exit(1)
	}

	templateCache, err := LoadTemplatesAsMap(
		templateFS,
		"pages/*.html",
		"partials/*.html",
		"base.html",
	)

	if err != nil {
		logger.Error("unable to create template cache", "error", err)
		os.Exit(1)
	}

	// fetch posts from content folder
	markdown, err := readMarkdownContent(os.DirFS(cfg.contentPath))
	if err != nil {
		logger.Error("error reading markdown content", "error", err)
	}

	// add all posts to the post list
	logger.Info("loading markdown posts")
	allPosts := []models.Post{}
	for _, p := range markdown {
		logger.Info("loaded post", "title", p.Metadata.Title, "slug", p.Metadata.Slug)
		allPosts = append(allPosts, p)
	}
	logger.Info("post loading complete", "count", len(allPosts))

	app := application{
		templateCache: templateCache,
		markdownCache: markdown,
		contentPath:   cfg.contentPath,
		postsPerPage:  *postsPerPage,
		postRepo:      posts.NewPostRepository(allPosts),
		cfg:           cfg,
		logger:        logger,
	}

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("web server started", "port", cfg.port)
	logger.Error("server error", "error", srv.ListenAndServe())
	os.Exit(1)
}
