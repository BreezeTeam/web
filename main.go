package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	"web"
)

func main() {
	w := web.New()
	//设置funcMap
	{
		w.SetFuncMap(template.FuncMap{
			"Time2String": Time2String,
		})
		//加载模板文件
		w.LoadTemplate("templates2/*")
		//将该目录设置为静态资源目录
		w.Static("/assets2", "/home/yons/go/src/web/assets2")
		//add middleware
		w.Use(web.Logger(),web.Recovery())

		w.GET("/panic", func(c *web.Context) {
			names := []string{"geektutu"}
			c.STRING(http.StatusOK, names[100])
		})

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

		//w这个组将assets2添加到了资源资源服务器中，但是assets没有
		//curl http://localhost:9999/assets/js/main.js
		w.GET("/assets/:filepath/:file", func(c *web.Context) {
			c.JSON(http.StatusOK, web.H{
				"filepath": c.Param("filepath"),
				"file":     c.Param("file"),
			})
		})
		//curl http://localhost:9999/assets2/js/main.js
		w.GET("/assets2/:filepath/:file", func(c *web.Context) {
			c.JSON(http.StatusOK, web.H{
				"filepath": c.Param("filepath"),
				"file":     c.Param("file"),
			})
		})

		//curl http://localhost:9999/date/dateShow
		w.GET("/date/:title", func(c *web.Context) {
			c.TEMPLATE(http.StatusOK,
				"timeShow.template",
				web.H{
					"title": c.Param("title"),
					"now":   time.Now(),
				})
		})

	}

	v1 := w.Group("/v1")
	{
		//w组将templates2加载了，但是没有加载templates，因此拿不到css.template
		//curl http://localhost:9999/v1/css
		v1.GET("/css", func(c *web.Context) {
			c.TEMPLATE(http.StatusOK, "css.template", nil)
		})

		//curl http://localhost:9999/v1/date/dateShow
		v1.GET("/date/:title", func(c *web.Context) {
			c.TEMPLATE(http.StatusOK,
				"timeShow.template",
				web.H{
					"title": c.Param("title"),
					"now":   time.Now(),
				})
		})

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
		//设置中间件
		v2.Use(v2handler2(), v2handler())
		//设置funcMap
		v2.SetFuncMap(template.FuncMap{
			"IntDouble": IntDouble,
		})
		//加载模板文件
		v2.LoadTemplate("templates/*")
		//加载静态资源
		v2.Static("/assets", "/home/yons/go/src/web/assets")
		//curl http://localhost:9999/v2/css
		v2.GET("/css", func(c *web.Context) {
			c.TEMPLATE(http.StatusOK, "css.template", nil)
		})

		//curl http://localhost:9999/v2/students/Show
		v2.GET("/students/:title", func(c *web.Context) {
			type student struct {
				Name string
				Age  int
			}
			stu1 := &student{Name: "Euraxluo", Age: 20}
			stu2 := &student{Name: "Test", Age: 22}
			c.TEMPLATE(http.StatusOK,
				"arr.template",
				web.H{
					"title":          c.Param("title"),
					"stuArr":         [2]*student{stu1, stu2},
					"testIntDoulble": 11,
				})
		})

		//curl http://localhost:9999/v2/date/dateShow
		v2.GET("/date/:title", func(c *web.Context) {
			c.TEMPLATE(http.StatusOK,
				"timeShow.template",
				web.H{
					"title": c.Param("title"),
					"now":   time.Now(),
				})
		})
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

func Time2String(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func IntDouble(i int) string {
	return fmt.Sprintf("%d", i*2)
}

func v2handler() web.HandlerFunc {
	return func(c *web.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Handle()
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func v2handler2() web.HandlerFunc {
	return func(c *web.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Handle()
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v22", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
