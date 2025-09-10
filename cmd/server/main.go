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
	port string
}

type application struct {
	templateCache map[string]*template.Template
	logger        *slog.Logger
}

// TODO: template is rendering if the template name is home.html not home... why

func main() {
	port := flag.Int("port", 3999, "server port")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: false}))

	cfg := config{
		port: fmt.Sprintf(":%d", *port),
	}

	templateCache, err := TemplateCache()
	if err != nil {
		fmt.Println("unable to create template cache, exiting to OS")
		os.Exit(1)
	}

	app := application{
		logger:        logger,
		templateCache: templateCache,
	}

	srv := http.Server{
		Addr:    cfg.port,
		Handler: app.routes(),
	}

	fmt.Printf("starting server on port%s\n", cfg.port)
	err = srv.ListenAndServe()
	if err != nil {
		app.logger.Error("Fatal Error:", err.Error(), time.Now())
		fmt.Println(err)
		os.Exit(1)
	}
}
