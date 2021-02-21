package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)
type H map[string]interface{}

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

//init
func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
		handlerIndex: -1,
	}
}
func (c *Context)Handle(){
	c.handlerIndex++
	s:=len(c.handlers)
	for ;c.handlerIndex<s;c.handlerIndex++ {
		c.handlers[c.handlerIndex](c)
	}
}
func (c *Context)Fail(code int,err string){
	c.handlerIndex = len(c.handlers)
	c.JSON(code,H{"message":err})
}

func (c *Context)Param(key string)string{
	value,_ := c.Params[key]
	return value
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

//String响应
func (c *Context) STRING(code int, format string, value ...interface{}) {
	c.SetHeader("Context-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, value...)))
}

//JSON 响应
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Context-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

//DATA响应
func (c *Context) DATA(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

//HTML响应
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Context-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

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