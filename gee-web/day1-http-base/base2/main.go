// http.ListenAndServe  ListenAndServe func(addr string, handler Handler) error
// 第一个参数是地址，设置监听端口，第二个参数代表处理所有的http请求的实例，nil代表使用标准库中的实例处理，
// 可以自己自定义数据实现 Handler 接口，这样第二个参数就可以使用自定义的实例处理http请求
/*
package http
type Handler interface {
	ServeHTTP(w ResponseWriter, r *Request)
}
func ListenAndServe(address string, h Handler) error
*/
// http.Request 包含了该http请求的所有信息，比如请求地址、请求头Header和请求体Body等
// http.ResponseWriter 可以构造针对该请求的响应

package main

import (
	"fmt"
	"log"
	"net/http"
)

// Engine is the defined handler for all request
type Engine struct{} // 实现Engine之后i，拦截了所有的HTTP请求，拥有了统一的控制入口，可以自由定义路由映射的规则，也可以统一添加一些处理逻辑，例如日志、异常处理等

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	case "/hello":
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL)
	}
}

func main() {
	engine := new(Engine)
	log.Fatal(http.ListenAndServe(":9999", engine))
}
