package main

import (
	"strings"
	"testing"
)

func TestCalculateDuration(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    int
	}{
		{"empty", "", 1},
		{"short", "hello world", 1},
		{"200 words", strings.Repeat("word ", 200), 2},
		{"500 words", strings.Repeat("word ", 500), 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateDuration(tt.content)
			if got != tt.want {
				t.Errorf("calcuateDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World.md", "hello-world"},
		{"My_post_title.md", "my-post-title"},
		{"  Spaces  and___underscores.md", "spaces-and-underscores"},
		{"Special!@#Characters.md", "specialcharacters"},
		{"MixedCASE123.md", "mixedcase123"},
		{"trailing-.md", "trailing"},
		{"--double--dash--.md", "double-dash"},
		{"file with spaces.md", "file-with-spaces"},
		{"file_name_with_underscores.md", "file-name-with-underscores"},
	}

	for _, tt := range tests {
		got := slugify(tt.input)
		if got != tt.expected {
			t.Errorf("slugify(%q) got %q; want %q", tt.input, got, tt.expected)
		}
	}
}
