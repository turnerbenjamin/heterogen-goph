package router

import "github.com/turnerbenjamin/heterogen-go/internal/handlers/htmlHandler"

func AuthRoutes(authHandler *htmlHandler.AuthHandler) Routes {

	return Routes{
		//REGISTRATION
		Get("/register", authHandler.RegistrationPage),
		Post("/register", authHandler.Register),

		//LOG-IN
		Get("/log-in", authHandler.LogInPage),
		Post("/log-in", authHandler.LogIn),

		//LOG-OUT
		Post("/log-out", authHandler.LogOut),
	}

}
