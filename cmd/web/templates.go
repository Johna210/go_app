package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"johna.net/snippetbox/internal/models"
	"johna.net/snippetbox/ui"
)

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

// newTemplateCache creates a new template cache by parsing the template files in the specified directory.
// It returns a map of template names to template pointers and an error if any occurred during the parsing process.
// The template files are expected to be located in the "./ui/html/pages/" directory.
// The base template file "./ui/html/base.tmpl" is parsed and added to each template set.
// Additionally, any partial templates located in the "./ui/html/partials/" directory are also parsed and added to each template set.
// The template files in the "./ui/html/pages/" directory are parsed individually and added to the cache with their respective names.
// The cache is a map where the keys are the names of the template files and the values are the corresponding template pointers.
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil

}
