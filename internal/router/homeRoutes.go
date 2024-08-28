package router

import (
	"github.com/turnerbenjamin/heterogen-go/internal/handlers/htmlHandler"
)

func Home() Routes {

	return Routes{
		//GET - /
		Get("/{$}", htmlHandler.HomeHandler),
	}

}
