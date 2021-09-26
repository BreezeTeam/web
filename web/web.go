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
	RouterGroup struct {
		prefix           string
		middlewares      []HandlerFunc //该分组的middleware
		parent           *RouterGroup
		web              *Web               //all groups share a Web instance
		htmlTemplates    *template.Template //html render
		funcMapTemplates template.FuncMap   //html render
	}
	Web struct {
		*RouterGroup //嵌入
		router       *router
		groups       []*RouterGroup //all groups
	}
)

//init Web,需要初始化最顶层的分组，然后放到groups中管理
func New() *Web {
	web := &Web{router: newRouter()}
	web.RouterGroup = &RouterGroup{web: web}
	web.groups = []*RouterGroup{web.RouterGroup}
	return web
}

//添加分组
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	web := group.web
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		web:    web,
	}
	web.groups = append(web.groups, newGroup)
	return newGroup
}

//添加中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

//添加处理静态资源的handler
func (group *RouterGroup) Static(relativePath string, root string) {
	//创建一个handler
	staticHandler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	group.GET(urlPattern, staticHandler)
}

/**
 * @Description: 静态资源处理中间件
 * @param group
 * @param relativePath
 * @param fs
 * @return HandlerFunc
 */
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	//找到绝对地址
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

//添加路由
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.web.router.addRoute(method, pattern, handler)
}

//custom render function
func (group *RouterGroup) SetFuncMap(funcMap template.FuncMap) {
	group.funcMapTemplates = funcMap
}

//根据pattern加载所有模板
func (group *RouterGroup) LoadTemplate(pattern string) {
	group.htmlTemplates = template.Must(
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
	var contextGroup *RouterGroup
	contextGroupMatchLen := -1
	for _, group := range web.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...) //找到上下文包含的group的所有中间件
		}
		if strings.HasPrefix(r.URL.Path, group.prefix) && len(group.prefix) > contextGroupMatchLen {
			contextGroupMatchLen = len(group.prefix)
			contextGroup = group
		}
	}
	//收集完所有的中间件后，就把这些函数添加到上下文对象的处理列表中
	c := newContext(w, r)
	c.handlers = middlewares
	c.group = contextGroup
	web.router.handle(c)
}

func (web *Web) Run(addr string) (err error) {
	log.Printf("Running in %s", addr)
	return http.ListenAndServe(addr, web)
}
