package main

import (
	"strings"
	"testing"
	"testing/fstest"
)

const testFrontMatterOne = `
---
title: "test1"
date: 2025-09-01
summary: "test summary"
tags:
  - "#tag1"
  - "#tag2"
---
# Test1
`

const testFrontMatterTwo = `
---
title: "test2"
date: 2025-09-01
summary: "test summary"
tags:
  - "#tag1"
  - "#tag2"
---
# Test2
`

// Tests readMarkdownContent pulls any markdown files
func TestReadContent(t *testing.T) {
	testFs := fstest.MapFS{
		"file1.md": {Data: []byte(testFrontMatterOne)},
		"file2.md": {Data: []byte(testFrontMatterTwo)},
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
