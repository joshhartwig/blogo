package main

import (
	"fmt"
	"net/http"

	"github.com/joshhartwig/blogo/internal/models"
)

func (app *application) ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

// notFoundHandler is intended for posts that are not foud, it will render the notfound template
func (app *application) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	if err := app.render(w, "notfound", nil); err != nil {
		app.logger.Error(err.Error(), "error", "error rendering")
		http.NotFound(w, r)
	}

}

// homeHandler is our default handler for the '/' route
func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	td := models.TemplateData{}
	for _, p := range app.markdownCache {
		td.Posts = append(td.Posts, p)
	}

	if err := app.render(w, "home", td); err != nil {
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

// postHandler renders our post content assuming the slug is found
func (app *application) postHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug") // fetch the slug
	fmt.Printf("in postHandler found slug %s\n", slug)
	if slug == "" {
		http.Redirect(w, r, "/", http.StatusNoContent) // if slug contains no content redirect to home
		return
	}

	post, ok := app.markdownCache[slug]
	if !ok {
		fmt.Println("post not found, redirecting")
		http.Redirect(w, r, "/notfound", http.StatusPermanentRedirect)
		return
	}

	app.render(w, "posts", post)
}
