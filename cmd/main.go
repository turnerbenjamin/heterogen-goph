package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/turnerbenjamin/heterogen-go/internal/database"
	"github.com/turnerbenjamin/heterogen-go/internal/dotenv"
	"github.com/turnerbenjamin/heterogen-go/internal/handlers/hg_middleware"
	"github.com/turnerbenjamin/heterogen-go/internal/handlers/web_app_handlers"
	"github.com/turnerbenjamin/heterogen-go/internal/hg_services"
	"github.com/turnerbenjamin/heterogen-go/internal/jwt"
	"github.com/turnerbenjamin/heterogen-go/internal/render"
	"github.com/turnerbenjamin/heterogen-go/internal/router"
	"github.com/turnerbenjamin/heterogen-go/internal/router/routers"
)

func main() {
	dotenv.Load()

	err := render.InitialiseTemplateCache()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db := database.GetDB()
	defer db.Close()

	token, _ := jwt.Sign("123", time.Now().Add(time.Hour*24*7))
	decoded, _ := jwt.Decode(token)

	log.Println(token)
	log.Println(decoded)

	//Middlewares
	router.Use(hg_middleware.Logger)

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
