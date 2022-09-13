package gee

import (
	"html/template"
	"net/http"
	"path"
	"strings"
)

// HandlerFunc 处理请求的函数类型
type HandlerFunc func(*Context)

// RouterGroup 当前路径所包含结构的组
type RouterGroup struct {
	prefix      string        // 前缀
	middlewares []HandlerFunc // 中间件
	parent      *RouterGroup  // 父分组
	engine      *Engine       // 所有分组共享一个Engine实例
}

type Engine struct {
	*RouterGroup
	router        *router            // 全局路由
	groups        []*RouterGroup     // 全局RouterGroup
	htmlTemplates *template.Template //html 渲染，将所有的模板加载进内存
	funcMap       template.FuncMap   //html 渲染，所有自定义模板渲染函数
}

// Default 配置  Logger & Recovery 中间件
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

// 创建一个处理静态路径的handler函数
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(ctx *Context) {
		file := ctx.GetParam("filepath")
		//Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			ctx.SetStatus(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(ctx.Writer, ctx.Req)
	}
}

// Static 将磁盘上的静态文件夹{root}映射到某个静态文件路径{relativePath}上，
// 使得请求在访问该静态路径{relativePath}时都会转发到静态文件夹{root}上
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")

	group.GET(urlPattern, handler)
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}  //(*Engine).engine是指向自己
	engine.groups = []*RouterGroup{engine.RouterGroup} //engine的组指向包含自身路由组的路由组，Engine拥有所有RouterGroup
	return engine
}

func (engine *Engine) addRouter(method string, pattern string, handler HandlerFunc) {
	engine.router.addRouter(method, pattern, handler)
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRouter("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRouter("POST", pattern, handler)
}

// Run defined the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// Use 将中间件添加到路由组Group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// ServeHTTP 接口，接管了所有的 HTTP 请求
// 当接收到一个具体请求时，利用前缀判断请求适用于哪些中间件（前缀是静态的话，是从左到右的查找过程，顺序可以保证）
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	ctx := newContext(w, req)
	ctx.handlers = middlewares
	ctx.engine = engine
	engine.router.addHandleByPath(ctx)
	ctx.Next()
}

// Group 创建一个新的子路由分组，属于同个族系的分组共享一个engine实例
func (group *RouterGroup) newChildGroup(prefix string) *RouterGroup {
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: group.engine,
	}

	group.engine.groups = append(group.engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRouter(method string, component string, handler HandlerFunc) {
	pattern := group.prefix + component
	group.engine.router.addRouter(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRouter("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRouter("POST", pattern, handler)
}
