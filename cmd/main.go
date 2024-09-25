package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	staticAssets "github.com/turnerbenjamin/heterogen-go/cmd/static"
	"github.com/turnerbenjamin/heterogen-go/cmd/templates"
	"github.com/turnerbenjamin/heterogen-go/internal/database"
	"github.com/turnerbenjamin/heterogen-go/internal/dotenv"
	"github.com/turnerbenjamin/heterogen-go/internal/handlers/htmlHandler"
	"github.com/turnerbenjamin/heterogen-go/internal/handlers/middleware"
	"github.com/turnerbenjamin/heterogen-go/internal/helpers"
	"github.com/turnerbenjamin/heterogen-go/internal/hg_services"
	"github.com/turnerbenjamin/heterogen-go/internal/render"
	"github.com/turnerbenjamin/heterogen-go/internal/routeMapping"
	"github.com/turnerbenjamin/heterogen-go/internal/router"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

func main() {
	dotenv.Load()
	mode := helpers.NewMode(os.Getenv("mode"))
	fmt.Printf("Starting Server in %s mode\n", mode)
	staticAssets.CompressFiles()

	staticFileServer := helpers.SelectValueByMode(mode, helpers.ValueSelector[http.Handler]{
		Development: http.FileServer(http.Dir("./cmd/static/")),
		Production:  statigz.FileServer(staticAssets.FileSystem, brotli.AddEncoding),
	})

	//*Template config
	templateDirConfig := render.TemplateConfig{
		FileSuffix: ".go.tmpl",
	}

	templateFs := helpers.SelectValueByMode(mode, helpers.ValueSelector[fs.ReadDirFS]{
		Development: helpers.GetDirFs("./cmd/templates"),
		Production:  templates.FileSystem,
	})

	doCache := mode != "development"
	if err := render.InitialiseTemplateStore(templateFs, templateDirConfig, doCache); err != nil {
		log.Fatal(err)
	}

	//*DB
	db := database.GetDB()
	defer db.Close()

	//*Services
	authService := hg_services.NewAuthService(db)
	userService := hg_services.NewUsersService(db)
	businessService := hg_services.NewBusinessServiceService(db)

	//*Handlers
	authHandler := htmlHandler.NewAuthHandler(authService)
	userHandler := htmlHandler.NewUsersHandler(userService)
	businessHandler := htmlHandler.NewBusinessesHandler(businessService)

	//*Middlewares
	router.Use(middleware.Logger)
	router.Use(middleware.GetUserAuthenticator(authService))

	//*routes
	routeMapping := routeMapping.Get(authHandler, userHandler, businessHandler)
	mux := router.GetMux(routeMapping, staticFileServer)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
