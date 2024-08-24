package router

import (
	"net/http"

	"github.com/turnerbenjamin/heterogen-go/internal/handlers/hg_middleware"
	"github.com/turnerbenjamin/heterogen-go/internal/router/httpMethods"
)

type Route struct {
	Method   httpMethods.HttpMethod
	Endpoint string
	Handler  http.Handler
}

type Routes []Route

/*
wrap middeware and handler calls
*/
func chain(f http.HandlerFunc, middlewares []hg_middleware.Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

/*
Define global middleware
*/
var using = []hg_middleware.Middleware{}

func Use(middlewares ...hg_middleware.Middleware) {
	using = append(using, middlewares...)
}

var globalMiddlewareRunner hg_middleware.Middleware = func(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chain(next, using)(w, r)
	}
}

/*
Create a route with a wrapper handler containing calls to global and route specific midlewares
*/
func createRoute(method httpMethods.HttpMethod, endpoint string, handler http.HandlerFunc, middlewares ...hg_middleware.Middleware) Route {

	middlewares = append(middlewares, globalMiddlewareRunner)

	return Route{
		Method:   method,
		Endpoint: endpoint,
		Handler:  chain(handler, middlewares),
	}
}

/*
	Helper functions to create method specific routes
*/

func Get(endpoint string, handler http.HandlerFunc, middlewares ...hg_middleware.Middleware) Route {
	return createRoute(httpMethods.Get, endpoint, handler, middlewares...)
}

func Post(endpoint string, handler http.HandlerFunc, middlewares ...hg_middleware.Middleware) Route {
	return createRoute(httpMethods.Post, endpoint, handler, middlewares...)
}

func Patch(endpoint string, handler http.HandlerFunc, middlewares ...hg_middleware.Middleware) Route {
	return createRoute(httpMethods.Patch, endpoint, handler, middlewares...)
}

func Put(endpoint string, handler http.HandlerFunc, middlewares ...hg_middleware.Middleware) Route {
	return createRoute(httpMethods.Put, endpoint, handler, middlewares...)
}

func Delete(endpoint string, handler http.HandlerFunc, middlewares ...hg_middleware.Middleware) Route {
	return createRoute(httpMethods.Delete, endpoint, handler, middlewares...)
}
