package main

import (
	"net/http"
	"web"
)

func main() {
	w := web.New()
	w.GET("/", func(c *web.Context) {
		c.HTML(http.StatusOK, "<h1>Hello World</h1>")
	})
	///hello?name=Euraxluo
	w.GET("/hello", func(c *web.Context) {
		c.STRING(http.StatusOK, "hello %s ,your path is %s\n", c.Query("name"), c.Path)
	})
	w.POST("/login", func(c *web.Context) {
		c.JSON(http.StatusOK, web.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	w.Run(":9999")
}
