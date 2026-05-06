package markdown

import (
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/joshhartwig/blogo/internal/models"
)

func TestConvertMarkdownToHtml(t *testing.T) {
	tests := []struct {
		name        string
		slug        string
		data        []byte
		expectError bool
		validate    func(t *testing.T, post models.Post)
	}{
		{
			name: "valid markdown with frontmatter",
			slug: "test-post",
			data: []byte(`---
title: Test Post
date: 2024-01-15
draft: false
summary: A test post
---

# Hello World

This is a test post with **bold** and *italic* text.`),
			expectError: false,
			validate: func(t *testing.T, post models.Post) {
				if post.Metadata.Title != "Test Post" {
					t.Errorf("expected title 'Test Post', got '%s'", post.Metadata.Title)
				}
				if post.Metadata.Slug != "test-post" {
					t.Errorf("expected slug 'test-post', got '%s'", post.Metadata.Slug)
				}
				if post.Content == "" {
					t.Error("expected content, got empty string")
				}
				if post.Metadata.ReadingDuration == 0 {
					t.Error("expected reading duration to be calculated")
				}
			},
		},
		{
			name: "markdown without frontmatter",
			slug: "no-frontmatter",
			data: []byte(`# Just Content

This is content without frontmatter.`),
			expectError: false,
			validate: func(t *testing.T, post models.Post) {
				if post.Metadata.Title != "post missing frontmatter" {
					t.Errorf("expected default title, got '%s'", post.Metadata.Title)
				}
				if post.Metadata.Draft != false {
					t.Error("expected draft to be false")
				}
			},
		},
		{
			name:        "empty markdown",
			slug:        "empty",
			data:        []byte(""),
			expectError: false,
			validate: func(t *testing.T, post models.Post) {
				if post.Metadata.Title != "post missing frontmatter" {
					t.Errorf("expected default title for empty content")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post, err := convertMarkdownToHtml(tt.slug, tt.data)
			if (err != nil) != tt.expectError {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectError && tt.validate != nil {
				tt.validate(t, post)
			}
		})
	}
}

func TestGetMarkdownFromFS(t *testing.T) {
	tests := []struct {
		name          string
		fileSystem    fs.FS
		expectedCount int
	}{
		{
			name: "folder-based posts",
			fileSystem: fstest.MapFS{
				"post1/index.md": &fstest.MapFile{
					Data: []byte(`---
title: Post 1
date: 2024-01-15
draft: false
summary: First post
---

# Post 1 Content`),
				},
				"post2/index.md": &fstest.MapFile{
					Data: []byte(`---
title: Post 2
date: 2024-01-16
draft: false
summary: Second post
---

# Post 2 Content`),
				},
			},
			expectedCount: 2,
		},
		{
			name: "flat markdown files",
			fileSystem: fstest.MapFS{
				"post1.md": &fstest.MapFile{
					Data: []byte(`---
title: Flat Post 1
date: 2024-01-15
draft: false
summary: Flat post
---

# Content`),
				},
			},
			expectedCount: 1,
		},
		{
			name: "mixed folder and flat files",
			fileSystem: fstest.MapFS{
				"folder-post/index.md": &fstest.MapFile{
					Data: []byte(`---
title: Folder Post
date: 2024-01-15
draft: false
summary: In folder
---

# Content`),
				},
				"flat-post.md": &fstest.MapFile{
					Data: []byte(`---
title: Flat Post
date: 2024-01-15
draft: false
summary: Flat file
---

# Content`),
				},
			},
			expectedCount: 2,
		},
		{
			name: "draft posts excluded",
			fileSystem: fstest.MapFS{
				"published/index.md": &fstest.MapFile{
					Data: []byte(`---
title: Published
date: 2024-01-15
draft: false
summary: Published
---

# Content`),
				},
				"draft/index.md": &fstest.MapFile{
					Data: []byte(`---
title: Draft Post
date: 2024-01-15
draft: true
summary: Draft
---

# Content`),
				},
			},
			expectedCount: 1,
		},
		{
			name:          "empty filesystem",
			fileSystem:    fstest.MapFS{},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			posts, err := GetMarkdownFromFS(tt.fileSystem)
			if err != nil {
				t.Errorf("error fetching markdown from filesystem: %v", err)
			}
			if len(posts) != tt.expectedCount {
				t.Errorf("expected %d posts, got %d", tt.expectedCount, len(posts))
			}
		})
	}
}

func TestCalculateDuration(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name:     "empty content",
			content:  "",
			expected: 1,
		},
		{
			name:     "short content",
			content:  "just a few words here",
			expected: 1,
		},
		{
			name:     "medium content",
			content:  strings.Repeat("word ", 200),
			expected: 2,
		},
		{
			name:     "longer content",
			content:  strings.Repeat("word ", 500),
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration := calculateDuration(tt.content)
			if duration != tt.expected {
				t.Errorf("expected duration %d, got %d", tt.expected, duration)
			}
		})
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple string",
			input:    "hello world",
			expected: "hello-world",
		},
		{
			name:     "with file extension",
			input:    "my-post.md",
			expected: "my-post",
		},
		{
			name:     "with underscores",
			input:    "my_post_name",
			expected: "my-post-name",
		},
		{
			name:     "with mixed case",
			input:    "MyPostTitle",
			expected: "myposttitle",
		},
		{
			name:     "with special characters",
			input:    "hello@world#test!",
			expected: "helloworldtest",
		},
		{
			name:     "multiple consecutive dashes",
			input:    "hello---world",
			expected: "hello-world",
		},
		{
			name:     "leading and trailing dashes",
			input:    "-hello-world-",
			expected: "hello-world",
		},
		{
			name:     "numbers included",
			input:    "post-2024-01",
			expected: "post-2024-01",
		},
		{
			name:     "only special characters",
			input:    "!@#$%",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := slugify(tt.input)
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// Compile test to ensure types match
func TestPostStructure(t *testing.T) {
	slug := "test"
	data := []byte(`---
title: Test
date: 2024-01-15
draft: false
summary: Test
---
Content`)

	post, err := convertMarkdownToHtml(slug, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify post structure
	if _, ok := interface{}(post.Metadata).(models.PostMetadata); !ok {
		t.Error("expected post.Metadata to be models.PostMetadata")
	}

	if post.Metadata.ReadingDuration < 1 {
		t.Error("expected ReadingDuration to be at least 1")
	}
}
