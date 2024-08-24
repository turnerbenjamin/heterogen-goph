package router

import (
	"fmt"
	"net/http"
)

type Router struct {
	Mux *http.ServeMux
}

func (router *Router) AddRoute(route Route) {
	pattern := fmt.Sprintf("%s %s", route.Method, route.Endpoint)
	(*router).Mux.Handle(pattern, route.Handler)
}
