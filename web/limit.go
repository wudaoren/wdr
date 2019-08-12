package web

import (
	"github.com/gin-gonic/gin"
)

//最大允许连接数
func MaxAllowed(n int) gin.HandlerFunc {
	type s struct{}
	sem := make(chan s, n)
	acquire := func() {
		sem <- s{}
	}
	release := func() {
		<-sem
	}
	return func(c *gin.Context) {
		acquire()       //before request
		defer release() //after request
		c.Next()
	}
}
