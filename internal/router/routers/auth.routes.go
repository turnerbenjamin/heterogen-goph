package routers

import (
	"github.com/turnerbenjamin/heterogen-go/internal/handlers/web_app_handlers"
	"github.com/turnerbenjamin/heterogen-go/internal/router"
)

func AuthRoutes(authController *web_app_handlers.AuthController) router.Routes {

	return router.Routes{
		//REGISTRATION
		router.Get("/register", authController.RegistrationPage),
		router.Post("/register", authController.Register),

		//LOG-IN
		router.Get("/log-in", authController.LogInPage),
		router.Post("/log-in", authController.LogIn),
	}

}
