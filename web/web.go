package web

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

//定义一种类型
//type HandlerFunc func(http.ResponseWriter, *http.Request)
type HandlerFunc func(*Context)

type (
	RouterGroup struct{
		prefix 		string
		middlewares []HandlerFunc//该分组的middleware
		parent 		*RouterGroup
		web 		*Web //all groups share a Web instance
		//todo 这里htmlTemplates是不是没有隔离性？？？
		htmlTemplates template.Template //html render
		funcMapTemplates template.FuncMap //html render
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
//添加分组
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
//添加中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc){
	group.middlewares = append(group.middlewares,middlewares...)
}
//添加处理静态资源的handler
func (group *RouterGroup)Static(relativePath string,root string){
	//创建一个handler
	staticHandler :=createStaticHandler(group,relativePath,http.Dir(root))
	urlPattern:=path.Join(relativePath,"/*filepath")
	group.GET(urlPattern,staticHandler)
}

//添加路由
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern:=group.prefix+comp
	log.Printf("Route %4s - %s", method, handler)
	group.web.router.addRoute(method, pattern, handler)
}
//custom render function
func (group *RouterGroup) SetFuncMap(funcMap template.FuncMap){
	group.funcMapTemplates = funcMap
}
//根据pattern加载所有模板
func (group *RouterGroup) LoadTemplate(pattern string){
	group.htmlTemplates = *template.Must(
		template.New("").Funcs(group.funcMapTemplates).ParseGlob(pattern))
}


func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}
func (group *RouterGroup) PUT(pattern string, handler HandlerFunc) {
	group.addRoute("PUT", pattern, handler)
}
func (group *RouterGroup) PATCH(pattern string, handler HandlerFunc) {
	group.addRoute("PATCH", pattern, handler)
}
func (group *RouterGroup) HEAD(pattern string, handler HandlerFunc) {
	group.addRoute("HEAD", pattern, handler)
}
func (group *RouterGroup) OPTIONS(pattern string, handler HandlerFunc) {
	group.addRoute("OPTIONS", pattern, handler)
}
func (group *RouterGroup) DELETE(pattern string, handler HandlerFunc) {
	group.addRoute("DELETE", pattern, handler)
}
func (group *RouterGroup) ANY(pattern string, handler HandlerFunc) {
	group.addRoute("ANY", pattern, handler)
}

//将Web实现为Handler接口
func (web *Web) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//循环web的每一个分组，如果请求的url包含group的前缀，那么就把这个分组定义的中间件收集起来
	var middlewares []HandlerFunc
	for _,group := range web.groups {
		if strings.HasPrefix(r.URL.Path,group.prefix){
			middlewares = append(middlewares,group.middlewares...)
		}
	}
	//收集完所有的中间件后，就把这些函数添加到上下文对象的处理列表中
	c := newContext(w, r)
	c.handlers = middlewares
	//todo 这里是不是应该让我们的上下文，能够知道自己所在的组？？
	c.web = web
	web.router.handle(c)
}

func (web *Web) Run(addr string) (err error) {
	return http.ListenAndServe(addr, web)
}
