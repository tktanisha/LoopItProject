package router

import "net/http"

type Router interface {
	Handle(pattern string, handler http.Handler)
}
