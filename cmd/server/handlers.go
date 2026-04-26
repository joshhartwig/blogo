package main

import (
	"encoding/xml"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joshhartwig/blogo/internal/models"
)

var (
	slugRegex = regexp.MustCompile(`^[a-z0-9-]+$`)
)

// ping is used for testing endpoints to ensure handlers are working
func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

// notFoundHandler is intended for posts that are not found, it will render the notfound template
func (app *application) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusNotFound, "notfound", nil)

}

// homeHandler is our default handler for the '/' route
func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	var pagePosts = []models.Post{}

	// if the post count is < 5 do the max else do 5
	if len(app.postRepo.Posts) < 5 {
		pagePosts = app.postRepo.GetTopPosts(min(len(app.postRepo.Posts)))
	} else {
		pagePosts = app.postRepo.GetTopPosts(5)
	}

	data := models.HomePageData{
		Posts: pagePosts,
	}

	app.render(w, http.StatusOK, "home", data)
}

// about handler renders the about page
func (app *application) aboutHandler(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "about", nil)
}

// route: /posts lists all posts with paging
func (app *application) listPostHandler(w http.ResponseWriter, r *http.Request) {
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if n, err := strconv.Atoi(p); err == nil && n > 0 {
			page = n
		}
	}

	pagination := models.NewPagination(page, app.postsPerPage, len(app.postRepo.Posts))
	start := pagination.PostsStart
	end := pagination.PostsEnd
	//fmt.Printf("start %d end %d %v", start, end, pagination)
	data := models.HomePageData{
		Posts:      app.postRepo.GetPostsBetweenRange(start, end),
		Pagination: pagination,
	}

	app.render(w, http.StatusOK, "posts", data)
}

// postHandler renders our post content assuming the slug is found
func (app *application) showPostHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimSpace(r.PathValue("slug")) // fetch the slug

	if slug == "" || len(slug) > 200 {
		app.logger.Warn("path not found", "info", slug, "handler", "post")
		http.Redirect(w, r, "/notfound", http.StatusSeeOther)
		return
	}

	if !isValidSlug(slug) {
		app.logger.Warn("invalid slug format", "slug", slug)
		http.Redirect(w, r, "/notfound", http.StatusSeeOther)
		return
	}

	post, ok := app.markdownCache[slug]
	if !ok {
		http.Redirect(w, r, "/notfound", http.StatusSeeOther)
		return
	}

	app.render(w, http.StatusOK, "post", post)
}

func (app *application) projectsHandler(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "projects", nil)
}

// rssHandler generates and writes an RSS XML feed built from the application's markdown cache.
// It iterates over app.markdownCache, converts each post's metadata into a models.Item
// (Title from post.Metadata.Title, Link from post.Metadata.Slug, Description from post.Metadata.Summary,
// Category set to "blog", GUID generated via uuid.New()), and collects those items into a models.RSS
// feed with feed-level metadata (Title "Josh's Blog", Link "https://localhost:3999", Description,
// Language "English", PubDate set to time.Now(), Category "blog").
// The final feed is marshaled to XML and written to the response with Content-Type "application/xml".
// If XML marshaling fails the error is logged using app.logger.Error and the handler returns
// without writing a response body or setting an explicit HTTP status code.
//
// Notes:
//   - The request parameter r is not used by this handler beyond satisfying the http.Handler signature.
//   - The handler does not write an explicit HTTP status code on success, does not include an XML
//     declaration header, and assumes app.markdownCache can be safely read without additional locking.
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/xml")
	w.Write([]byte(xml.Header)) // fix missing header
	w.Write(data)
}

// searchHandler handles HTTP requests that perform a search for posts.
// It extracts the "q" query parameter from the request URL, invokes the
// application's post repository to search for matching posts, and renders
// the "search" template with a data object containing the found posts and
// the original search term. Rendering errors are logged via the application's
// logger. The handler writes the rendered output to the provided http.ResponseWriter.
func (app *application) searchHandler(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("q")

	if term == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if len(term) > 200 {
		app.logger.Error("seach term too long", "length", len(term))
		http.Error(w, "Search term too long", http.StatusBadRequest)
		return
	}

	results := app.postRepo.SearchPosts(term)

	data := struct {
		Posts []models.Post
		Term  string
	}{
		Posts: results,
		Term:  term,
	}

	app.render(w, http.StatusOK, "search", data)
}

// isValidSlug reports whether slug is a non-empty, ASCII-only slug suitable for URLs
func isValidSlug(slug string) bool {
	return slugRegex.MatchString(slug)
}
