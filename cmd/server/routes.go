package main

import "net/http"

func (app *application) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/ping", app.ping)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	router.Handle("/posts/{slug}", http.HandlerFunc(app.postHandler))
	router.Handle("/", app.middleware(http.HandlerFunc(app.homeHandler)))
	router.Handle("/notfound", http.HandlerFunc(app.notFoundHandler))
	router.Handle("/about", http.HandlerFunc(app.aboutHandler))

	return router
}
