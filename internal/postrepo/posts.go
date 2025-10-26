package postrepo

import (
	"fmt"
	"slices"
	"strings"

	"github.com/joshhartwig/blogo/internal/models"
)

type PostRepository struct {
	Posts []models.Post
}

func NewPostRepository(posts []models.Post) PostRepository {
	sorted := []models.Post{}
	sorted = append(sorted, posts...)
	slices.SortFunc(sorted, sortPostsByDate)

	return PostRepository{
		Posts: sorted,
	}
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
func (p *PostRepository) GetPostsBetweenRange(x, y int) []models.Post {
	// if we have no posts return empty slice
	if len(p.Posts) == 0 {
		return []models.Post{}
	}
	// set upper bounds
	upperBounds := len(p.Posts)

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

	return p.Posts[x:y]
}

// SearchPosts searches for posts containing the specified term in their content,
// slug, summary, title, or tags. It returns a slice of posts that match the term.
func (p *PostRepository) SearchPosts(term string) []models.Post {
	results := []models.Post{}
	if len(p.Posts) == 0 {
		return results
	}

	// remove any spacing and lower case everything
	term = strings.TrimSpace(strings.ToLower(term))

	for _, v := range p.Posts {
		if strings.Contains(strings.ToLower(string(v.Content)), term) {
			results = append(results, v)
			continue
		}

		if strings.Contains(strings.ToLower(string(v.Metadata.Slug)), term) {
			results = append(results, v)
			continue
		}

		if strings.Contains(strings.ToLower(string(v.Metadata.Summary)), term) {
			results = append(results, v)
			continue
		}

		if strings.Contains(strings.ToLower(string(v.Metadata.Title)), term) {
			results = append(results, v)
			continue
		}

		for _, y := range v.Metadata.Tags {
			if strings.Contains(strings.ToLower(string(y)), term) {
				results = append(results, v)
				continue
			}
		}
	}

	return results
}

// GetTopPosts returns a slice containing the top 'count' posts from the repository.
// If 'count' exceeds the number of available posts, it is adjusted to avoid out-of-bounds errors.
// The posts are returned in the order they are stored in the repository.
func (p *PostRepository) GetTopPosts(count int) []models.Post {
	if len(p.Posts) == 0 {
		return []models.Post{}
	}

	if count > len(p.Posts) {
		return p.Posts[0 : len(p.Posts)-1]
	}
	fmt.Printf("GetTopPosts: count:%d postCount:%d", count, len(p.Posts))
	return p.Posts[0:count]
}

func (p *PostRepository) GetAllPostsInOrder() []models.Post {
	if len(p.Posts) == 0 {
		return []models.Post{}
	}
	return p.Posts
}
