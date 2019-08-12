package web

import (
	"fmt"
	"sync"
	"wdr/utils"

	"github.com/gin-gonic/gin"
)

//签名中间价
func SignHandler(tokenKey, signKey, secret string) gin.HandlerFunc {
	signsMap := new(sync.Map)
	length := 0
	return func(c *gin.Context) {
		length++
		if length > 100000 {
			signsMap = new(sync.Map)
			length = 0
		}
		token := c.GetHeader(tokenKey)
		sign := c.GetHeader(signKey)
		if _, ok := signsMap.Load(sign); ok { //签名重复使用
			c.AbortWithStatus(403)
			return
		}
		if utils.Sha256(fmt.Sprint(token, secret)) != sign { //签名未通过
			c.AbortWithStatus(403)
			return
		}
		signsMap.Store(sign, true)
	}
}
