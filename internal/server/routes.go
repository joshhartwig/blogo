package server

import (
	"io/fs"
	"net/http"

	"github.com/joshhartwig/blogo/ui"
)

func (s *Server) Routes() error {

	staticFS, err := fs.Sub(ui.Files, "static")
	if err != nil {
		return err
	}

	// file server for UI static assets
	s.mux.Handle("GET /static/",
		http.StripPrefix("/static",
			http.FileServer(http.FS(staticFS))))

	// file server for content assets (images in post directories)
	s.mux.Handle("GET /content/",
		http.StripPrefix("/content",
			http.FileServer(http.Dir(s.cfg.ContentDir))))

	// html routes - using security and common headers
	htmlHandler := func(name string, h http.HandlerFunc) http.Handler {
		return s.logHandler(name, s.setSecurity(s.setCommon(h)))
	}
	s.mux.Handle("/", htmlHandler("home", s.homeHandler))
	s.mux.Handle("/notfound", htmlHandler("notfound", s.notFoundHandler))
	s.mux.Handle("/search", htmlHandler("search", s.searchHandler))
	s.mux.Handle("/about", htmlHandler("about", s.aboutHandler))
	s.mux.Handle("/projects", htmlHandler("projects", s.projectsHandler))
	s.mux.Handle("/posts/", htmlHandler("listPost", s.listPostHandler))
	s.mux.Handle("/posts/{slug}", htmlHandler("showPost", s.showPostHandler))

	// api routes - set security
	apiHandler := func(h http.HandlerFunc) http.Handler {
		return s.setSecurity(h)
	}

	s.mux.Handle("/ping", apiHandler(ping))
	s.mux.Handle("/rss", apiHandler(s.rssHandler))

	return nil
}
