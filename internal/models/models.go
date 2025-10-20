package models

import (
	"html/template"
	"time"
)

type TemplateData struct {
	Title string
	Posts []Post
}

type HomePageData struct {
	Title          string
	Posts          []Post
	PaginationData PaginationData
}

type Post struct {
	Metadata PostMetadata
	Content  template.HTML
}

// Pagination represents pagination details for a collection of posts.
// It includes information about the current, next, and previous pages,
// the total number of posts and pages, and flags indicating the presence
// of next and previous pages. It also tracks the range of posts displayed
// on the current page.
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

	pageCount := (postCount + postsPerPage - 1) / postsPerPage
	if pageCount < 1 {
		pageCount = 1
	}

	// check the bounds of current page and make sure it falls between our post counts
	if currentPage < 1 {
		currentPage = 1
	}

	start, end := calculatePostStartAndEnd(currentPage, postsPerPage, postCount)

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

// calculatePostStartAndEnd calculates the start and end indices for paginating a list of posts.
// It takes the current page number, the number of posts per page, and the total post count as input.
// The function returns the start and end indices (1-based, inclusive) for the posts to display on the current page.
// If the calculated end index exceeds the total post count, it is capped at postCount.
func calculatePostStartAndEnd(currentPage, postsPerPage int, postCount int) (start, end int) {
	start = (currentPage-1)*postsPerPage + 1
	end = postsPerPage * currentPage
	if end > postCount {
		end = postCount
	}

	return
}

// min returns the smaller of two integers a and b.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the greater of two integer values a and b.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type PaginationData struct {
	CurrentPage  int
	PageCount    int
	HasNext      bool
	HasPrev      bool
	NextPage     int
	PrevPage     int
	PostCount    int
	PostsPerPage int
}

// PostMetadata represents the metadata associated with a blog post, including its title,
// publication date, tags, summary, slug, and estimated reading duration.
type PostMetadata struct {
	Title           string    `yaml:"title"`
	Date            time.Time `yaml:"date"`
	Tags            []string  `yaml:"tags"`
	Summary         string    `yaml:"summary"`
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
