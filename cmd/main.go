package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"reflect"

	staticAssets "github.com/turnerbenjamin/heterogen-go/cmd/static"
	"github.com/turnerbenjamin/heterogen-go/cmd/templates"
	"github.com/turnerbenjamin/heterogen-go/internal/database"
	"github.com/turnerbenjamin/heterogen-go/internal/dotenv"
	"github.com/turnerbenjamin/heterogen-go/internal/handlers/htmlHandler"
	"github.com/turnerbenjamin/heterogen-go/internal/handlers/middleware"
	"github.com/turnerbenjamin/heterogen-go/internal/helpers"
	"github.com/turnerbenjamin/heterogen-go/internal/hg_services"
	"github.com/turnerbenjamin/heterogen-go/internal/render"
	"github.com/turnerbenjamin/heterogen-go/internal/router"
	"github.com/vearutop/statigz"
)

func main() {
	dotenv.Load()
	mode := helpers.NewMode(os.Getenv("mode"))
	fmt.Printf("Starting Server in %s mode\n", mode)
	staticAssets.CompressFiles()

	staticFileServer := helpers.SelectValueByMode(mode, helpers.ValueSelector[http.Handler]{
		Development: http.FileServer(http.Dir("./cmd/static/")),
		Production:  statigz.FileServer(staticAssets.FileSystem),
	})

	fmt.Println(reflect.TypeOf(staticFileServer))

	//*Template config
	templateDirPaths := render.TemplateDirPaths{
		Layouts:    render.TemplateDirPath{Path: "layouts/", FileSuffix: ".layout.go.tmpl"},
		Pages:      render.TemplateDirPath{Path: "pages/", FileSuffix: ".page.go.tmpl"},
		Components: render.TemplateDirPath{Path: "components/", FileSuffix: ".component.go.tmpl"},
	}

	templateFs := helpers.SelectValueByMode(mode, helpers.ValueSelector[fs.ReadDirFS]{
		Development: helpers.GetDirFs("./cmd/templates"),
		Production:  templates.FileSystem,
	})

	log.Println(templateFs)

	doCache := mode != "development"
	if err := render.InitialiseTemplateCache(templateFs, templateDirPaths, doCache); err != nil {
		log.Fatal(err)
	}

	//*DB
	db := database.GetDB()
	defer db.Close()

	//*Middlewares
	router.Use(middleware.Logger)
	router.Use(middleware.ParseAuthCookie)

	//*Services
	authService := hg_services.NewAuthService(db)

	authHandler := htmlHandler.NewAuthHandler(authService)

	//*routes
	routes := router.Routes{}
	routes = append(routes, router.Home()...)
	routes = append(routes, router.AuthRoutes(authHandler)...)
	mux := router.GetMux(routes, staticFileServer)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
