package main

import "net/http"

func (app *application) routes() http.Handler {
	router := http.NewServeMux()

	// file server
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	router.Handle("/ping", http.HandlerFunc(ping))
	router.Handle("/rss", http.HandlerFunc(app.rssHandler))

	router.Handle("/", http.HandlerFunc(app.homeHandler))
	router.Handle("/notfound", http.HandlerFunc(app.notFoundHandler))
	router.Handle("/search", http.HandlerFunc(app.searchHandler))
	router.Handle("/about", http.HandlerFunc(app.aboutHandler))
	router.Handle("/projects", http.HandlerFunc(app.projectsHandler))
	router.Handle("/posts", http.HandlerFunc(app.listPostHandler))
	router.Handle("/posts/{slug}", http.HandlerFunc(app.showPostHandler))

	return app.setSecurity(app.setCommon(app.logRequests(router)))
}
