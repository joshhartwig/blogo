package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joshhartwig/blogo/internal/config"
	"github.com/joshhartwig/blogo/internal/logger"
	"github.com/joshhartwig/blogo/internal/server"
)

func main() {
	switch os.Args[1] {
	case "serve":
		runServe(os.Args[2:])
	default:
		log.Fatal("unknown command exiting")
	}
}

func runServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	configPath := fs.String("config", "", "path to config")
	addr := fs.String("addr", ":8080", "server address")
	fs.Parse(args)

	logger := logger.New()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("error loading config path: %v", err)
	}

	srv, err := server.New(cfg, logger)
	if err != nil {
		log.Fatalf("error creating new server: %v", err)
	}

	fmt.Printf("starting blogo server at %v\n", *addr)
	log.Fatal(srv.ListenAndServe(*addr))
}
