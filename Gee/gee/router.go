package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node       // roots['GET'],roots['POST']
	handlers map[string]HandlerFunc // handlers['GET-/p/:lang/doc'],handlers['POST-/p/book']
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 将URL以/拆分为各个部分，只允许一个通配符*存在
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 从根节点开始链接所有子路由节点
// 对该路由设置handler函数
func (r *router) addRouter(method string, absolutePattern string, handler HandlerFunc) {
	parts := parsePattern(absolutePattern)

	key := method + "-" + absolutePattern
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}

	// 链接所有以子路径节点
	r.roots[method].linkNode(absolutePattern, parts, 0)
	r.handlers[key] = handler
}

// path 为不包含通配符的具体请求路径
// pattern 为可能包含通配符的路径
func (r *router) getRouter(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	if tempNode := root.searchNode(searchParts, 0); tempNode != nil {
		params = matchParams(tempNode.pattern, path)
		return tempNode, params
	}
	return nil, nil
}

// 处理当路径模式存在通配符的情况
// 根据路径模式来匹配出请求路径对应的参数
func matchParams(pattern string, path string) map[string]string {
	params := map[string]string{}
	pathParts := parsePattern(path)
	patternParts := parsePattern(pattern)
	for index, part := range patternParts {
		// pattern=[/:lang/];path=[/cn/];param=[cn]
		if part[0] == ':' {
			params[part[1:]] = pathParts[index]
		}
		// 存在*通配符，则匹配从*号开始之后的所有路径
		// pattern=[/*filepath];path=[/cn/province];param=[cn/province]
		if part[0] == '*' && len(part) > 1 {
			params[part[1:]] = strings.Join(pathParts[index:], "/")
			break
		}
	}
	return params
}

// 将请求路径对应的handlerFunc添加到执行handlers列表里
func (r *router) addHandleByPath(c *Context) {
	pathNode, param := r.getRouter(c.Method, c.Path)
	if pathNode != nil {
		c.Params = param
		key := c.Method + "-" + pathNode.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.WriteString(http.StatusNotFound, "404 NOT FOUND:%s\n", c.Path)
		})
	}
}
