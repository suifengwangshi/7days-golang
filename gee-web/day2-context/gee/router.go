package gee

import (
	"log"
	"net/http"
)

// 将路由部分从gee.go中提取出来，方便之后对路由的功能进行增强，例如支持动态路由

// router implement the interface of ServeHTTP
type router struct {
	handlers map[string]HandlerFunc
}

// newRouter is the constructor of gee.router
func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

// 加路由信息到路由哈希表中
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// 根据Context中封装的信息构造出key，调用对应的处理函数
func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
