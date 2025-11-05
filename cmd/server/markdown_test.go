package main

import (
	"fmt"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/joshhartwig/blogo/internal/models"
)

const testFrontMatterOne = `
---
title: "test1"
date: 2025-09-01
summary: "test summary"
draft: false
tags:
  - "#tag1"
  - "#tag2"
---
# Test1
`

const testFrontMatterTwo = `
---
title: "test2"
date: 2025-09-02
summary: "test summary"
draft: false
tags:
  - "#tag1"
  - "#tag2"
---
# Test2
`

const testFrontMatterThree = `
# Test3
`

const testFrontMatterFour = `
---
title: "test4"
date: 2025-09-06
summary: "test summary"
draft: true
tags:
  - "#tag1"
  - "#tag2"
---
# Test4
`

func TestConvertMarkdownToHTML(t *testing.T) {
	want := models.Post{
		Content: `<h1>Test1</h1>`,
		Metadata: models.PostMetadata{
			Title: "test1",
		},
	}

	got, err := convertMarkdownToHtml("test1", []byte(testFrontMatterOne))
	if err != nil {
		t.Errorf("got error: %v when converting markdown to html", err)
	}

	// note we need to trim the string as goldmark will add a trailing line feed
	if strings.Compare(strings.TrimSpace(string(got.Content)), strings.TrimSpace(string(want.Content))) != 0 {
		t.Errorf("got %s want %s", got.Content, want.Content)
	}

	if got.Metadata.Title != want.Metadata.Title {
		t.Errorf("got title:%s\n, want title:%s\n", got.Metadata.Title, want.Metadata.Title)
	}
}

func TestReadMarkdownReturnPostsInOrder(t *testing.T) {
	testFs := fstest.MapFS{
		"file1.md": {Data: []byte(testFrontMatterOne)},
		"file2.md": {Data: []byte(testFrontMatterTwo)},
	}

	sortedPosts, err := readMarkdownReturnPostsInOrder(testFs)
	if err != nil {
		t.Errorf("should not have received error %v", err)
	}

	if len(sortedPosts) < len(testFs) {
		t.Error("should be more then 2 posts")
	}

	fmt.Println(sortedPosts[0].Metadata.Date)
}

// Tests readMarkdownContent pulls any markdown files
func TestReadMarkdownContent(t *testing.T) {
	testFs := fstest.MapFS{
		"file1.md": {Data: []byte(testFrontMatterOne)},
		"file2.md": {Data: []byte(testFrontMatterTwo)},
		"file3.md": {Data: []byte(testFrontMatterThree)}, // missing front matter
		"file4.md": {Data: []byte(testFrontMatterFour)},  // draft = true should not show
	}
	markdownCache, err := readMarkdownContent(testFs)
	if err != nil {
		t.Errorf("unable to read markdown content from files %v", err)
	}

	if len(markdownCache) < 1 {
		t.Errorf("markdown cache is empty")
	}

	tt := []struct {
		wantSlug         string
		wantDataContains string
	}{
		{
			wantSlug:         "file1",
			wantDataContains: "<h1>Test1</h1>",
		},
		{
			wantSlug:         "file2",
			wantDataContains: "<h1>Test2</h1>",
		},
		{
			wantSlug:         "file3",
			wantDataContains: "<h1>Test3</h1>", // ensure file with no front matter is processed
		},
	}

	for _, test := range tt {
		post, ok := markdownCache[test.wantSlug]
		if !ok {
			t.Errorf("expected to find slug:%s in markdown cache", test.wantSlug)
		}

		if !strings.Contains(string(post.Content), test.wantDataContains) {
			t.Errorf("expected test data: %s\n got: %s\n slug: %s", test.wantDataContains, post.Content, test.wantSlug)
		}
	}

}
