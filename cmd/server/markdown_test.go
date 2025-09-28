package main

import (
	"fmt"
	"testing"
)

// TODO: readContent works in main but not in test
func TestReadContent(t *testing.T) {
	mdCache, err := readMarkdownContent("./content/")
	if err != nil {
		t.Errorf("unable to read markdown content from files %v", err)
	}

	if mdCache == nil {
		t.Errorf("markdown cache is nil")
	}
	fmt.Println(mdCache)
}
