package router

import (
	"net/http"

	"github.com/turnerbenjamin/heterogen-go/internal/handlers/hg_middleware"
	"github.com/turnerbenjamin/heterogen-go/internal/httpErrors"
	"github.com/turnerbenjamin/heterogen-go/internal/router/httpMethods"
)

type Route struct {
	Method   httpMethods.HttpMethod
	Endpoint string
	Handler  http.HandlerFunc
}

type Routes []Route

// Wrap middeware and handler calls
func chain(f httpErrors.ReqHandler, middlewares []hg_middleware.Middleware) httpErrors.ReqHandler {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

// Define global middleware
var using = []hg_middleware.Middleware{}

func Use(middlewares ...hg_middleware.Middleware) {
	using = append(using, middlewares...)
}

var globalMiddlewareRunner hg_middleware.Middleware = func(next httpErrors.ReqHandler) httpErrors.ReqHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		err := chain(next, using)(w, r)
		return err
	}
}

// Create a route with a wrapper handler containing calls to global and route specific midlewares
func createRoute(method httpMethods.HttpMethod, endpoint string, handler httpErrors.ReqHandler, middlewares ...hg_middleware.Middleware) Route {

	middlewares = append(middlewares, globalMiddlewareRunner)

	return Route{
		Method:   method,
		Endpoint: endpoint,
		Handler:  httpErrors.HandleErrors(chain(handler, middlewares)),
	}
}

// *Helper functions to create method specific routes
func Get(endpoint string, handler httpErrors.ReqHandler, middlewares ...hg_middleware.Middleware) Route {
	return createRoute(httpMethods.Get, endpoint, handler, middlewares...)
}

func Post(endpoint string, handler httpErrors.ReqHandler, middlewares ...hg_middleware.Middleware) Route {
	return createRoute(httpMethods.Post, endpoint, handler, middlewares...)
}

func Patch(endpoint string, handler httpErrors.ReqHandler, middlewares ...hg_middleware.Middleware) Route {
	return createRoute(httpMethods.Patch, endpoint, handler, middlewares...)
}

func Put(endpoint string, handler httpErrors.ReqHandler, middlewares ...hg_middleware.Middleware) Route {
	return createRoute(httpMethods.Put, endpoint, handler, middlewares...)
}

func Delete(endpoint string, handler httpErrors.ReqHandler, middlewares ...hg_middleware.Middleware) Route {
	return createRoute(httpMethods.Delete, endpoint, handler, middlewares...)
}
