package router

import "net/http"

type MuxRouter struct {
	mux *http.ServeMux
}

func NewMuxRouter(mux *http.ServeMux) *MuxRouter {
	return &MuxRouter{mux: mux}
}

func (r *MuxRouter) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	r.mux.HandleFunc(pattern, handler)
}
