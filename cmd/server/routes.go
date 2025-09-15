package main

import "net/http"

func (app *application) routes() http.Handler {
	router := http.NewServeMux()

	// file server
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	router.Handle("/", app.middleware(http.HandlerFunc(app.homeHandler)))

	router.Handle("/posts/{slug}", http.HandlerFunc(app.postHandler))
	router.HandleFunc("/ping", http.HandlerFunc(app.ping))
	router.Handle("/notfound", app.middleware(http.HandlerFunc(app.notFoundHandler)))
	router.Handle("/about", app.middleware(http.HandlerFunc(app.aboutHandler)))

	return router
}
