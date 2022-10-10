package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

// Context 封装了HTML/String/JSON函数，能快速根据请求构造HTTP响应
// Context 随着每一个请求的出现而产生，请求的结束而销毁，和当前请求强相关的信息都应由 Context 承载（路由的处理函数、中间件处理、参数等）。
type Context struct {
	//origin objects
	Writer http.ResponseWriter //响应
	Req    *http.Request       //请求
	//request info，请求信息，路径、方法和参数
	Path   string            //请求路径
	Method string            //请求方法
	Params map[string]string //请求参数

	//response info，响应信息，状态码
	StatusCode int //响应状态码

	//middleware，中间件
	handlers []HandlerFunc
	index    int //记录当前执行到第几个中间件

	//engine pointer，可以通过Context访问Engine中的HTML模板
	engine *Engine
}

// newContext 定义了新建一个Context的方法
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Path:   req.URL.Path,
		Method: req.Method,
		Req:    req,
		Writer: w,
		index:  -1,
	}
}

// Next 定义了顺序调用中间件的方法
// Next 中间件调权用时，控制交给下一个中间件（调用下一个handler），直到调用到最后一个中间件，然后在从后往前，调用每个中间件在Next方法之后定义的部分
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

// GetParam 定义了获取Context.Params的方法
func (ctx *Context) GetParam(key string) string {
	value, _ := ctx.Params[key]
	return value
}

// GetPostForm 定义了获取Context.Request请求的表单key对应信息方法
func (ctx *Context) GetPostForm(key string) string {
	return ctx.Req.FormValue(key)
}

// Query 定义了Context.Request.URL查询方法
func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get(key)
}

// SetStatus 定义了设置Context状态码的方法
func (ctx *Context) SetStatus(code int) {
	ctx.StatusCode = code
	ctx.Writer.WriteHeader(code)
}

// SetHeader 定义了设置Context响应头的方法
func (ctx *Context) SetHeader(key string, value string) {
	ctx.Writer.Header().Set(key, value)
}

// WriteString 定义了快速构造HTTP响应为String格式的方法
func (ctx *Context) WriteString(code int, format string, values ...interface{}) {
	ctx.SetHeader("Content-Type", "text/plain")
	ctx.SetStatus(code)
	_, err := ctx.Writer.Write([]byte(fmt.Sprintf(format, values...)))
	if err != nil {
		ctx.ReturnFail(500, err.Error())
	}
}

// WriteJSON 定义了快速构造HTTP响应为JSON格式的方法
func (ctx *Context) WriteJSON(code int, obj interface{}) {
	ctx.SetHeader("Content-Type", "application/json")
	ctx.SetStatus(code)

	// writer作为缓冲区传入到json的NewEncoder方法中，使得Encode方法直接将obj编码后写入到Writer里
	encoder := json.NewEncoder(ctx.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(ctx.Writer, err.Error(), code)
	}
}

// WriteJSON 定义了快速构造HTTP响应为Bytes格式的方法
func (ctx *Context) WriteBytes(code int, data []byte) {
	ctx.SetStatus(code)
	_, err := ctx.Writer.Write(data)
	if err != nil {
		ctx.ReturnFail(500, err.Error())
	}
}

// WriteJSON 定义了快速构造HTTP响应为HTML格式的方法
func (ctx *Context) WriteHTML(code int, fileName string, data interface{}) {
	ctx.SetHeader("Content-Type", "text/html")
	ctx.SetStatus(code)
	if err := ctx.engine.htmlTemplates.ExecuteTemplate(ctx.Writer, fileName, data); err != nil {
		ctx.ReturnFail(500, err.Error())
	}
}
