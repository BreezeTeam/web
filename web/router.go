package web

import (
	"net/http"
	"strings"
)

type router struct {
	roots map[string]*node
	handlers map[string]HandlerFunc
}

//init handler
func newRouter() *router {
	return &router{
		roots: make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}
func parsePattern(pattern string) []string  {
	patternSplit:=strings.Split(pattern,"/")
	parts:=make([]string, 0)
	for _,item := range patternSplit{
		if item != ""{
			parts= append(parts,item)
			if item[0] =='*'{
				break
			}
		}
	}
	return parts
}

//add router
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	//设置handlers
	key := method + "-" + pattern
	r.handlers[key] = handler

	//将pattern解析后添加到对应请求的路径中
	parts:=parsePattern(pattern)
	//如果该请求类型的根节点为空那么初始化
	_,ok:=r.roots[method]
	if !ok{
		r.roots[method]=&node{}
	}
	r.roots[method].insert(pattern,parts,0)
}


func (r *router) getRoute(method string, path string)(*node,map[string]string){
	//看这个方法被不被支持
	root,ok:=r.roots[method]
	if !ok{
		return nil,nil
	}
	//将路径解析为parts，并找到对应的节点
	pathParts:=parsePattern(path)
	pathNode:=root.search(pathParts,0)

	//解析两种匹配符的参数
	params:=make(map[string]string)
	if pathNode !=nil{
		nodeParts:=parsePattern(pathNode.pattern)
		for index, nodePart := range nodeParts {
			if nodePart[0] == ':'{
				params[nodePart[1:]] = pathParts[index]
			}
			if nodePart[0] == '*'&&len(nodePart)>1{
				params[nodePart[1:]] = strings.Join(pathParts[index:],"/")
				break
			}
		}
		return pathNode,params
	}
	return nil,nil
}
func (r *router) handle(c *Context) {
	pathNode,params:=r.getRoute(c.Method,c.Path)
	if pathNode !=nil{
		c.Params = params
		key := c.Method + "-" + pathNode.pattern
		r.handlers[key](c)
	}else {
		c.STRING(http.StatusNotFound, "404 Not Found: %s\n", c.Path)
	}
}
