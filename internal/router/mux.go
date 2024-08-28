package router

import (
	"log"
	"net/http"
	"os"

	staticAssets "github.com/turnerbenjamin/heterogen-go/cmd/static"
	"github.com/vearutop/statigz"
)

func GetMux(routes Routes) *http.ServeMux {
	router := Router{Mux: http.NewServeMux()}

	log.Println(os.Getenv("mode"))
	//Static Files
	var fileServer http.Handler
	if os.Getenv("mode") == "development" {
		fileServer = http.FileServer(http.Dir("./cmd/static/"))
	} else {
		fileServer = statigz.FileServer(staticAssets.FileSystem)
	}

	router.Mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	//Main paths
	for _, route := range routes {
		router.AddRoute(route)
	}
	return router.Mux
}
