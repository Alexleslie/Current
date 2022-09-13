package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	//origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	//request info
	Path   string
	Method string
	Params map[string]string

	//response info
	StatusCode int

	//middleware
	handlers []HandlerFunc
	index    int //记录当前执行到第几个中间件

	//engine pointer，可以通过Context访问Engine中的HTML模板
	engine *Engine
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Path:   req.URL.Path,
		Method: req.Method,
		Req:    req,
		Writer: w,
		index:  -1,
	}
}

// Next 中间件调权用时，控制交给下一个中间件（调用下一个handler），直到调用到最后一个中间件
// 然后在从后往前，调用每个中间件在Next方法之后定义的部分
func (ctx *Context) Next() {
	ctx.index++
	s := len(ctx.handlers)
	if ctx.index < s {
		ctx.handlers[ctx.index](ctx)
	}
}

func (ctx *Context) ReturnFail(code int, err string) {
	ctx.index = len(ctx.handlers)
	ctx.WriteJSON(code, H{"message": err})
}

func (ctx *Context) GetParam(key string) string {
	value, _ := ctx.Params[key]
	return value
}

func (ctx *Context) GetPostForm(key string) string {
	return ctx.Req.FormValue(key)
}

func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get(key)
}

func (ctx *Context) SetStatus(code int) {
	ctx.StatusCode = code
	ctx.Writer.WriteHeader(code)
}

func (ctx *Context) SetHeader(key string, value string) {
	ctx.Writer.Header().Set(key, value)
}

func (ctx *Context) WriteString(code int, format string, values ...interface{}) {
	ctx.SetHeader("Content-Type", "text/plain")
	ctx.SetStatus(code)
	_, err := ctx.Writer.Write([]byte(fmt.Sprintf(format, values...)))
	if err != nil {
		ctx.ReturnFail(500, err.Error())
	}
}

func (ctx *Context) WriteJSON(code int, obj interface{}) {
	ctx.SetHeader("Content-Type", "application/json")
	ctx.SetStatus(code)

	// writer作为缓冲区传入到json的NewEncoder方法中，使得Encode方法直接将obj编码后写入到Writer里
	encoder := json.NewEncoder(ctx.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(ctx.Writer, err.Error(), code)
	}
}

func (ctx *Context) WriteBytes(code int, data []byte) {
	ctx.SetStatus(code)
	_, err := ctx.Writer.Write(data)
	if err != nil {
		ctx.ReturnFail(500, err.Error())
	}
}

func (ctx *Context) WriteHTML(code int, fileName string, data interface{}) {
	ctx.SetHeader("Content-Type", "text/html")
	ctx.SetStatus(code)
	if err := ctx.engine.htmlTemplates.ExecuteTemplate(ctx.Writer, fileName, data); err != nil {
		ctx.ReturnFail(500, err.Error())
	}
}
