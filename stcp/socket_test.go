package stcp

import (
	"log"
	"testing"
	"time"
)

func TestSocket(t *testing.T) {
	go server()
	go client()
	time.Sleep(time.Hour * 24)
}

//设置心跳内容
var heartbeat = []byte("@#")
var addr = "localhost:9010"

type ServerHandle struct {
}

//
func (this *ServerHandle) OnConnect(c *Context) bool {
	log.Println("服务区收到连接")
	c.Write([]byte("ok"))
	return true
}

//
func (this *ServerHandle) OnMessage(c *Context) {
	log.Println("服务器段收到消息", string(c.Data()))
}

//
func (this *ServerHandle) OnClose(c *Context) {
	log.Println("服务器关闭连接", c.Id())
}

//
func (this *ServerHandle) OnError(c *Context) {
	log.Println("服务器遇到错误", c.Error())
}

//
func (this *ServerHandle) OnHeartbeat(c *Context) {
	log.Println("服务器端收到心跳", string(c.Data()))
}

//
func server() {
	sev := NewServer(&ServerConfig{
		ListenAddr: addr,
		Heartbeat:  heartbeat,
		Deadline:   5,
		Handler:    new(ServerHandle),
	})
	sev.Run()
}

type ClientHandle struct {
}

//
var index = 0

func (this *ClientHandle) OnConnect(c *Context) {
	log.Println("客户端连接成功", string(c.Data()))
	go func() {
		index++
		if index == 1 {
			time.AfterFunc(time.Second*10, func() {
				c.Close()
			})
			return
		}
		var i int
		for {
			i++
			time.Sleep(time.Second * 1)
			if i > 10 {
				return
			}
			c.Write([]byte("发送1111"))
		}
	}()
}

//
func (this *ClientHandle) OnMessage(c *Context) {
	log.Println("客户端收到消息", string(c.Data()))
}

//
func (this *ClientHandle) OnClose(c *Context) {
	log.Println("客户端关闭连接")
	go client()
}

//
func (this *ClientHandle) OnError(c *Context) {
	log.Println("客户端遇到错误", c.Error())
}

//
func (this *ClientHandle) OnHeartbeat(c *Context) {
	log.Println("客户端收到心跳", string(c.Data()))
}

//
func client() {
	cli := NewClient(&ClientConfig{
		ServerAddr:   addr,
		Heartbeat:    heartbeat,
		RequestFirst: []byte("#¥ADFEJSDF"),
		Deadline:     2,
		Handler:      new(ClientHandle),
	})
	cli.Connect()
}
