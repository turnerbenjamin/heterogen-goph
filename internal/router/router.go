package router

import (
	"fmt"
	"log"
	"net/http"

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

func Make(reqHandler ReqHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := models.ResponseModal{}
		err := reqHandler(w, r, &m)
		if err == nil {
			return
		}

		httpError, ok := err.(httpErrors.HttpError)
		if !ok {
			httpError = httpErrors.ServerFail()
		}
		log.Print(err)

		errorModal := map[string][]httpErrors.ErrorMessage{"Errors": httpError.Msgs}

		if err = render.Template(w, r, "errorMessage.component.go.tmpl", errorModal); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
	}
}
