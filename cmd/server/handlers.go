package main

import (
	"fmt"
	"net/http"
)

func (app *application) ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {

	err := app.render(w, "home", nil)

	if err != nil {
		app.logger.Error("error rendering template", err.Error(), "error")
		http.Error(w, "error rendering template", http.StatusInternalServerError)
	}
}

// middleware outputs the method and uri the request is hitting on the server
func (app *application) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("<= %s | %s \n", r.Method, r.URL)
		w.Header().Set("Content-Type", "text/html")
		next.ServeHTTP(w, r)
	})
}

func (app *application) postHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	if slug == "" {
		http.Redirect(w, r, "/", http.StatusNoContent)
		return
	}

	data, ok := app.markdownCache[slug]
	if !ok {
		http.Redirect(w, r, "/", http.StatusNoContent)
		return
	}

	app.render(w, "posts", data)
}
