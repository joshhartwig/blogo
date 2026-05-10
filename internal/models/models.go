// Package models provides core data structures for the blog application,
// including Post, Pagination, and RSS feed types.
package models

import (
	"html/template"
	"time"
)

// PostListData carries the post slice and pagination for list views.
// Site config is injected automatically by render() — no need to include it here.
type PostListData struct {
	Posts      []Post
	Pagination Pagination
}

// SearchData carries search results and the original query term.
type SearchData struct {
	Posts []Post
	Term  string
}

// Post represents a blog post, containing its metadata and HTML content.
type Post struct {
	Metadata PostMetadata
	Content  template.HTML
}

// Pagination represents pagination details for a collection of posts.
type Pagination struct {
	CurrentPage int  // current page number
	NextPage    int  // next page number
	PrevPage    int  // prev page number
	PostCount   int  // post count
	PageCount   int  // our page count
	HasNext     bool // if we should show prev
	HasPrev     bool // if we should show next
	PostsStart  int  // tracks posts starting
	PostsEnd    int  // tracks posts ends
}

// NewPagination creates and returns a Pagination struct based on the provided
// current page, number of posts per page, and total post count. It ensures
// that the current page and posts per page are within valid bounds, calculates
// the total number of pages, and determines the start and end indices for the
// posts on the current page. The returned Pagination struct includes
// information about the current, next, and previous pages, as well as flags
// indicating the presence of next and previous pages.
func NewPagination(currentPage, postsPerPage, postCount int) Pagination {
	if postsPerPage <= 0 {
		postsPerPage = 5
	}

	if currentPage < 1 {
		currentPage = 1
	}

	pageCount := getPageCount(postCount, postsPerPage)

	if currentPage > pageCount && pageCount > 0 {
		currentPage = pageCount
	}

	start := (currentPage - 1) * postsPerPage
	end := min(start+postsPerPage, postCount)

	return Pagination{
		CurrentPage: currentPage,
		NextPage:    min(pageCount, currentPage+1),
		PrevPage:    max(1, currentPage-1),
		PostCount:   postCount,
		PageCount:   pageCount,
		HasNext:     currentPage < pageCount,
		HasPrev:     currentPage > 1,
		PostsStart:  start,
		PostsEnd:    end,
	}
}

func getPageCount(postCount, postsPerPage int) int {
	return (postCount + postsPerPage - 1) / postsPerPage
}

// PostMetadata represents the metadata associated with a blog post, including its title,
// publication date, tags, summary, slug, and estimated reading duration.
type PostMetadata struct {
	Title           string    `yaml:"title"`
	Date            time.Time `yaml:"date"`
	Tags            []string  `yaml:"tags"`
	Summary         string    `yaml:"summary"`
	Draft           bool      `yaml:"draft"`
	Slug            string
	ReadingDuration int
}

// RSS represents the structure of an RSS feed, containing metadata such as title, link, description, language, publication date, category, and a list of items.
type RSS struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"`
	PubDate     time.Time `xml:"pubdate"`
	Category    string    `xml:"category"`
	Item        []Item    `xml:"items"`
}

// Item represents a single entry in an RSS feed, containing metadata such as title, link, description,
// category, unique identifier (GUID), and publication date.
type Item struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Category    string    `xml:"category"`
	GUID        string    `xml:"guid"`
	PubDate     time.Time `xml:"pubDate"`
}
