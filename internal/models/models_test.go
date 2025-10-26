package models

import (
	"testing"
)

func TestNewPaginationReturnsProperValues(t *testing.T) {
	tests := []struct {
		name         string
		currentPage  int
		postsPerPage int
		postCount    int
		want         Pagination
	}{
		{
			name:         "Normal case, middle page",
			currentPage:  2,
			postsPerPage: 5,
			postCount:    20,
			want: Pagination{
				CurrentPage: 2,
				NextPage:    3,
				PrevPage:    1,
				PostCount:   20,
				PageCount:   4,
				HasNext:     true,
				HasPrev:     true,
				PostsStart:  5,
				PostsEnd:    10,
			},
		},
		{
			name:         "Current page out of bounds (too high)",
			currentPage:  10,
			postsPerPage: 5,
			postCount:    20,
			want: Pagination{
				CurrentPage: 4,
				NextPage:    4,
				PrevPage:    3,
				PostCount:   20,
				PageCount:   4,
				HasNext:     false,
				HasPrev:     true,
				PostsStart:  15,
				PostsEnd:    20,
			},
		},
		{
			name:         "Current page out of bounds (negative)",
			currentPage:  -1,
			postsPerPage: 5,
			postCount:    20,
			want: Pagination{
				CurrentPage: 1,
				NextPage:    2,
				PrevPage:    1,
				PostCount:   20,
				PageCount:   4,
				HasNext:     true,
				HasPrev:     false,
				PostsStart:  0,
				PostsEnd:    5,
			},
		},
		{
			name:         "First page, has next only",
			currentPage:  1,
			postsPerPage: 10,
			postCount:    15,
			want: Pagination{
				CurrentPage: 1,
				NextPage:    2,
				PrevPage:    1,
				PostCount:   15,
				PageCount:   2,
				HasNext:     true,
				HasPrev:     false,
				PostsStart:  0,
				PostsEnd:    10,
			},
		},
		{
			name:         "Last page, no next",
			currentPage:  4,
			postsPerPage: 5,
			postCount:    20,
			want: Pagination{
				CurrentPage: 4,
				NextPage:    4,
				PrevPage:    3,
				PostCount:   20,
				PageCount:   4,
				HasNext:     false,
				HasPrev:     true,
				PostsStart:  15,
				PostsEnd:    20,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := NewPagination(tc.currentPage, tc.postsPerPage, tc.postCount)
			if got.CurrentPage != tc.want.CurrentPage {
				t.Errorf("CurrentPage: got %d, want %d", got.CurrentPage, tc.want.CurrentPage)
			}
			if got.NextPage != tc.want.NextPage {
				t.Errorf("NextPage: got %d, want %d", got.NextPage, tc.want.NextPage)
			}
			if got.PrevPage != tc.want.PrevPage {
				t.Errorf("PrevPage: got %d, want %d", got.PrevPage, tc.want.PrevPage)
			}
			if got.PostCount != tc.want.PostCount {
				t.Errorf("PostCount: got %d, want %d", got.PostCount, tc.want.PostCount)
			}
			if got.PageCount != tc.want.PageCount {
				t.Errorf("PageCount: got %d, want %d", got.PageCount, tc.want.PageCount)
			}
			if got.HasNext != tc.want.HasNext {
				t.Errorf("HasNext: got %v, want %v", got.HasNext, tc.want.HasNext)
			}
			if got.HasPrev != tc.want.HasPrev {
				t.Errorf("HasPrev: got %v, want %v", got.HasPrev, tc.want.HasPrev)
			}
			if got.PostsStart != tc.want.PostsStart {
				t.Errorf("PostsStart: got %d, want %d", got.PostsStart, tc.want.PostsStart)
			}
			if got.PostsEnd != tc.want.PostsEnd {
				t.Errorf("PostsEnd: got %d, want %d", got.PostsEnd, tc.want.PostsEnd)
			}
		})
	}
}
