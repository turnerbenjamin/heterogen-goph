package router

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/turnerbenjamin/heterogen-go/internal/helpers"
	"github.com/turnerbenjamin/heterogen-go/internal/httpErrors"
	"github.com/turnerbenjamin/heterogen-go/internal/models"
	"github.com/turnerbenjamin/heterogen-go/internal/render"
)

type Router struct {
	Mux *http.ServeMux
}

type Middleware func(ReqHandler) ReqHandler

type ReqHandler func(http.ResponseWriter, *http.Request, *models.ResponseModal) error

func (router *Router) AddRoute(route Route) {
	pattern := fmt.Sprintf("%s %s", route.Method, route.Endpoint)
	(*router).Mux.Handle(pattern, route.Handler)
}

func GetMux(routes Routes, staticFileServer http.Handler) *http.ServeMux {
	router := Router{Mux: http.NewServeMux()}

	router.Mux.Handle("/static/", http.StripPrefix("/static/", staticFileServer))

	for _, route := range routes {
		router.AddRoute(route)
	}
	return router.Mux
}

func Handle(reqHandler ReqHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		m := models.ResponseModal{
			IsProduction: os.Getenv("mode") == string(helpers.Production),
		}
		err := reqHandler(w, r, &m)
		if err == nil {
			return
		}

		log.Println("ERROR: ", err.Error())
		httpError, ok := err.(httpErrors.HttpError)
		if !ok {
			httpError = httpErrors.ServerFail()
		}

		m.Errors = httpError.Msgs

		w.Header().Set("Status", string(httpError.StatusCode))
		if err = render.Component(w, r, "errorMessageList", m, httpError.StatusCode); err == nil {

			return
		}

		http.Error(w, "Server error", 500)
	}
}
