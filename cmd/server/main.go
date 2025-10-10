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
)

type config struct {
	port string
}

type application struct {
	templateCache map[string]*template.Template // use to search for template by name
	markdownCache map[string]models.Post        // used to search for markdown by slug name
	contentPath   string
	logger        *slog.Logger
	posts         []models.Post // used to keep track off all posts
}

func main() {
	port := flag.Int("port", 3999, "port to listen on")
	contentPath := flag.String("content", "./content/", "path on file sytem for content")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: false}))

	cfg := config{
		port: fmt.Sprintf(":%d", *port),
	}

	templateCache, err := TemplateCache()
	if err != nil {
		fmt.Println("Unable to create template cache, exiting to OS")
		os.Exit(1)
	}

	// fetch posts from content folder
	markdown, err := readMarkdownContent(os.DirFS(*contentPath))
	if err != nil {
		fmt.Println("Error reading markdown content: ", err)
	}

	// add all posts to the post list
	posts := []models.Post{}
	for _, p := range markdown {
		posts = append(posts, p)
	}

	app := application{
		logger:        logger,
		templateCache: templateCache,
		markdownCache: markdown,
		contentPath:   *contentPath,
		posts:         posts,
	}

	srv := http.Server{
		Addr:         cfg.port,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("Starting server on port%s\n", cfg.port)
	log.Fatal(srv.ListenAndServe())
}
