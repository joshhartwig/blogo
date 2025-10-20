package models

import (
	"testing"
)

func TestCalculatePostStartAndEnd(t *testing.T) {
	tt := []struct {
		name         string
		wantStart    int
		wantEnd      int
		postCount    int
		postsPerPage int
		currentPage  int
	}{
		{
			name:         "40 p, 5 ppp, 8 pages, 2 cp",
			wantStart:    6,
			wantEnd:      10,
			postCount:    40,
			postsPerPage: 5,
			currentPage:  2,
		},
		{
			name:         "20 p, 3 ppp, 6 pages, 2 cp",
			wantStart:    4,
			wantEnd:      6,
			postCount:    20,
			postsPerPage: 3,
			currentPage:  2,
		},
		{
			name:         "103 posts, 3 ppp, 34 pages, 13 cp",
			wantStart:    37,
			wantEnd:      39,
			postCount:    103,
			postsPerPage: 3,
			currentPage:  13,
		},
	}

	for _, tst := range tt {
		gs, ge := calculatePostStartAndEnd(tst.currentPage, tst.postsPerPage, tst.postCount)
		if gs != tst.wantStart || ge != tst.wantEnd {
			t.Errorf(`
			error with test:%s 
			postCount: %d, postsPerPage: %d currentPage: %d
			wantStart: %d gotStart: %d
			wantEnd: %d gotEnd: %d`, tst.name, tst.postCount, tst.postsPerPage, tst.currentPage, tst.wantStart, gs, tst.wantEnd, ge)
		}
	}
}

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
				PostsStart:  6,
				PostsEnd:    10,
			},
		},
		{
			name:         "Current page out of bounds (too high)",
			currentPage:  10,
			postsPerPage: 5,
			postCount:    20,
			want: Pagination{
				CurrentPage: 1,
				NextPage:    1,
				PrevPage:    1,
				PostCount:   20,
				PageCount:   4,
				HasNext:     true,
				HasPrev:     true,
				PostsStart:  1,
				PostsEnd:    5,
			},
		},
		{
			name:         "Current page out of bounds (negative)",
			currentPage:  -1,
			postsPerPage: 5,
			postCount:    20,
			want: Pagination{
				CurrentPage: 1,
				NextPage:    1,
				PrevPage:    1,
				PostCount:   20,
				PageCount:   4,
				HasNext:     true,
				HasPrev:     true,
				PostsStart:  1,
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
				NextPage:    1,
				PrevPage:    1,
				PostCount:   15,
				PageCount:   1,
				HasNext:     false,
				HasPrev:     true,
				PostsStart:  11,
				PostsEnd:    20,
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
				PostsStart:  16,
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

func TestCalculatePostStartAndEnd_MoreCases(t *testing.T) {
	tests := []struct {
		name         string
		currentPage  int
		postsPerPage int
		postCount    int
		wantStart    int
		wantEnd      int
	}{
		{
			name:         "First page, 10 per page",
			currentPage:  1,
			postsPerPage: 10,
			postCount:    100,
			wantStart:    1,
			wantEnd:      10,
		},
		{
			name:         "Second page, 10 per page",
			currentPage:  2,
			postsPerPage: 10,
			postCount:    100,
			wantStart:    11,
			wantEnd:      20,
		},
		{
			name:         "Last page, 1 per page",
			currentPage:  5,
			postsPerPage: 1,
			postCount:    5,
			wantStart:    5,
			wantEnd:      5,
		},
		{
			name:         "Zero page, 5 per page",
			currentPage:  0,
			postsPerPage: 5,
			postCount:    25,
			wantStart:    -4,
			wantEnd:      0,
		},
		{
			name:         "Negative page, 5 per page",
			currentPage:  -1,
			postsPerPage: 5,
			postCount:    25,
			wantStart:    -9,
			wantEnd:      -5,
		},
		{
			name:         "Large page, 7 per page",
			currentPage:  15,
			postsPerPage: 7,
			postCount:    105,
			wantStart:    99,
			wantEnd:      105,
		},
		{
			name:         "Single post, single page",
			currentPage:  1,
			postsPerPage: 1,
			postCount:    1,
			wantStart:    1,
			wantEnd:      1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			start, end := calculatePostStartAndEnd(tc.currentPage, tc.postsPerPage, tc.postCount)
			if start != tc.wantStart || end != tc.wantEnd {
				t.Errorf("For %q: got start=%d, end=%d; want start=%d, end=%d", tc.name, start, end, tc.wantStart, tc.wantEnd)
			}
		})
	}
}
