package main

import (
	"flag"
	"fmt"
	"html/template"

	"log/slog"
	"net/http"
	"os"
	"time"
)

type config struct {
	port     string
	logger   *slog.Logger
	tplCache map[string]*template.Template
}

// TODO: template is rendering if the template name is home.html not home... why

func main() {
	port := flag.Int("port", 3999, "server port")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: false}))

	cfg := config{
		port:   fmt.Sprintf(":%d", *port),
		logger: logger,
	}

	srv := http.Server{
		Addr:    cfg.port,
		Handler: cfg.routes(),
	}

	templateCache, err := TemplateCache()
	if err != nil {
		cfg.logger.Error("error", err.Error(), "error parsing template cache")
		os.Exit(1)
	}

	cfg.tplCache = templateCache

	fmt.Printf("starting server on port%s\n", cfg.port)
	err = srv.ListenAndServe()
	if err != nil {
		cfg.logger.Error("Fatal Error:", err.Error(), time.Now())
		fmt.Println(err)
		os.Exit(1)
	}
}
