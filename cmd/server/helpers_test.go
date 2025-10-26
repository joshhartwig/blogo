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
