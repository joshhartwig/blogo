package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"

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
}

type application struct {
	port          string
	templateCache map[string]*template.Template // use to search for template by name
	markdownCache map[string]models.Post        // used to search for markdown by slug name
	contentPath   string
	logger        *slog.Logger
	posts         []models.Post // used to keep track off all posts
	postsPerPage  int
	postRepo      posts.PostRepository
	cfg           config
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 3999, "Server Port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|production)")
	flag.StringVar(&cfg.contentPath, "content", "./content/", "path on file sytem for content")

	postsPerPage := flag.Int("posts per page", 5, "sets the default count of posts per page")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: false}))

	templateCache, err := TemplateCache(os.DirFS(
		"./ui/templates/"),
		"pages/*.html",
		"partials/*.html",
		"base.html",
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
	allPosts := []models.Post{}
	for _, p := range markdown {
		allPosts = append(allPosts, p)
	}

	app := application{
		templateCache: templateCache,
		markdownCache: markdown,
		contentPath:   cfg.contentPath,
		posts:         allPosts,
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

	fmt.Printf("Starting server on port%s\n", app.port)
	log.Fatal(srv.ListenAndServe())
}
