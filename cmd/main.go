package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	staticAssets "github.com/turnerbenjamin/heterogen-go/cmd/static"
	"github.com/turnerbenjamin/heterogen-go/internal/database"
	"github.com/turnerbenjamin/heterogen-go/internal/dotenv"
	"github.com/turnerbenjamin/heterogen-go/internal/handlers/hg_middleware"
	"github.com/turnerbenjamin/heterogen-go/internal/handlers/web_app_handlers"
	"github.com/turnerbenjamin/heterogen-go/internal/hg_services"
	"github.com/turnerbenjamin/heterogen-go/internal/render"
	"github.com/turnerbenjamin/heterogen-go/internal/router"
	"github.com/turnerbenjamin/heterogen-go/internal/router/routers"
)

func main() {
	dotenv.Load()
	staticAssets.CompressFiles()

	err := render.InitialiseTemplateCache()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db := database.GetDB()
	defer db.Close()

	//Middlewares
	router.Use(hg_middleware.Logger)
	router.Use(hg_middleware.PrintUserId)

	//Services
	authService := hg_services.NewAuthService(db)

	//Controllers
	authController := web_app_handlers.NewAuthController(authService)

	//routes
	routes := router.Routes{}
	routes = append(routes, routers.Home()...)
	routes = append(routes, routers.AuthRoutes(authController)...)

	mux := router.GetMux(routes)

	log.Fatal(http.ListenAndServe(":8080", mux))

}
