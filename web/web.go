package web

import (
	"fmt"
	"log"
	"net/http"
)

//定义一种类型
type HandlerFunc func(http.ResponseWriter, *http.Request)

type Web struct {
	router map[string]HandlerFunc
}

func New() *Web {
	return &Web{make(map[string]HandlerFunc)}
}
func (web *Web) addRouter(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	log.Printf("Route %4s - %s", method, handler)
	web.router[key] = handler
}

func (web *Web) GET(pattern string, handler HandlerFunc) {
	web.addRouter("GET", pattern, handler)
}

func (web *Web) POST(pattern string, handler HandlerFunc) {
	web.addRouter("POST", pattern, handler)
}
func (web *Web) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Method + "-" + r.URL.Path
	if handler, ok := web.router[key]; ok {
		handler(w, r)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL)
	}
}

func (web *Web) Run(addr string) (err error) {
	return http.ListenAndServe(addr, web)
}
