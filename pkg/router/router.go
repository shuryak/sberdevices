package router

import (
	"log"
	"net/http"
)

type Router struct {
	mux *http.ServeMux
	log *log.Logger
}

func New(log *log.Logger) *Router {
	return &Router{
		mux: http.NewServeMux(),
		log: log,
	}
}

type Registrable interface {
	register(r *Router)
}

func (r *Router) Add(handlers ...Registrable) {
	for _, v := range handlers {
		v.register(r)
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
