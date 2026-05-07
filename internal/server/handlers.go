package server

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joshhartwig/blogo/internal/config"
	"github.com/joshhartwig/blogo/internal/models"
)

var (
	slugRegex = regexp.MustCompile(`^[a-z0-9-]+$`)
)

type PageData struct {
	Site config.SiteConfig
	Page any
}

// ping is used for testing endpoints to ensure handlers are working
func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

// notFoundHandler is intended for posts that are not found, it will render the notfound template
func (s *Server) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	s.render(w, http.StatusNotFound, "notfound", nil)
}

// homeHandler is our default handler for the '/' route
func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	var pagePosts = []models.Post{}

	// if the post count is < 5 do the max else do 5
	var postCount = s.postRepo.Count()
	if postCount < 5 {
		pagePosts = s.postRepo.GetTopPosts(min(postCount))
	} else {
		pagePosts = s.postRepo.GetTopPosts(5)
	}

	data := models.HomePageData{
		Posts: pagePosts,
	}

	s.render(w, http.StatusOK, "home", data)
}

// about handler renders the about page
func (s *Server) aboutHandler(w http.ResponseWriter, r *http.Request) {
	s.render(w, http.StatusOK, "about", nil)
}

// route: /posts lists all posts with paging
func (s *Server) listPostHandler(w http.ResponseWriter, r *http.Request) {
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if n, err := strconv.Atoi(p); err == nil && n > 0 {
			page = n
		}
	}

	pagination := models.NewPagination(page, s.cfg.PostsPerPage, s.postRepo.Count())
	start := pagination.PostsStart
	end := pagination.PostsEnd
	//fmt.Printf("start %d end %d %v", start, end, pagination)
	data := models.HomePageData{
		Posts:      s.postRepo.GetPostsBetweenRange(start, end),
		Pagination: pagination,
	}

	s.render(w, http.StatusOK, "posts", data)
}

// postHandler renders our post content assuming the slug is found
func (s *Server) showPostHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimSpace(r.PathValue("slug")) // fetch the slug

	if slug == "" || len(slug) > 200 {
		s.logger.Warn("path not found", "info", slug, "handler", "post")
		http.Redirect(w, r, "/notfound", http.StatusSeeOther)
		return
	}

	if !isValidSlug(slug) {
		s.logger.Warn("invalid slug format", "slug", slug)
		http.Redirect(w, r, "/notfound", http.StatusSeeOther)
		return
	}

	post, err := s.postRepo.GetPostBySlug(slug)
	if err != nil {
		http.Redirect(w, r, "/notfound", http.StatusSeeOther)
		return
	}

	s.render(w, http.StatusOK, "post", post)
}

func (s *Server) projectsHandler(w http.ResponseWriter, r *http.Request) {
	s.render(w, http.StatusOK, "projects", nil)
}

// rssHandler generates and writes an RSS XML feed built from the slication's markdown cache.
// It iterates over s.markdownCache, converts each post's metadata into a models.Item
// (Title from post.Metadata.Title, Link from post.Metadata.Slug, Description from post.Metadata.Summary,
// Category set to "blog", GUID generated via uuid.New()), and collects those items into a models.RSS
// feed with feed-level metadata (Title "Josh's Blog", Link "https://localhost:3999", Description,
// Language "English", PubDate set to time.Now(), Category "blog").
// The final feed is marshaled to XML and written to the response with Content-Type "slication/xml".
// If XML marshaling fails the error is logged using s.logger.Error and the handler returns
// without writing a response body or setting an explicit HTTP status code.
//
// Notes:
//   - The request parameter r is not used by this handler beyond satisfying the http.Handler signature.
//   - The handler does not write an explicit HTTP status code on success, does not include an XML
//     declaration header, and assumes s.markdownCache can be safely read without additional locking.
func (s *Server) rssHandler(w http.ResponseWriter, r *http.Request) {
	items := []models.Item{}
	feedData := models.RSS{
		Title:       s.cfg.SiteTitle,
		Link:        "https://localhost:3999",
		Description: s.cfg.Description,
		Language:    "English",
		PubDate:     time.Now(),
		Category:    "blog",
		Item:        items,
	}

	for _, post := range s.postRepo.GetAllPostsInOrder() {
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
		s.logger.Error(err.Error(), "error marshaling feed data to xml", "rssHanlder")
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
func (s *Server) searchHandler(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("q")

	if term == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if len(term) > 200 {
		s.logger.Error("seach term too long", "length", len(term))
		http.Error(w, "Search term too long", http.StatusBadRequest)
		return
	}

	results := s.postRepo.SearchPosts(term)

	data := struct {
		Posts []models.Post
		Term  string
	}{
		Posts: results,
		Term:  term,
	}

	s.render(w, http.StatusOK, "search", data)
}

// isValidSlug reports whether slug is a non-empty, ASCII-only slug suitable for URLs
func isValidSlug(slug string) bool {
	return slugRegex.MatchString(slug)
}

// render executes a template with the passed in data
func (s *Server) render(w http.ResponseWriter, status int, templateName string, data any) {
	v, ok := s.templateCache[templateName]
	if !ok {
		err := fmt.Errorf("template not found: %s", templateName)
		s.logger.Error(err.Error())
		return
	}

	pageData := PageData{
		Site: s.cfg,
		Page: data,
	}

	buf := new(bytes.Buffer)
	err := v.ExecuteTemplate(buf, "base", pageData)
	if err != nil {
		s.logger.Error("template execution failed", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	// write contents of buffer to response writer
	buf.WriteTo(w)
}
