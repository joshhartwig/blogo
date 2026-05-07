package posts

import (
	"testing"
)

func TestNewPostRepo(t *testing.T) {
	testPosts := getSampleTestPosts(10)
	repo := Repository{
		Posts: testPosts,
	}

	wantPostCount := 10
	if wantPostCount != len(repo.Posts) {
		t.Errorf("Wanted a post count of %d but got %d", wantPostCount, len(repo.Posts))
	}

}

func TestSearchPosts(t *testing.T) {
	testPosts := getSampleTestPosts(10)
	repo := Repository{
		Posts: testPosts,
	}

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

func TestNewPostRepository(t *testing.T) {
	testFS := createTestFS()
	repo, err := NewPostRepository(testFS)
	if err != nil {
		t.Errorf("error creating new post repository: %v", err)
	}

	if len(repo.Posts) != len(testFS) {
		t.Errorf("post repo contents should be %d but got %d", len(testFS), len(repo.Posts))
	}
}
