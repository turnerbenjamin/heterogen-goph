package render

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/turnerbenjamin/heterogen-go/internal/helpers"
)

type TemplateCache = map[string]*template.Template

var templateCache = TemplateCache{}

type TemplateDirPath struct {
	Path       string
	FileSuffix string
}
type TemplateDirPaths struct {
	Layouts    TemplateDirPath
	Pages      TemplateDirPath
	Components TemplateDirPath
}

var directory fs.ReadDirFS
var dirPaths TemplateDirPaths

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

	var tp TemplateDirPath
	if strings.HasSuffix(t, dirPaths.Pages.FileSuffix) {
		tp = dirPaths.Pages
	} else if strings.HasSuffix(t, dirPaths.Components.FileSuffix) {
		tp = dirPaths.Components
	} else {
		log.Fatalf("template name (%s) must end in %s or %s ", t, dirPaths.Pages.FileSuffix, dirPaths.Components.FileSuffix)
	}

	tmpl, err := template.ParseFS(directory, tp.addSuffix(t))
	if err != nil {
		return nil, err
	}

	tmpl, err = tmpl.ParseFS(directory, dirPaths.Layouts.allFilesGlob())
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

func InitialiseTemplateCache(templateDir fs.ReadDirFS, templateDirPaths TemplateDirPaths, doCache bool) error {
	directory = templateDir
	dirPaths = templateDirPaths

	paths, err := helpers.GetFilesFromDir(templateDir)
	if err != nil {
		return err
	}

	for _, path := range paths {
		name := filepath.Base(path)
		tmpl, err := template.New(name).ParseFS(directory, path)
		if err != nil {
			return err
		}

		tmpl, err = tmpl.ParseFS(directory, dirPaths.Layouts.allFilesGlob())
		if err != nil {
			return err
		}

		templateCache[name] = tmpl
	}
	return nil
}

func (tp *TemplateDirPath) addSuffix(sx string) string {
	return fmt.Sprintf("%s%s", tp.Path, sx)
}

func (tp *TemplateDirPath) allFilesGlob() string {
	return fmt.Sprintf("%s*%s", tp.Path, tp.FileSuffix)
}
