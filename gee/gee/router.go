package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

//roots key eg,roots['GET'] roots['POST']
//handlers key eg,handlers['GET-/p/:lang/doc'],handlers['POST=/p/book']

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Only one * is allowed
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

func (r *router) addRouter(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRouter(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

// 将从路由匹配到的Handle添加到c.handlers列表中，执行c.Next()从而一一调用中间件和请求再逐一从后往前执行Next()之后的操作
func (r *router) handle(c *Context) {
	n, param := r.getRouter(c.Method, c.Path)
	log.Printf("[handle]n=%v", n)
	log.Printf("[handle]param=%v", param)

	if n != nil {
		c.Params = param
		key := c.Method + "-" + n.pattern
		log.Printf("[handle]key=%v", key)
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.WriteString(http.StatusNotFound, "404 NOT FOUND:%s\n", c.Path)
		})
	}
	c.Next()
}
