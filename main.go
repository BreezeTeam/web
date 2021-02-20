package main

import (
	"web"
	"net/http"
)

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
