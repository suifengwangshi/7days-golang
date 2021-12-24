package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{} // 起了一个别名gee.H，构建JSON数据时显得更简洁

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string // 存储根据路由解析出来的参数信息
	// response info
	StatusCode int
}

// NewContext is the constructor of gee.Context
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

// 在 HandlerFunc 中，希望能够访问到解析的参数，因此，需要对 Context 对象增加一个属性和方法，
// 来提供对路由参数的访问。我们将解析后的参数存储到Params中，通过c.Param("lang")的方式获取到对应的值。
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// 提供查询PostForm方法  可以额获取 url 中? 后面的请求参数
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// 提供查询Query方法
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 设置响应信息的状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// 设置响应头信息
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 必须先 Header().Set 再 WriteHeader() 最后再 Wrute

// 快速返回String类型的HTTP响应
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// 快速返回JSON类型的HTTP响应
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer) // NewEncoder returns a new encoder that writes to c.Writer
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// 快速返回Data类型的HTTP响应
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// 快速返回HTML类型的HTTP响应
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))

}
