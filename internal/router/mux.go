package router

import (
	"net/http"
)

func GetMux(routes Routes) *http.ServeMux {
	router := Router{Mux: http.NewServeMux()}

	//Static Files
	fileServer := http.FileServer(http.Dir("web/static/assets"))
	router.Mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	//Main paths
	for _, route := range routes {
		router.AddRoute(route)
	}
	return router.Mux
}
