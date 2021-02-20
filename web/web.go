package web

import (
	"log"
	"net/http"
)

//定义一种类型
//type HandlerFunc func(http.ResponseWriter, *http.Request)
type HandlerFunc func(*Context)

type Web struct {
	router *router
}

//init Web
func New() *Web {
	return &Web{router: newRouter()}
}

func (web *Web) addRouter(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, handler)
	web.router.addRoute(method, pattern, handler)
}

func (web *Web) GET(pattern string, handler HandlerFunc) {
	web.addRouter("GET", pattern, handler)
}

func (web *Web) POST(pattern string, handler HandlerFunc) {
	web.addRouter("POST", pattern, handler)
}

//将Web实现为Handler接口
func (web *Web) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	web.router.handle(c)
}

func (web *Web) Run(addr string) (err error) {
	return http.ListenAndServe(addr, web)
}
