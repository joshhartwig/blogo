package main

import (
	"io/fs"
	"net/http"

	"github.com/joshhartwig/blogo/ui"
	_ "github.com/joshhartwig/blogo/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	staticFS, _ := fs.Sub(ui.Files, "static")

	// file server for UI static assets
	mux.Handle("GET /static/",
		http.StripPrefix("/static",
			http.FileServer(http.FS(staticFS))))

	// file server for content assets (images in post directories)
	mux.Handle("GET /content/",
		http.StripPrefix("/content",
			http.FileServer(http.Dir(app.contentPath))))

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
