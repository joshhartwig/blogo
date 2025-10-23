package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/joshhartwig/blogo/internal/models"
)

// ping is used for testing endpoints to ensure handlers are working
func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
	w.WriteHeader(http.StatusOK)
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
	var pagePosts = []models.Post{}

	// if the post count is < 5 do the max else do 5
	if len(app.postRepo.Posts) < 5 {
		pagePosts = app.postRepo.GetTopPosts(len(app.postRepo.Posts) - 1)
	} else {
		pagePosts = app.postRepo.GetTopPosts(5)
	}

	data := models.HomePageData{
		Posts: pagePosts,
	}

	if err := app.render(w, "home", data); err != nil {
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

// about handler renders the about page
func (app *application) aboutHandler(w http.ResponseWriter, r *http.Request) {
	if err := app.render(w, "about", nil); err != nil {
		app.logger.Error("render failed", "error", err, "handler", "about")
		return
	}
}

// route: /posts lists all posts with paging
func (app *application) listPostHandler(w http.ResponseWriter, r *http.Request) {
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if n, err := strconv.Atoi(p); err == nil && n > 0 {
			page = n
		}
	}
	fmt.Println("page:", page)
	fmt.Printf("param:%d postsPerPage:%d posts:%d", page, app.postsPerPage, len(app.postRepo.Posts))
	pagination := models.NewPagination(page, app.postsPerPage, len(app.postRepo.Posts))
	start := pagination.PostsStart
	end := pagination.PostsEnd
	fmt.Println("pagination:", pagination)
	data := models.HomePageData{
		Posts:      app.postRepo.GetPostsBetweenRange(start, end),
		Pagination: pagination,
	}

	if err := app.render(w, "posts", data); err != nil {
		app.logger.Error("error rendering template", err.Error(), "error")
		http.Error(w, "error rendering template", http.StatusInternalServerError)
	}
}

// postHandler renders our post content assuming the slug is found
func (app *application) showPostHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug") // fetch the slug

	if slug == "" {
		app.logger.Info("path not found", "info", slug, "handler", "post")
		http.Redirect(w, r, "/", http.StatusNoContent) // if slug contains no content redirect to home
		return
	}

	post, ok := app.markdownCache[slug]
	if !ok {
		fmt.Println("post not found, redirecting")
		http.Redirect(w, r, "/notfound", http.StatusSeeOther)
		return
	}

	app.render(w, "post", post)
}

func (app *application) projectsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (app *application) rssHandler(w http.ResponseWriter, r *http.Request) {
	items := []models.Item{}
	feedData := models.RSS{
		Title:       "Josh's Blog",
		Link:        "https://localhost:3999",
		Description: "description",
		Language:    "English",
		PubDate:     time.Now(),
		Category:    "blog",
		Item:        items,
	}

	for _, post := range app.markdownCache {
		item := models.Item{
			Title:       post.Metadata.Title,
			Link:        post.Metadata.Slug,
			Description: post.Metadata.Summary,
			Category:    "blog",
			GUID:        uuid.New().String(),
		}

		feedData.Item = append(feedData.Item, item)
	}

	data, err := xml.Marshal(feedData)
	if err != nil {
		app.logger.Error(err.Error(), "error marshaling feed data to xml", "rssHanlder")
		return
	}

	w.Header().Add("Content-Type", "application/xml")
	w.Write(data)
}

func (app *application) searchHandler(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("q")
	results := app.searchPosts(term)

	data := struct {
		Posts []models.Post
		Term  string
	}{
		Posts: results,
		Term:  term,
	}

	err := app.render(w, "search", data)
	if err != nil {
		app.logger.Error(err.Error(), "error fetching posts", "error")
		return
	}
}
