package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"path/filepath"

	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joshhartwig/blogo/internal/models"
	"github.com/joshhartwig/blogo/internal/posts"
)

type config struct {
	port        int
	env         string
	contentPath string
	theme       string
	themesPath  string
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
	flag.StringVar(&cfg.themesPath, "theme path", "themes", "Path to the directory containing all theme folders (e.g., themes)")
	flag.StringVar(&cfg.theme, "theme", "default", "Name of the theme folder to use (e.g., default, darkblue). Should match a folder under the themes path.")

	postsPerPage := flag.Int("PostsPerPage", 5, "sets the default count of posts per page")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: false}))

	fmt.Printf("\nStarting Server...\n")

	themeDir := os.DirFS(filepath.Join(cfg.themesPath, cfg.theme))
	fmt.Printf("Blog Theme: %s\n", themeDir)

	templateCache, err := LoadTemplatesAsMap(
		themeDir,
		"templates/pages/*.html",
		"templates/partials/*.html",
		"templates/base.html",
	)

	if err != nil {
		fmt.Println("Unable to create template cache, exiting to OS", err)
		os.Exit(1)
	}

	// fetch posts from content folder
	markdown, err := readMarkdownContent(os.DirFS(cfg.contentPath))
	if err != nil {
		fmt.Println("Error reading markdown content: ", err)
	}

	// add all posts to the post list
	fmt.Printf("Found Following MarkDown Posts\n")
	allPosts := []models.Post{}
	for _, p := range markdown {
		fmt.Printf("Title: %s Slug: %s\n", p.Metadata.Title, p.Metadata.Slug)
		allPosts = append(allPosts, p)
	}
	fmt.Printf("\ntotal post count: %d\n", len(allPosts))

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

	fmt.Printf("\nWeb Server running on port :%d\n", cfg.port)
	log.Fatal(srv.ListenAndServe())
}
