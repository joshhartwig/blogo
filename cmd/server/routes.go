package main

import "net/http"

func (app *application) routes() http.Handler {
	router := http.NewServeMux()

	// file server
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	router.Handle("/", app.middleware(http.HandlerFunc(app.homeHandler)))

	router.Handle("/ping", app.middleware(http.HandlerFunc(app.ping)))
	router.Handle("/notfound", app.middleware(http.HandlerFunc(app.notFoundHandler)))
	router.Handle("/about", app.middleware(http.HandlerFunc(app.aboutHandler)))
	router.Handle("/projects", app.middleware(http.HandlerFunc(app.projectsHandler)))
	router.Handle("/posts/{slug}", app.middleware(http.HandlerFunc(app.postHandler)))

	return router
}
