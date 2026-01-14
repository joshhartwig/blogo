package main

import (
	"net/http"
	"path/filepath"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// file server
	mux.Handle("GET /static/",
		http.StripPrefix("/static",
			//http.FileServer(http.Dir(fmt.Sprintf("./themes/%s/static", app.cfg.theme)))))
			http.FileServer(http.Dir(filepath.Join(app.cfg.themesPath, app.cfg.theme)))))

	// html routes - using security and common headers
	htmlHandler := func(name string, h http.HandlerFunc) http.Handler {
		return app.logHandler(name, app.setSecurity(app.setCommon(h)))
	}
	mux.Handle("/", htmlHandler("home", app.homeHandler))
	mux.Handle("/notfound", htmlHandler("notfound", app.notFoundHandler))
	mux.Handle("/search", htmlHandler("search", app.searchHandler))
	mux.Handle("/about", htmlHandler("about", app.aboutHandler))
	mux.Handle("/projects", htmlHandler("projects", app.projectsHandler))
	mux.Handle("/posts/", htmlHandler("listPost", app.listPostHandler))
	mux.Handle("/posts/{slug}", htmlHandler("showPost", app.showPostHandler))

	// api routes - set security
	apiHandler := func(h http.HandlerFunc) http.Handler {
		return app.setSecurity(h)
	}

	mux.Handle("/ping", apiHandler(ping))
	mux.Handle("/rss", apiHandler(app.rssHandler))

	return app.logRequests(mux)
}
