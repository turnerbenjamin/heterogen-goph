package router

import (
	"net/http"

	"github.com/turnerbenjamin/heterogen-go/internal/models"
	"github.com/turnerbenjamin/heterogen-go/internal/router/httpMethods"
)

type Route struct {
	Method   httpMethods.HttpMethod
	Endpoint string
	Handler  http.HandlerFunc
}

type Routes []Route

// Wrap middeware and handler calls
func chain(f ReqHandler, middlewares []Middleware) ReqHandler {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

// Define global middleware
var using = []Middleware{}

func Use(middlewares ...Middleware) {
	using = append(using, middlewares...)
}

var globalMiddlewareRunner Middleware = func(next ReqHandler) ReqHandler {
	return func(w http.ResponseWriter, r *http.Request, m *models.ResponseModel) error {
		err := chain(next, using)(w, r, m)
		return err
	}
}

// Create a route with a wrapper handler containing calls to global and route specific midlewares
func createRoute(method httpMethods.HttpMethod, endpoint string, handler ReqHandler, middlewares ...Middleware) Route {

	middlewares = append(middlewares, globalMiddlewareRunner)

	return Route{
		Method:   method,
		Endpoint: endpoint,
		Handler:  Handle(chain(handler, middlewares)),
	}
}

// *Helper functions to create method specific routes
func Get(endpoint string, handler ReqHandler, middlewares ...Middleware) Route {
	return createRoute(httpMethods.Get, endpoint, handler, middlewares...)
}

func Post(endpoint string, handler ReqHandler, middlewares ...Middleware) Route {
	return createRoute(httpMethods.Post, endpoint, handler, middlewares...)
}

func Patch(endpoint string, handler ReqHandler, middlewares ...Middleware) Route {
	return createRoute(httpMethods.Patch, endpoint, handler, middlewares...)
}

func Put(endpoint string, handler ReqHandler, middlewares ...Middleware) Route {
	return createRoute(httpMethods.Put, endpoint, handler, middlewares...)
}

func Delete(endpoint string, handler ReqHandler, middlewares ...Middleware) Route {
	return createRoute(httpMethods.Delete, endpoint, handler, middlewares...)
}
