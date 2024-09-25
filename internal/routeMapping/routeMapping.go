package routeMapping

import (
	"github.com/turnerbenjamin/heterogen-go/internal/handlers/htmlHandler"
	"github.com/turnerbenjamin/heterogen-go/internal/handlers/jsonHandler"
	"github.com/turnerbenjamin/heterogen-go/internal/handlers/middleware"
	"github.com/turnerbenjamin/heterogen-go/internal/router"
)

func Get(authHandler *htmlHandler.AuthHandler, userHandler *htmlHandler.UsersHandler, businessHandler *htmlHandler.BusinessHandler) router.Routes {

	return router.Routes{

		//*HOME PAGES
		router.Get("/", htmlHandler.HomeHandler),
		router.Get("/dashboard", htmlHandler.DashboardHandler, middleware.RequireAuthentication()),

		//*AUTHENTICATION
		router.Get("/register", authHandler.RegistrationPage),
		router.Post("/register", authHandler.Register),
		router.Post("/log-in", authHandler.LogIn),
		router.Post("/log-out", authHandler.LogOut),

		//*USERS RESOURCES
		router.Get("/users", userHandler.UsersPage, middleware.RequireAuthentication(), middleware.RequireAdmin()),
		router.Get("/users/table", userHandler.UsersTable, middleware.RequireAuthentication(), middleware.RequireAdmin()),

		//*BUSINESSES RESOURCES
		router.Get("/add-business", businessHandler.AddBusinessPage, middleware.RequireAuthentication(), middleware.RequireAdmin()),
		//NOT FOUND
		// router.Get("/", htmlHandler.HomeHandler),

		//*API
		router.Get("/api/geocoding", jsonHandler.GeoCodingHandler, middleware.RequireAuthentication()),
	}

}
