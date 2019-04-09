package web

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Sequence struct {
	gced      bool
	requested *sync.Map
}

//防止重复提交的中间件
func SequenceHandler(gcTime int64, onrepeat func(*gin.Context)) gin.HandlerFunc {
	var seq = new(Sequence)
	seq.requested = new(sync.Map)
	seq.gc(gcTime)
	return func(c *gin.Context) {
		defer c.Next()
		key := fmt.Sprint(c.Request.Method, c.Request.RemoteAddr, c.Request.URL, c.Request.UserAgent(), c.Request.ContentLength)
		if _, ok := seq.requested.Load(key); ok {
			onrepeat(c)
			c.Abort()
		}
		seq.requested.Store(key, time.Now().Unix())
	}
}

//当请求的时间超过gc时间则清除该请求
func (this *Sequence) gc(wait int64) {
	if this.gced {
		return
	}
	this.gced = true
	go func() {
		defer recover()
		for {
			time.Sleep(time.Second * time.Duration(wait))
			now := time.Now().Unix()
			this.requested.Range(func(k, v interface{}) bool {
				if visitedTime, ok := v.(int64); ok && now-visitedTime > wait {
					this.requested.Delete(k)
				}
				return true
			})
		}
	}()
}
