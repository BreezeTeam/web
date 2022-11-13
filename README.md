# web
web by go
框架代码位于web目录中

how to use

`git clone this repo`

`cd in dir`

`go run .`

test:
you can find test code in `main.go`


### 特性1. 仿Gin 接口设计
```go
w := web.New()
v1 := w.Group("/v1")
v1.GET("/panic", func(c *web.Context) {
    names := []string{"geektutu"}
    c.STRING(http.StatusOK, names[100])
})
```

### 特性2. 基于Trie实现的动态路由
实现了基于trie 前缀树的动态路由，并且支持了两种模式:name和*filepath
```go
type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isFuzzy  bool    // 是否模糊匹配，part 含有 : 或 * 时为true
}
//使用bfs来寻找parts对应的node
func (n *node)search(parts []string,height int) *node{
	//递归结束条件,这里不能直接取第一个字符，因为有可能第一个字符不存在
	//如果是*开始的，那么就可以把这里设置为一个节点，之后只要匹配到这节点的，都能匹配成功
	if len(parts) == height || strings.HasPrefix(n.part, "*"){
		//如果pattern，说明，这个路径没对应的method，匹配失败
		if n.pattern == ""{
			return nil
		}
		return n
	}
	part:=parts[height]
	children:=n.matchChildren(part)
	//循环该节点的每一个子节点，其中，对每一个节点再进行搜索
	for _, child := range children {
		result :=child.search(parts,height+1)
		if result!=nil{
			return result
		}
	}
	return nil
}

//插入
func (n *node) insert(pattern string, parts []string, height int) {
	//递归结束，插入完成
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	//如果children为空，表示该节点还没有对应的子节点，需要进行创建
	if child == nil {
		child = &node{part: part, isFuzzy: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	//进行下一层的插入
	child.insert(pattern,parts,height+1)
}

//通过part匹配节点的第一个子节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isFuzzy {
			return child
		}
	}
	return nil
}

//通过part匹配到全部的子节点
func (n *node) matchChildren(part string) []*node {
	children := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isFuzzy {
			children = append(children, child)
		}
	}
	return children
}
```

domo 
```go
func main() {
	w := web.New()
	//curl http://localhost:9999/
	w.GET("/", func(c *web.Context) {
		c.HTML(http.StatusOK, "<h1>Hello World</h1>")
	})
	//curl http://localhost:9999/hello?name=Euraxluo
	w.GET("/hello", func(c *web.Context) {
		c.STRING(http.StatusOK, "hello %s ,your path is %s\n", c.Query("name"), c.Path)
	})
	//curl http://localhost:9999/hello/Euraxluo
	w.GET("/hello/:name", func(c *web.Context) {
		c.STRING(http.StatusOK, "hello %s ,your path is %s\n", c.Param("name"), c.Path)
	})
	//curl "http://localhost:9999/login" -X POST -d 'username=Euraxluo&password=1234'
	w.POST("/login", func(c *web.Context) {
		c.JSON(http.StatusOK, web.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	//curl http://localhost:9999/assets/js/main.js
	w.GET("/assets/:filepath/:file", func(c *web.Context) {
		c.JSON(http.StatusOK, web.H{
			"filepath": c.Param("filepath"),
			"file": c.Param("file"),
		})
	})
	w.Run(":9999")
}
```

### 特性3.　实现了路由分组控制
分组控制(Group Control)是 Web 框架应提供的基础功能之一。
所谓分组，是指路由的分组。如果没有路由分组，我们需要针对每一个路由进行控制。
大部分情况下的路由分组都是根据前缀进行区分的，因此我们的分组控制也是根据路由的相同前缀，进行分组。
并且分组应该支持嵌套，也即分组下面还可以进行分组。
```go
type (
	RouterGroup struct {
		prefix           string
		middlewares      []HandlerFunc //该分组的middleware
		parent           *RouterGroup
		web              *Web               //all groups share a Web instance
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
//添加路由
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.web.router.addRoute(method, pattern, handler)
}
```
通过addRoute函数，我们调用的group.web.router.addRoute,也即将该路由添加到了该分组对象所管理的Web对象中

使用demo 
```go

func main() {
	w := web.New()
	//curl http://localhost:9999/
	w.GET("/", func(c *web.Context) {
		c.HTML(http.StatusOK, "<h1>Hello World</h1>")
	})
	//curl http://localhost:9999/hello?name=Euraxluo
	w.GET("/hello", func(c *web.Context) {
		c.STRING(http.StatusOK, "hello %s ,your path is %s\n", c.Query("name"), c.Path)
	})
	//curl http://localhost:9999/hello/Euraxluo
	w.GET("/hello/:name", func(c *web.Context) {
		c.STRING(http.StatusOK, "hello %s ,your path is %s\n", c.Param("name"), c.Path)
	})
	//curl "http://localhost:9999/login" -X POST -d 'username=Euraxluo&password=1234'
	w.POST("/login", func(c *web.Context) {
		c.JSON(http.StatusOK, web.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	//curl http://localhost:9999/assets/js/main.js
	w.GET("/assets/:filepath/:file", func(c *web.Context) {
		c.JSON(http.StatusOK, web.H{
			"filepath": c.Param("filepath"),
			"file": c.Param("file"),
		})
	})

	v1 := w.Group("/v1")
	{

		//curl http://localhost:9999/v1/
		v1.GET("/", func(c *web.Context) {
			c.HTML(http.StatusOK, "<h1>Hello World</h1>")
		})
		//curl http://localhost:9999/v1/hello?name=Euraxluo
		v1.GET("/hello", func(c *web.Context) {
			c.STRING(http.StatusOK, "hello %s ,your path is %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := w.Group("/v2")
	{
		//curl http://localhost:9999/v2/hello/Euraxluo
		v2.GET("/hello/:name", func(c *web.Context) {
			c.STRING(http.StatusOK, "hello %s ,your path is %s\n", c.Param("name"), c.Path)
		})
		//curl "http://localhost:9999/v2/login" -X POST -d 'username=Euraxluo&password=1234'
		v2.POST("/login", func(c *web.Context) {
			c.JSON(http.StatusOK, web.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}
	w.Run(":9999")
}
```

### 特性4. 实现了中间件机制
接收到请求后，查找所有应作用于该路由的中间件，保存在Context中，依次进行调用。
我们的中间件设计中，中间件不仅作用在处理流程前，也可以作用在处理流程后，即在用户定义的 Handler 处理完毕后，还可以执行剩下的操作。
```go
type Context struct {
	//origin obj
	Writer     http.ResponseWriter
	Req        *http.Request
	//request info
	Path       string
	Method     string
	Params map[string]string
	//response info
	StatusCode int
	//middleware
	handlers []HandlerFunc
	handlerIndex int
	//group ptr
	group *RouterGroup
}
func (c *Context)Handle(){
	c.handlerIndex++
	s:=len(c.handlers)
	for ;c.handlerIndex<s;c.handlerIndex++ {
		c.handlers[c.handlerIndex](c)
	}
}
```

###特性5.　template
基于Go语言内置的html/template模板标准库，实现了支持普通变量渲染、列表渲染、对象渲染等功能的Template功能
```go
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
//添加处理静态资源的handler
func (group *RouterGroup)Static(relativePath string,root string){
	//创建一个handler
	staticHandler :=createStaticHandler(group,relativePath,http.Dir(root))
	urlPattern:=path.Join(relativePath,"/*filepath")
	group.GET(urlPattern,staticHandler)
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
```
```go
//TEMPLATE响应
func (c *Context) TEMPLATE(code int, htmlName string,data interface{}) {
	c.SetHeader("Context-Type", "text/html")
	c.Status(code)
	group :=c.group
	templateRender(group,c,htmlName,data,nil)
}
//递归render
func templateRender(group *RouterGroup,c *Context,htmlName string,data interface{},err error){
	if group == nil{
		c.Fail(http.StatusInternalServerError,err.Error())
		return
	}
	if group.htmlTemplates == nil{
		templateRender(group.parent,c,htmlName,data,err)
	}else{
		if renderError:=group.htmlTemplates.ExecuteTemplate(c.Writer,htmlName,data);renderError != nil{
			if err!=nil{
				renderError=err
			}
			templateRender(group.parent,c,htmlName,data,renderError)
		}
	}
}
```
