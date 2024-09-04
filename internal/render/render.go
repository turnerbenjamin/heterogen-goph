package render

import (
	"bufio"
	"compress/gzip"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"text/template"

	"github.com/turnerbenjamin/heterogen-go/internal/helpers"
	"github.com/turnerbenjamin/heterogen-go/internal/httpErrors"
)

type TemplateConfig struct {
	FileSuffix string
}

type templateInfo struct {
	Name       string
	FileName   string
	BuildFiles []string
	Template   *template.Template
}

var templateFuncs = template.FuncMap{
	"has": func(slice []string, item string) bool {
		return slices.Contains(slice, item)
	},
}

var useCache bool

var templateInfoStore = map[string]templateInfo{}

var directory fs.ReadDirFS

var config TemplateConfig

func Page(w http.ResponseWriter, r *http.Request, t string, m any, sc httpErrors.StatusCode) error {
	fn := t + ".page" + config.FileSuffix
	tmpl := getTemplates(fn)
	return render(w, r, tmpl, m, sc)
}

func Component(w http.ResponseWriter, r *http.Request, t string, m any, sc httpErrors.StatusCode) error {
	fn := t + ".component" + config.FileSuffix
	tmpl := getTemplates(fn)
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

func getTemplates(fn string) *template.Template {
	ti, ok := templateInfoStore[fn]

	if !ok {
		log.Fatal("Template not found")
	}
	if useCache {
		return ti.Template
	} else {
		tmpl, err := getTemplateFromDisk(ti)
		if err != nil {
			log.Fatal(err)
		}
		return tmpl
	}
}

func acceptsGz(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
}

func getTemplateFromDisk(ti templateInfo) (*template.Template, error) {
	return template.New(ti.Name).Funcs(templateFuncs).ParseFS(directory, ti.BuildFiles...)
}

// Build a template store contining info about each template in the template directory
func InitialiseTemplateStore(templateDir fs.ReadDirFS, templateConfig TemplateConfig, doCache bool) error {

	useCache = doCache
	directory = templateDir
	config = templateConfig

	paths, err := helpers.GetFilesFromDir(templateDir)
	if err != nil {
		return err
	}

	//Iterate through each file in the template directory
	for _, path := range paths {

		//Skip not template files
		if !strings.Contains(path, config.FileSuffix) {
			continue
		}
		//Skip files aready processed
		if _, ok := templateInfoStore[path]; ok {
			continue
		}

		//Add template info to cache
		AddTemplateToStore(path)

	}

	return nil
}

func AddTemplateToStore(path string) {

	//Recursive build dependency generation
	shallowDependencies := parseDependencies(path)
	buildFiles := []string{path}
	for _, d := range shallowDependencies {
		if _, ok := templateInfoStore[d]; !ok {
			AddTemplateToStore(d)
		}
		buildFiles = append(buildFiles, templateInfoStore[d].BuildFiles...)
	}

	//Initialise template
	ti := templateInfo{
		Name:       filepath.Base(path),
		FileName:   path,
		BuildFiles: buildFiles,
	}

	//Cache template build if useCache is true
	if useCache {
		tmpl, err := template.New(ti.Name).Funcs(templateFuncs).ParseFS(directory, ti.BuildFiles...)
		if err != nil {
			log.Fatal(err)
		}
		ti.Template = tmpl
	}

	//Add template to store
	templateInfoStore[ti.Name] = ti
}

// Search file for template declarations and add these to dependencies
func parseDependencies(filepath string) []string {
	dependencies := []string{}

	file, err := directory.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	templatePattern := regexp.MustCompile(`template "(?P<name>.*)"`)

	for scanner.Scan() {
		if templatePattern.MatchString(scanner.Text()) {
			declaration := templatePattern.FindStringSubmatch(scanner.Text())
			if len(declaration) != 2 {
				continue
			}
			dependency := declaration[1] + ".go.tmpl"
			dependencies = append(dependencies, dependency)
		}
	}

	return dependencies

}
