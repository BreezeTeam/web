package web

import (
	"log"
	"net/http"
	"path"
	"time"
)

/**
 * @Description: 默认中间件
 * @return HandlerFunc
 */
func Logger() HandlerFunc {
	return func(c *Context){
		t:=time.Now()
		c.Handle()
		log.Printf("[%d] %s in %v",c.StatusCode,c.Req.RequestURI,time.Since(t))
	}
}

//添加静态资源处理函数
func createStaticHandler(group *RouterGroup,relativePath string,fs http.FileSystem) HandlerFunc{
	//找到绝对地址
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath,http.FileServer(fs))
	return func(c *Context) {
		file:=c.Param("filepath")
		if _,err:=fs.Open(file);err != nil{
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer,c.Req)
	}
}
