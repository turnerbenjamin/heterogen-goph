package httpMethods

type HttpMethod string

const (
	Get     = HttpMethod("GET")
	Head    = HttpMethod("HEAD")
	Post    = HttpMethod("POST")
	Put     = HttpMethod("PUT")
	Delete  = HttpMethod("DELETE")
	Connect = HttpMethod("CONNECT")
	Options = HttpMethod("OPTIONS")
	Trace   = HttpMethod("Trace")
	Patch   = HttpMethod("PATCH")
)
