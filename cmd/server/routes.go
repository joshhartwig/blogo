package main

import "net/http"

func (c *config) routes() http.Handler {
	router := http.NewServeMux()
	router.Handle("/static", c.middleware(http.FileServer(http.Dir("/ui/static"))))
	router.Handle("/", c.middleware(http.HandlerFunc(c.homeHandler)))
	router.HandleFunc("/blog", c.blogHandler)

	return router
}
