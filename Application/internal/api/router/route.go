package router

import "net/http"

//go:generate mockgen -source=route.go -destination=../../mock/mock_route.go -package=mock

type Router interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}
