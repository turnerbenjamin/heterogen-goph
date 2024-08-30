package render

import (
	"compress/gzip"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/turnerbenjamin/heterogen-go/internal/httpErrors"
)

type TemplateDirPath struct {
	Path       string
	FileSuffix string
}
type TemplateDirPaths struct {
	Layouts    TemplateDirPath
	Pages      TemplateDirPath
	Components TemplateDirPath
}

var useCache bool

var templateCache *template.Template

var directory fs.ReadDirFS
var dirPaths TemplateDirPaths

func Page(w http.ResponseWriter, r *http.Request, t string, m any, sc httpErrors.StatusCode) error {
	tmpl := getTemplates(t, dirPaths.Pages)
	return render(w, r, tmpl, m, sc)
}

func Component(w http.ResponseWriter, r *http.Request, t string, m any, sc httpErrors.StatusCode) error {
	tmpl := getTemplates(t, dirPaths.Components)
	return render(w, r, tmpl, m, sc)
}

func render(w http.ResponseWriter, r *http.Request, t *template.Template, m any, sc httpErrors.StatusCode) error {

	w.Header().Add("Content-Type", "text/html")

	if acceptsGz(r) {
		w.Header().Add("Content-Encoding", "gzip")
		gzw := gzip.NewWriter(w)
		defer gzw.Close()
		w.WriteHeader(int(sc))
		return t.Execute(gzw, m)
	}

	w.WriteHeader(int(sc))
	return t.Execute(w, m)
}

func getTemplates(t string, p TemplateDirPath) *template.Template {
	if useCache {
		log.Println("TODO: FIX CACHING")
		tmpl, err := getTemplateFromDisk(t, p)
		if err != nil {
			log.Fatal(err)
		}
		return tmpl
		// return templateCache
	} else {
		tmpl, err := getTemplateFromDisk(t, p)
		if err != nil {
			log.Fatal(err)
		}
		return tmpl
	}
}

func acceptsGz(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
}

func getTemplateFromDisk(t string, p TemplateDirPath) (*template.Template, error) {
	tp := p.Path + t + p.FileSuffix
	filesToParse := []string{tp}

	if p == dirPaths.Pages {
		filesToParse = append(filesToParse, dirPaths.Layouts.allFilesGlob(), dirPaths.Components.allFilesGlob())
	}
	return template.ParseFS(directory, filesToParse...)
}

func InitialiseTemplateCache(templateDir fs.ReadDirFS, templateDirPaths TemplateDirPaths, doCache bool) error {
	useCache = doCache
	directory = templateDir
	dirPaths = templateDirPaths

	templateCache = template.Must(
		template.ParseFS(templateDir,
			templateDirPaths.Layouts.allFilesGlob(),
			templateDirPaths.Pages.allFilesGlob(),
			templateDirPaths.Components.allFilesGlob()))

	return nil
	// dirPaths = templateDirPaths

	// paths, err := helpers.GetFilesFromDir(templateDir)
	// if err != nil {
	// 	return err
	// }

	// for _, path := range paths {
	// 	name := filepath.Base(path)
	// 	tmpl, err := template.New(name).ParseFS(directory, path)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	tmpl, err = tmpl.ParseFS(directory, dirPaths.Layouts.allFilesGlob())
	// 	if err != nil {
	// 		return err
	// 	}

	// 	tmpl, err = tmpl.ParseFS(directory, dirPaths.Components.allFilesGlob())
	// 	if err != nil {
	// 		return err
	// 	}

	// 	templateCache[name] = tmpl
	// }
	// return nil
}

func (tp *TemplateDirPath) addSuffix(sx string) string {
	return fmt.Sprintf("%s%s", tp.Path, sx)
}

func (tp *TemplateDirPath) allFilesGlob() string {
	return fmt.Sprintf("%s*%s", tp.Path, tp.FileSuffix)
}
