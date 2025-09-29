package main

import (
	"testing"
)

// TODO: readContent works in main but not in test
func TestReadContent(t *testing.T) {
	markdownCache, err := readMarkdownContent("../../content/")
	if err != nil {
		t.Errorf("unable to read markdown content from files %v", err)
	}

	if len(markdownCache) < 1 {
		t.Errorf("markdown cache is empty")
	}

}
