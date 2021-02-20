package web

import "net/http"

type router struct {
	handlers map[string]HandlerFunc
}

//init handler
func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

//add router
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.STRING(http.StatusNotFound, "404 Not Found: %s\n", c.Path)
	}
}
