package posts

import (
	"fmt"
	"testing/fstest"

	"github.com/joshhartwig/blogo/internal/models"
	tu "github.com/joshhartwig/blogo/internal/testutils"
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

func createTestFS() fstest.MapFS {
	testFS := fstest.MapFS{
		"postOne.md":   {Data: []byte(testFrontMatterOne)},
		"postTwo.md":   {Data: []byte(testFrontMatterTwo)},
		"postThree.md": {Data: []byte(testFrontMatterThree)},
	}

	return testFS
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
