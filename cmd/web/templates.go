package main

import (
	"html/template"
	"path/filepath"
)

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("ui/html/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// All templates must include base + the page
		ts, err := template.ParseFiles("ui/html/base.tmpl", page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
