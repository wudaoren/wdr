package stcp

import (
	"fmt"
	"time"
)

//服务器端处理接口
type ServerHandler interface {
	OnConnect(*Context) bool
	OnMessage(*Context)
	OnClose(*Context)
	OnError(*Context)
	OnHeartbeat(*Context)
}

//客户端处理接口
type ClientHandler interface {
	OnConnect(*Context)
	OnMessage(*Context)
	OnClose(*Context)
	OnError(*Context)
	OnHeartbeat(*Context)
}

var (
	DEBUG = true
)

//调试输出
func debug(arg ...interface{}) {
	if DEBUG {
		now := time.Now().Format("[2006-01-02 15:04:05]")
		def := []interface{}{now}
		fmt.Println(append(def, arg...)...)
	}
}

//自定义模拟类型
type Error interface{}

//新建协程
func goTry(call func(), errcall func(Error)) {
	go func() {
		defer func() {
			if err := recover(); err != nil && errcall != nil {
				errcall(err)
			}
		}()
		call()
	}()
}
