package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func (app *application) render(w http.ResponseWriter, name string, data any) error {
	v, ok := app.templateCache[name]
	if !ok {
		return errors.New("template not found")
	}

	err := v.ExecuteTemplate(w, "base", data)
	if err != nil {

		return err
	}
	return nil
}

func readMarkdownContent() (map[string][]byte, error) {
	markdown := make(map[string][]byte)
	dir, _ := os.Getwd()
	fmt.Println(dir)
	files, err := filepath.Glob("./content/*.md")
	fmt.Println("found files:", files)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		filename := filepath.Base(file)             // get just the base file example.html
		name := strings.TrimSuffix(filename, ".md") // get the title name for the map
		bytes, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		markdown[name] = convertMarkdownToHtml(bytes) // add to map and convert mardown to html
	}

	return markdown, nil
}

// converts markdown content to html
func convertMarkdownToHtml(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	render := html.NewRenderer(opts)

	return markdown.Render(doc, render)
}
