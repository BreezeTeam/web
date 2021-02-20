package web

import (
	"log"
	"net/http"
)

//定义一种类型
//type HandlerFunc func(http.ResponseWriter, *http.Request)
type HandlerFunc func(*Context)

type (
	RouterGroup struct{
		prefix string
		middleware []HandlerFunc//该分组的middleware
		parent *RouterGroup
		web *Web //all groups share a Web instance
	}
	Web struct {
		*RouterGroup //嵌入
		router *router
		groups []*RouterGroup //all groups
	}
)

//init Web,需要初始化最顶层的分组，然后放到groups中管理
func New() *Web {
	web := &Web{router: newRouter()}
	web.RouterGroup = &RouterGroup{web:web}
	web.groups = []*RouterGroup{web.RouterGroup}
	return web
}
func (group *RouterGroup) Group(prefix string) *RouterGroup{
	web := group.web
	newGroup:=&RouterGroup{
		prefix:group.prefix+prefix,
		parent:group,
		web:web,
	}
	web.groups = append(web.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRouter(method string, comp string, handler HandlerFunc) {
	pattern:=group.prefix+comp
	log.Printf("Route %4s - %s", method, handler)
	group.web.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRouter("GET", pattern, handler)
}
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRouter("POST", pattern, handler)
}
func (group *RouterGroup) PUT(pattern string, handler HandlerFunc) {
	group.addRouter("PUT", pattern, handler)
}
func (group *RouterGroup) PATCH(pattern string, handler HandlerFunc) {
	group.addRouter("PATCH", pattern, handler)
}
func (group *RouterGroup) HEAD(pattern string, handler HandlerFunc) {
	group.addRouter("HEAD", pattern, handler)
}
func (group *RouterGroup) OPTIONS(pattern string, handler HandlerFunc) {
	group.addRouter("OPTIONS", pattern, handler)
}
func (group *RouterGroup) DELETE(pattern string, handler HandlerFunc) {
	group.addRouter("DELETE", pattern, handler)
}
func (group *RouterGroup) ANY(pattern string, handler HandlerFunc) {
	group.addRouter("ANY", pattern, handler)
}

//将Web实现为Handler接口
func (web *Web) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	web.router.handle(c)
}

func (web *Web) Run(addr string) (err error) {
	return http.ListenAndServe(addr, web)
}
