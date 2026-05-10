package posts

import (
	"errors"
	"io/fs"
	"slices"
	"strings"

	"github.com/joshhartwig/blogo/internal/markdown"
	"github.com/joshhartwig/blogo/internal/models"
)

// ErrPostNotFound is returned by GetPostBySlug when no post matches the slug.
var ErrPostNotFound = errors.New("no post found")

type Repository struct {
	Posts  []models.Post
	bySlug map[string]models.Post
}

func NewPostRepository(fs fs.FS) (*Repository, error) {
	posts, err := markdown.GetMarkdownFromFS(fs) // fetch posts from content directory
	if err != nil {
		return &Repository{}, err
	}

	slices.SortFunc(posts, sortPostsByDate)

	bySlug := make(map[string]models.Post, len(posts))
	for _, p := range posts {
		bySlug[p.Metadata.Slug] = p
	}

	return &Repository{
		Posts:  posts,
		bySlug: bySlug,
	}, nil
}

func sortPostsByDate(x, y models.Post) int {
	if x.Metadata.Date.After(y.Metadata.Date) {
		return -1
	}
	if x.Metadata.Date.Before(y.Metadata.Date) {
		return 1
	}
	return 0
}

// GetPostsBetweenRange returns posts[x:y], clamping both indices to valid bounds.
func (r *Repository) GetPostsBetweenRange(x, y int) []models.Post {
	n := len(r.Posts)
	x = max(0, min(x, n))
	y = max(0, min(y, n))
	if x > y {
		x = y
	}
	return r.Posts[x:y]
}

// SearchPosts searches for posts containing the specified term in their content,
// slug, summary, title, or tags. It returns a slice of posts that match the term.
func (r *Repository) SearchPosts(term string) []models.Post {
	results := []models.Post{}
	if len(r.Posts) == 0 {
		return results
	}

	term = strings.TrimSpace(strings.ToLower(term))

	seen := make(map[string]bool)
	for _, v := range r.Posts {
		matched := strings.Contains(strings.ToLower(string(v.Content)), term) ||
			strings.Contains(strings.ToLower(v.Metadata.Slug), term) ||
			strings.Contains(strings.ToLower(v.Metadata.Summary), term) ||
			strings.Contains(strings.ToLower(v.Metadata.Title), term)

		if !matched {
			for _, tag := range v.Metadata.Tags {
				if strings.Contains(strings.ToLower(tag), term) {
					matched = true
					break
				}
			}
		}

		if matched && !seen[v.Metadata.Slug] {
			seen[v.Metadata.Slug] = true
			results = append(results, v)
		}
	}

	// sort posts by date
	slices.SortFunc(results, sortPostsByDate)

	return results
}

// GetTopPosts returns the first count posts.
func (r *Repository) GetTopPosts(count int) []models.Post {
	return r.GetPostsBetweenRange(0, count)
}

func (r *Repository) GetAllPostsInOrder() []models.Post {
	if len(r.Posts) == 0 {
		return []models.Post{}
	}
	return r.Posts
}

func (r *Repository) GetPostBySlug(slug string) (models.Post, error) {
	post, ok := r.bySlug[slug]
	if !ok {
		return models.Post{}, ErrPostNotFound
	}
	return post, nil
}

func (r *Repository) Count() int {
	return len(r.Posts)
}
