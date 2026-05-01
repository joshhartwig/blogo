package server

import (
	"io/fs"
	"net/http"

	"github.com/joshhartwig/blogo/ui"
)

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	staticFS, _ := fs.Sub(ui.Files, "static")

	// file server for UI static assets
	mux.Handle("GET /static/",
		http.StripPrefix("/static",
			http.FileServer(http.FS(staticFS))))

	// file server for content assets (images in post directories)
	mux.Handle("GET /content/",
		http.StripPrefix("/content",
			http.FileServer(http.Dir(s.cfg.ContentDir))))

	// html routes - using security and common headers
	htmlHandler := func(name string, h http.HandlerFunc) http.Handler {
		return s.logHandler(name, s.setSecurity(s.setCommon(h)))
	}
	mux.Handle("/", htmlHandler("home", s.homeHandler))
	mux.Handle("/notfound", htmlHandler("notfound", s.notFoundHandler))
	mux.Handle("/search", htmlHandler("search", s.searchHandler))
	mux.Handle("/about", htmlHandler("about", s.aboutHandler))
	mux.Handle("/projects", htmlHandler("projects", s.projectsHandler))
	mux.Handle("/posts/", htmlHandler("listPost", s.listPostHandler))
	mux.Handle("/posts/{slug}", htmlHandler("showPost", s.showPostHandler))

	// api routes - set security
	apiHandler := func(h http.HandlerFunc) http.Handler {
		return s.setSecurity(h)
	}

	mux.Handle("/ping", apiHandler(ping))
	mux.Handle("/rss", apiHandler(s.rssHandler))

	return mux
}
