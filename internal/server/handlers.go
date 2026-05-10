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
func (s *Server) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	s.render(w, http.StatusNotFound, "notfound", nil)
}

// homeHandler is our default handler for the '/' route
func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	s.render(w, http.StatusOK, "home", models.PostListData{Posts: s.postRepo.GetTopPosts(5)})
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
	s.render(w, http.StatusOK, "posts", models.PostListData{
		Posts:      s.postRepo.GetPostsBetweenRange(start, end),
		Pagination: pagination,
	})
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

// rssHandler writes an RSS 2.0 XML feed of all posts.
func (s *Server) rssHandler(w http.ResponseWriter, r *http.Request) {
	items := []models.Item{}
	feedData := models.RSS{
		Title:       s.cfg.SiteTitle,
		Link:        s.cfg.BaseURL,
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

	s.render(w, http.StatusOK, "search", models.SearchData{Posts: results, Term: term})
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
