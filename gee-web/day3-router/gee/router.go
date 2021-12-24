/*
在 trie.go 文件中实现了 Trie树的插入和查找，该文件是将Trie树应用到路由中去。我们使用roots来存储每种请求方式
的Trie树根节点，使用handlers存储每种请求方式的HandlerFunc.
getRoute函数中，还解析了 : 和 * 两种匹配符的参数，返回一个哈希表map，比如/p/go/doc匹配到/p/:lang/doc，解析
结果为{lang:"go"}，/static/css/geektutu.css匹配到/static/*filepath，解析结果为{filepath:"css/geektutu.css"}
*/

package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Only one * is allowed 将路由分解成各个部分组成的parts数组
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/") // 将路由按照分割符分割开

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			// 遇到了通配符*
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 添加路由信息
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern) // 先将路由信息进行分割成[]string数组

	key := method + "-" + pattern // 哈希表的key
	_, ok := r.roots[method]      // 判断系统中是否有method请求方式(GET/POST等)的Trie树根节点
	// 没有的话就创建一个
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0) // 插入当前路由，默认根节点height=0
	r.handlers[key] = handler                 // 将路由对应的处理函数加上
}

// 根据路由获取对应的处理函数
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path) // 先分割路由成[]string
	params := make(map[string]string) // 哈希表
	root, ok := r.roots[method]       // 判断请求方法是否存在

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0) // Trie树上查询

	// Trie树上有对应的路由信息
	if n != nil {
		parts := parsePattern(n.pattern) // 分析查询到的节点的pttern
		for index, part := range parts {
			// 遇到了 :
			if part[0] == ':' {
				params[part[1:]] = searchParts[index] // /p/go/doc 匹配到 /p/:lang/doc 返回 {lang:"go"}
			}
			// 遇到了 *
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		// 在调用匹配到的handler前，将解析出来的路由参数赋值给了c.Params 这样就能够在handler中，
		// 通过Context对象访问到具体的值了。
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
