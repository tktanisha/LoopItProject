package router

import "net/http"

type MuxRouter struct {
	mux *http.ServeMux
}

func NewMuxRouter(mux *http.ServeMux) *MuxRouter {
	return &MuxRouter{mux: mux}
}

func (r *MuxRouter) Handle(pattern string, handler http.Handler) {
	r.mux.Handle(pattern, handler) // registers route
}
