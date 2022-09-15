package gee

import (
	"log"
	"time"
)

// Logger 全局中间件，在所有中间件被调用后，Next方法之后定义的部分可以被从后往前被调用
func Logger() HandlerFunc {
	return func(ctx *Context) {
		//Start timer
		t := time.Now()
		//Process request 等待执行其他中间件或者用户的Handle
		ctx.Next()
		//Calculate resolution time 计算时间
		log.Printf("[HandleFunc][%d]%s in %v", ctx.StatusCode, ctx.Req.RequestURI, time.Since(t))
	}
}
