package stcp

import (
	"net"
	"time"
)

//客户端配置
type ClientConfig struct {
	ServerAddr   string //服务器端地址
	Heartbeat    []byte //自动心跳数据
	Deadline     int64  //连接超时时间,单位：秒
	RequestFirst []byte //请求头
	Handler      ClientHandler
}

//socket客户端
type Client struct {
	Conn    *Conn //继承连接对象
	handler ClientHandler
	config  *ClientConfig
	closed  chan bool
}

//创建一个客户端
func NewClient(config *ClientConfig) *Client {
	p := new(Client)
	p.config = config
	p.handler = config.Handler
	p.closed = make(chan bool)
	return p
}

//发送数据
func (this *Client) Write(bt []byte) (int, error) {
	return this.Conn.Write(bt)
}

//关闭连接
func (this *Client) Close() {
	this.closed <- true
	this.Conn.Close()
}

//
func (this *Client) checktimeout() {
	defer recover()
	if this.config.Deadline <= 0 {
		return
	}
	wait := time.Second * time.Duration(this.config.Deadline)
	for {
		select {
		case <-time.After(wait):
			now := time.Now().Unix()
			if (this.Conn.LastAccessTime.Unix() + this.config.Deadline) < now {
				this.Write(this.config.Heartbeat)
			}
		case <-this.closed:
			return
		}
	}
}

//发起连接
//输入：addr=发起请求的地址
func (this *Client) Connect() error {
	var err error
	var data []byte
	conn, e1 := net.Dial("tcp", this.config.ServerAddr)
	if e1 != nil { //连接失败
		return e1
	}
	this.Conn = newConn(conn)
	defer func() {
		context := newCloseContext(this.Conn, err, recover())
		if context.Error() != nil {
			this.handler.OnError(context)
		}
		this.handler.OnClose(context)
		this.Close()
	}()
	this.Write(this.config.RequestFirst)
	//if data, err = this.Conn.Read(); err != nil {
	//return err
	//}
	context := newContext(this.Conn, data, err)
	this.handler.OnConnect(context)
	go this.checktimeout()
	for {
		if ndata, err := this.Conn.Read(); err != nil {
			break
		} else {
			goTry(func() {
				this.Conn.LastAccessTime = time.Now() //更新链接时间
				context := newContext(this.Conn, ndata, nil)
				if this.isHeartbeat(ndata) {
					this.handler.OnHeartbeat(context)
				} else {
					this.Conn.triggerReadCall(ndata)
					this.handler.OnMessage(context)
				}
			}, func(e Error) {
				context := newContext(this.Conn, nil, e)
				this.handler.OnError(context)
			})
		}
	}
	return nil
}

func (this *Client) isHeartbeat(data []byte) bool {
	return len(data) == len(this.config.Heartbeat) && string(data) == string(this.config.Heartbeat)
}
