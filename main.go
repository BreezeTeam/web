package main

import (
	"log"
	"time"
	"web"
	"net/http"
)

func main() {
	w := web.New()

	//add middleware
	w.Use(web.Logger())


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
	v1.Use()
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
	v2.Use(v2handler2(),v2handler())
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

func v2handler() web.HandlerFunc{
	return func(c *web.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Handle()
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func v2handler2() web.HandlerFunc{
	return func(c *web.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Handle()
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v22", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}