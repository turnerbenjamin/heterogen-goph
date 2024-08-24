package render

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type TemplateCache = map[string]*template.Template

var templateCache = TemplateCache{}

var doCache = false

func Template(w http.ResponseWriter, t string, data any) {
	parsedTemplate := getTemplate(t)
	err := parsedTemplate.Execute(w, data)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
}

func getTemplate(t string) *template.Template {
	if !doCache {
		return getTemplateFromDisk(t)
	}
	return getTemplateFromCache(t)
}

func getTemplateFromDisk(t string) *template.Template {

	filepath := "./web/views/pages/"
	if strings.Contains(t, ".component") {
		filepath = "./web/views/components/"
	}

	tmpl, err := template.ParseFiles(filepath + t)

	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}

	tmpl, err = tmpl.ParseGlob("./web/views/layouts/*.layout.tmpl")
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}

	return tmpl
}

func getTemplateFromCache(t string) *template.Template {
	tmpl, inMap := templateCache[t]
	if !inMap {
		fmt.Printf("Template (%s) not found in cache", t)
		os.Exit(1)
	}
	return tmpl

}

func InitialiseTemplateCache() error {

	doCache = os.Getenv("mode") == "production"

	pagePaths, err := filepath.Glob("./web/views/pages/*.page.tmpl")
	if err != nil {
		return err
	}

	componentPaths, err := filepath.Glob("./web/views/components/*.component.tmpl")
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

		tmpl, err = tmpl.ParseGlob("./web/views/layouts/*.layout.tmpl")
		if err != nil {
			return err
		}

		templateCache[name] = tmpl
	}
	return nil
}
