package main

import "net/http"

func (c *config) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/ping", c.ping)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	router.Handle("/", c.middleware(http.HandlerFunc(c.homeHandler)))
	router.HandleFunc("/blog", c.blogHandler)

	return router
}
