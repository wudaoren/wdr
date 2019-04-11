package web

import (
	"net/http/httputil"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

//异常恢复中间件
func Recover(fn func(c *gin.Context, errs ...interface{})) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				fn(c, err, string(httprequest), string(stack))
			}
		}()
		c.Next()
	}
}
