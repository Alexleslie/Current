package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// 使用defer挂载错误恢复函数，调用recover()捕获panic，并将堆栈信息打印在日志中，向用户返回Internal Server Error
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.ReturnFail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}

// 打印debug的堆栈
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) //为了日志简洁一点，跳过前三个Caller
	//Callers 用来返回调用栈的程序计数器, 第 0 个 Caller 是 Callers 本身，第 1 个是上一层 trace，第 2 个是再上一层的 defer func。

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)   //获取对应的函数
		file, line := fn.FileLine(pc) //获取该函数的文件名和行号，打印
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
