package render

import (
	"compress/gzip"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type TemplateCache = map[string]*template.Template

var templateCache = TemplateCache{}

var doCache = false

func Template(w http.ResponseWriter, r *http.Request, t string, data any) error {

	parsedTemplate, err := getTemplate(t)
	if err != nil {
		return err
	}

	if acceptsGz(r) {
		w.Header().Add("Content-Type", "text/html")
		w.Header().Add("Content-Encoding", "gzip")

		gzw := gzip.NewWriter(w)
		defer gzw.Close()
		return parsedTemplate.Execute(gzw, data)

	}
	return parsedTemplate.Execute(w, data)
}

func acceptsGz(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
}

func getTemplate(t string) (*template.Template, error) {
	if !doCache {
		return getTemplateFromDisk(t)
	}
	return getTemplateFromCache(t)
}

func getTemplateFromDisk(t string) (*template.Template, error) {

	filepath := "./web/views/pages/"
	if strings.Contains(t, ".component") {
		filepath = "./web/views/components/"
	}

	tmpl, err := template.ParseFiles(filepath + t)

	if err != nil {
		return nil, err
	}

	tmpl, err = tmpl.ParseGlob("./web/views/layouts/*.layout.go.tmpl")
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func getTemplateFromCache(t string) (*template.Template, error) {
	tmpl, inMap := templateCache[t]
	if !inMap {
		return nil, errors.New("template not found")
	}
	return tmpl, nil

}

func InitialiseTemplateCache() error {

	doCache = os.Getenv("mode") == "production"

	pagePaths, err := filepath.Glob("./web/views/pages/*.page.go.tmpl")
	if err != nil {
		return err
	}

	componentPaths, err := filepath.Glob("./web/views/components/*.component.go.tmpl")
	if err != nil {
		return err
	}

	paths := append(pagePaths, componentPaths...)

	for _, path := range paths {
		name := filepath.Base(path)
		tmpl, err := template.New(name).ParseFiles(path)
		if err != nil {
			return err
		}

		tmpl, err = tmpl.ParseGlob("./web/views/layouts/*.layout.go.tmpl")
		if err != nil {
			return err
		}

		templateCache[name] = tmpl
	}
	return nil
}
