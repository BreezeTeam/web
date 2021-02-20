package web

import (
	"log"
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
