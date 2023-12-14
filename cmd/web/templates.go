package main

import (
	"html/template"
	"path/filepath"
	"mayuraandrew.tech/snippetbox/pkg/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// At the moment it only contains one field, but we'll add more
// to it as the build progresses.


type templateData struct {
	Snippet *models.Snippet
	Snippets []*models.Snippet
	
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// initalize a new map to act as the cache
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// loop through the pages one-by-one
	for _, page := range pages {
		// extract the file name (like 'home.page.tmpl') from the file path
		// and assign it to the name variable.
		name := filepath.Base(page)

		// parse the page template file in to a template set.
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// use the ParseGlob method to add any 'layout' templates to the 
		// template set (in our case, it's just the "base" layout at the moment.)
		// parse the page template file in to a template set.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// use the ParseGlob method to add any 'partial' templates to the
		// template set (in our case, it;s just the 'footer' partial at the moment.)

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// add the tempalte set to the cache, using the name of the page
		// (like 'home.page.tmpl') as the key.
		cache[name] = ts
	}
	return cache, nil
}

