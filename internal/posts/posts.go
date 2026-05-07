package posts

import (
	"errors"
	"io/fs"
	"slices"
	"strings"

	"github.com/joshhartwig/blogo/internal/markdown"
	"github.com/joshhartwig/blogo/internal/models"
)

type Repository struct {
	Posts  []models.Post
	bySlug map[string]models.Post
}

func NewPostRepository(fs fs.FS) (Repository, error) {
	posts, err := markdown.GetMarkdownFromFS(fs) // fetch posts from content directory
	if err != nil {
		return Repository{}, err
	}

	slices.SortFunc(posts, sortPostsByDate)

	return Repository{
		Posts:  posts,
		bySlug: make(map[string]models.Post),
	}, nil
}

// sortPostsByDate compares two posts by their metadata date.
// It returns -1 if the first post is newer, 1 if the second post is newer,
// and 0 if both posts have the same date.
// This function is useful for sorting posts in descending order by date.
func sortPostsByDate(x, y models.Post) int {
	if x.Metadata.Date.After(y.Metadata.Date) {
		return -1
	}
	if x.Metadata.Date.Before(y.Metadata.Date) {
		return 1
	}
	return 0
}

// GetPostsBetweenRange returns a slice of posts within the specified index range [x, y).
// It ensures the indices are within valid bounds. If the repository has no posts,
// it returns an empty slice. Negative indices are clamped to zero, and indices
// exceeding the upper bound are clamped to the last post index.
func (r *Repository) GetPostsBetweenRange(x, y int) []models.Post {
	// if we have no posts return empty slice
	if len(r.Posts) == 0 {
		return []models.Post{}
	}
	// set upper bounds
	upperBounds := len(r.Posts)

	// guard against negatives
	if x < 0 {
		x = 0
	} else if x > upperBounds {
		x = upperBounds
	}

	// gaurd against negatives
	if y < 0 {
		y = 0
	} else if y > upperBounds {
		y = upperBounds
	}

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

	// remove any spacing and lower case everything
	term = strings.TrimSpace(strings.ToLower(term))

	// create a map to hold posts that match
	temp := make(map[string]models.Post)

	for _, v := range r.Posts {
		// search content for term
		if strings.Contains(strings.ToLower(string(v.Content)), term) {
			//results = append(results, v)
			temp[v.Metadata.Slug] = v
			continue
		}

		// search slug for term
		if strings.Contains(strings.ToLower(string(v.Metadata.Slug)), term) {
			//results = append(results, v)
			temp[v.Metadata.Slug] = v
			continue
		}

		// search summary for term
		if strings.Contains(strings.ToLower(string(v.Metadata.Summary)), term) {
			//results = append(results, v)
			temp[v.Metadata.Slug] = v
			continue
		}

		// search title for term
		if strings.Contains(strings.ToLower(string(v.Metadata.Title)), term) {
			//results = append(results, v)
			temp[v.Metadata.Slug] = v
			continue
		}

		// search tags
		for _, y := range v.Metadata.Tags {
			if strings.Contains(strings.ToLower(string(y)), term) {
				//results = append(results, v)
				temp[v.Metadata.Slug] = v
				continue
			}
		}
	}

	// take from map and add to results
	for _, p := range temp {
		results = append(results, p)
	}

	// sort posts by date
	slices.SortFunc(results, sortPostsByDate)

	return results
}

// GetTopPosts returns a slice containing the top 'count' posts from the repository.
// If 'count' exceeds the number of available posts, it is adjusted to avoid out-of-bounds errors.
// The posts are returned in the order they are stored in the repository.
func (r *Repository) GetTopPosts(count int) []models.Post {
	if count < 0 {
		count = 0
	}
	if len(r.Posts) == 0 {
		return []models.Post{}
	}

	if count > len(r.Posts) {
		return r.Posts[0:len(r.Posts)]
	}

	return r.Posts[0:count]
}

func (r *Repository) GetAllPostsInOrder() []models.Post {
	if len(r.Posts) == 0 {
		return []models.Post{}
	}
	return r.Posts
}

func (r *Repository) GetPostBySlug(slug string) (models.Post, error) {
	bySlug := make(map[string]models.Post)
	for _, post := range r.Posts {
		bySlug[post.Metadata.Slug] = post
	}

	found, ok := bySlug[slug]
	if !ok {
		return models.Post{}, errors.New("no post found")
	}
	return found, nil
}

func (r *Repository) Count() int {
	return len(r.Posts)
}
