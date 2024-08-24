package routers

import (
	"github.com/turnerbenjamin/heterogen-go/internal/handlers/web_app_handlers"
	"github.com/turnerbenjamin/heterogen-go/internal/router"
)

func Home() router.Routes {

	return router.Routes{
		//GET - /
		router.Get("/{$}", web_app_handlers.HomeHandler),
	}

}
