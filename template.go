package main

import (
	"embed"
	"html/template"
)

var (
	//go:embed index.html
	embedFS embed.FS

	tmpl *template.Template
)

type tmplData struct {
	Links  []string
	Upload bool
}

func parseTemplate() (*template.Template, error) {
	return template.ParseFS(embedFS, "*.html")
}
