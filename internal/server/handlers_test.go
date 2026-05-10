package server

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/joshhartwig/blogo/internal/config"
	"github.com/joshhartwig/blogo/internal/models"
	"github.com/joshhartwig/blogo/internal/posts"
	tu "github.com/joshhartwig/blogo/internal/testutils"
)

// fakeRepo satisfies PostRepository with an in-memory slice of posts.
type fakeRepo struct {
	posts []models.Post
}

func (f *fakeRepo) Count() int { return len(f.posts) }

func (f *fakeRepo) GetTopPosts(n int) []models.Post {
	if n > len(f.posts) {
		n = len(f.posts)
	}
	return f.posts[:n]
}

func (f *fakeRepo) GetPostsBetweenRange(x, y int) []models.Post {
	if x < 0 {
		x = 0
	}
	if y > len(f.posts) {
		y = len(f.posts)
	}
	if x > y {
		return []models.Post{}
	}
	return f.posts[x:y]
}

func (f *fakeRepo) GetPostBySlug(slug string) (models.Post, error) {
	for _, p := range f.posts {
		if p.Metadata.Slug == slug {
			return p, nil
		}
	}
	return models.Post{}, posts.ErrPostNotFound
}

func (f *fakeRepo) SearchPosts(term string) []models.Post {
	results := []models.Post{}
	for _, p := range f.posts {
		if strings.Contains(strings.ToLower(p.Metadata.Title), strings.ToLower(term)) {
			results = append(results, p)
		}
	}
	return results
}

func (f *fakeRepo) GetAllPostsInOrder() []models.Post {
	return f.posts
}

func newTestServer(t *testing.T, repo PostRepository) *Server {
	t.Helper()
	tc, err := newTemplateCache()
	if err != nil {
		t.Fatalf("template cache: %v", err)
	}
	return &Server{
		cfg: config.SiteConfig{
			PostsPerPage: 5,
			SiteTitle:    "Test Blog",
		},
		mux:           http.NewServeMux(),
		templateCache: tc,
		postRepo:      repo,
		logger:        slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

func TestPing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	ping(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("got %d, want 200", w.Code)
	}
	if w.Body.String() != "pong" {
		t.Errorf("got %q, want %q", w.Body.String(), "pong")
	}
}

func TestNotFoundHandler(t *testing.T) {
	s := newTestServer(t, &fakeRepo{})
	req := httptest.NewRequest(http.MethodGet, "/notfound", nil)
	w := httptest.NewRecorder()
	s.notFoundHandler(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("got %d, want 404", w.Code)
	}
}

func TestHomeHandler_WithPosts(t *testing.T) {
	posts := []models.Post{
		tu.NewPost(tu.WithTitle("First Post"), tu.WithSlug("first-post")),
		tu.NewPost(tu.WithTitle("Second Post"), tu.WithSlug("second-post")),
	}
	s := newTestServer(t, &fakeRepo{posts: posts})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	s.homeHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("got %d, want 200", w.Code)
	}
	if !strings.Contains(w.Body.String(), "First Post") {
		t.Error("body missing expected post title")
	}
}

func TestHomeHandler_NoPosts(t *testing.T) {
	s := newTestServer(t, &fakeRepo{})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	s.homeHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("got %d, want 200", w.Code)
	}
}

func TestAboutHandler(t *testing.T) {
	s := newTestServer(t, &fakeRepo{})
	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	w := httptest.NewRecorder()
	s.aboutHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("got %d, want 200", w.Code)
	}
}

func TestProjectsHandler(t *testing.T) {
	s := newTestServer(t, &fakeRepo{})
	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	w := httptest.NewRecorder()
	s.projectsHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("got %d, want 200", w.Code)
	}
}

func TestListPostHandler(t *testing.T) {
	posts := []models.Post{
		tu.NewPost(tu.WithTitle("Post 1"), tu.WithSlug("post-1")),
		tu.NewPost(tu.WithTitle("Post 2"), tu.WithSlug("post-2")),
	}
	s := newTestServer(t, &fakeRepo{posts: posts})
	req := httptest.NewRequest(http.MethodGet, "/posts/", nil)
	w := httptest.NewRecorder()
	s.listPostHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("got %d, want 200", w.Code)
	}
}

func TestListPostHandler_PageParam(t *testing.T) {
	posts := make([]models.Post, 10)
	for i := range posts {
		posts[i] = tu.NewPost(tu.WithSlug(fmt.Sprintf("slug-%d", i)))
	}
	s := newTestServer(t, &fakeRepo{posts: posts})
	req := httptest.NewRequest(http.MethodGet, "/posts/?page=2", nil)
	w := httptest.NewRecorder()
	s.listPostHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("got %d, want 200", w.Code)
	}
}

func TestShowPostHandler_ValidSlug(t *testing.T) {
	posts := []models.Post{
		tu.NewPost(tu.WithTitle("Hello World"), tu.WithSlug("hello-world")),
	}
	s := newTestServer(t, &fakeRepo{posts: posts})
	req := httptest.NewRequest(http.MethodGet, "/posts/hello-world", nil)
	req.SetPathValue("slug", "hello-world")
	w := httptest.NewRecorder()
	s.showPostHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("got %d, want 200", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Hello World") {
		t.Error("body missing post title")
	}
}

func TestShowPostHandler_InvalidSlug_Redirects(t *testing.T) {
	tests := []struct {
		name string
		slug string
	}{
		{"uppercase letters", "Invalid-Slug"},
		{"spaces", "has spaces"},
		{"empty", ""},
		{"too long", strings.Repeat("a", 201)},
		{"special chars", "slug!@#"},
	}
	s := newTestServer(t, &fakeRepo{})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/posts/x", nil)
			req.SetPathValue("slug", tt.slug)
			w := httptest.NewRecorder()
			s.showPostHandler(w, req)
			if w.Code != http.StatusSeeOther {
				t.Errorf("slug %q: got %d, want 303", tt.slug, w.Code)
			}
		})
	}
}

func TestShowPostHandler_SlugNotFound_Redirects(t *testing.T) {
	s := newTestServer(t, &fakeRepo{})
	req := httptest.NewRequest(http.MethodGet, "/posts/no-such-post", nil)
	req.SetPathValue("slug", "no-such-post")
	w := httptest.NewRecorder()
	s.showPostHandler(w, req)
	if w.Code != http.StatusSeeOther {
		t.Errorf("got %d, want 303", w.Code)
	}
}

func TestSearchHandler_EmptyTerm_Redirects(t *testing.T) {
	s := newTestServer(t, &fakeRepo{})
	req := httptest.NewRequest(http.MethodGet, "/search", nil)
	w := httptest.NewRecorder()
	s.searchHandler(w, req)
	if w.Code != http.StatusSeeOther {
		t.Errorf("got %d, want 303", w.Code)
	}
}

func TestSearchHandler_TermTooLong(t *testing.T) {
	s := newTestServer(t, &fakeRepo{})
	req := httptest.NewRequest(http.MethodGet, "/search?q="+strings.Repeat("a", 201), nil)
	w := httptest.NewRecorder()
	s.searchHandler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want 400", w.Code)
	}
}

func TestSearchHandler_WithResults(t *testing.T) {
	posts := []models.Post{
		tu.NewPost(tu.WithTitle("Go Programming"), tu.WithSlug("go-programming")),
		tu.NewPost(tu.WithTitle("Python Basics"), tu.WithSlug("python-basics")),
	}
	s := newTestServer(t, &fakeRepo{posts: posts})
	req := httptest.NewRequest(http.MethodGet, "/search?q=go+programming", nil)
	w := httptest.NewRecorder()
	s.searchHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("got %d, want 200", w.Code)
	}
}

func TestRSSHandler_ContentType(t *testing.T) {
	posts := []models.Post{
		tu.NewPost(tu.WithTitle("RSS Post"), tu.WithSlug("rss-post")),
	}
	s := newTestServer(t, &fakeRepo{posts: posts})
	req := httptest.NewRequest(http.MethodGet, "/rss", nil)
	w := httptest.NewRecorder()
	s.rssHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("got %d, want 200", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/xml" {
		t.Errorf("Content-Type: got %q, want %q", ct, "application/xml")
	}
	if !strings.Contains(w.Body.String(), "<?xml") {
		t.Error("body missing XML declaration")
	}
}

func TestRSSHandler_NoPosts(t *testing.T) {
	s := newTestServer(t, &fakeRepo{})
	req := httptest.NewRequest(http.MethodGet, "/rss", nil)
	w := httptest.NewRecorder()
	s.rssHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("got %d, want 200", w.Code)
	}
}

func TestIsValidSlug(t *testing.T) {
	tests := []struct {
		slug  string
		valid bool
	}{
		{"hello-world", true},
		{"abc123", true},
		{"a", true},
		{"123", true},
		{"Hello-World", false},
		{"has spaces", false},
		{"has_underscores", false},
		{"", false},
		{"special!chars", false},
		{"with/slash", false},
	}
	for _, tt := range tests {
		t.Run(tt.slug, func(t *testing.T) {
			if got := isValidSlug(tt.slug); got != tt.valid {
				t.Errorf("isValidSlug(%q) = %v, want %v", tt.slug, got, tt.valid)
			}
		})
	}
}
