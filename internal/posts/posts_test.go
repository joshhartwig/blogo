package posts

import (
	"fmt"
	"testing"

	"github.com/joshhartwig/blogo/internal/models"
	tu "github.com/joshhartwig/blogo/internal/testutils"
)

func TestNewPostRepo(t *testing.T) {
	testPosts := getSampleTestPosts(10)
	repo := NewPostRepository(testPosts)

	wantPostCount := 10
	if wantPostCount != len(repo.Posts) {
		t.Errorf("Wanted a post count of %d but got %d", wantPostCount, len(repo.Posts))
	}

}

func TestSearchPosts(t *testing.T) {
	testPosts := getSampleTestPosts(10)
	repo := NewPostRepository(testPosts)

	results := repo.SearchPosts("title 1")

	wantCount := 1
	gotCount := len(results)
	if gotCount != wantCount {
		t.Errorf("wanted %d got %d", wantCount, gotCount)
	}

	// Additional tests for SearchPosts
	results = repo.SearchPosts("content 2")

	if len(results) != 1 || results[0].Metadata.Title != "title 2" {
		t.Errorf("SearchPosts failed to find post by content")
	}

	results = repo.SearchPosts("slug3")

	if len(results) != 1 || results[0].Metadata.Title != "title 3" {
		t.Errorf("SearchPosts failed to find post by slug")
	}

	results = repo.SearchPosts("summary 4")

	if len(results) != 1 || results[0].Metadata.Title != "title 4" {
		t.Errorf("SearchPosts failed to find post by summary")
	}

	results = repo.SearchPosts("#test")

	if len(results) < 1 {
		t.Errorf("SearchPosts failed to find post by tag")
	}

	results = repo.SearchPosts("nonexistent")

	if len(results) != 0 {
		t.Errorf("SearchPosts should return 0 for nonexistent term")
	}

}

// getSampleTestPosts generates a slice of sample models.Post instances for testing purposes.
// The function creates 'c' posts, each with unique title, content, summary, slug, and reading duration.
// Other fields such as tags and date are set to default or random values.
// Returns a slice of models.Post containing the generated test posts.
func getSampleTestPosts(c int) []models.Post {
	posts := []models.Post{}
	for i := range c {
		t := tu.NewPost(
			tu.WithTitle(fmt.Sprintf("title %d", i)),
			tu.WithContent(fmt.Sprintf("content %d", i)),
			tu.WithRandomDate(),
			tu.WithTags("#test"),
			tu.WithReadingDuration(i),
			tu.WithSummary(fmt.Sprintf("summary %d", i)),
			tu.WithSlug(fmt.Sprintf("slug%d", i)),
		)
		posts = append(posts, t)
	}

	return posts
}
