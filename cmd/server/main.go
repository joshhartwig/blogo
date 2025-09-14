package main

import (
	"flag"
	"fmt"
	"html/template"

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
	templateCache map[string]*template.Template
	markdownCache map[string]models.Post
	logger        *slog.Logger
}

func main() {
	port := flag.Int("port", 3999, "server port")
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
	posts, err := readMarkdownContent()
	if err != nil {
		fmt.Println("Error reading markdown content: ", err)
	}

	app := application{
		logger:        logger,
		templateCache: templateCache,
		markdownCache: posts,
	}

	srv := http.Server{
		Addr:    cfg.port,
		Handler: app.routes(),
	}

	fmt.Printf("Starting server on port%s\n", cfg.port)
	err = srv.ListenAndServe()
	if err != nil {
		app.logger.Error("Fatal Error:", err.Error(), time.Now())
		fmt.Println(err)
		os.Exit(1)
	}
}
