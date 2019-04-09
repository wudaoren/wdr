package stcp

//处理socket服务器端

import (
	"net"
	"sync"
	"sync/atomic"
	"time"
)

//服务器端配置
type ServerConfig struct {
	ListenAddr string //服务器端监听地址
	MaxSize    int64  //最大允许连接数
	Deadline   int64  //连接超时时间,单位：秒
	Heartbeat  []byte //自动心跳数据
	Handler    ServerHandler
}

//tcp服务器
type Server struct {
	conns         *sync.Map
	config        *ServerConfig
	handler       ServerHandler
	connectLength int64 //连接数量
	index         int64 //连接索引
	stop          chan bool
	listener      net.Listener
}

//创建一个server,并初始化连接池
func NewServer(config *ServerConfig) *Server {
	p := new(Server)
	p.conns = new(sync.Map)
	p.config = config
	p.handler = config.Handler
	p.stop = make(chan bool)
	return p
}

//通过连接编号获取某个连接对象
func (this *Server) GetConn(id int64) (*Conn, bool) {
	if conn, ok := this.conns.Load(id); ok {
		return conn.(*Conn), true
	}
	return nil, false
}

//连接数量
func (this *Server) ConnLength() int64 {
	return this.connectLength
}

//遍历所有连接对象
func (this *Server) EachConn(call func(c *Conn) bool) {
	this.conns.Range(func(id, conn interface{}) bool {
		return call(conn.(*Conn))
	})
}

//发起监听
func (this *Server) Run() error {
	listener, err := net.Listen("tcp", this.config.ListenAddr)
	if err != nil {
		return err
	}
	this.listener = listener
	go this.checktimeout()
	for {
		//等待客户端接入
		conn, err := listener.Accept()
		if err != nil {
			continue
		} else {
			go this.handlerConnect(conn)
		}
	}
	return nil
}

//停止服务
func (this *Server) Stop() {
	this.listener.Close()
	this.stop <- true
}

//处理每个连接
func (this *Server) handlerConnect(c net.Conn) {
	if this.ConnLength() > this.config.MaxSize {
		c.Close()
		return
	}
	var err error
	var data []byte
	conn := newConn(c)
	atomic.AddInt64(&this.connectLength, 1)
	conn.id = this.getNewIndex() //使用连接池长度作为地址编号
	defer func() {
		context := newCloseContext(conn, err, recover())
		if context.Error() != nil {
			this.handler.OnError(context)
		}
		conn.Close()
		this.handler.OnClose(context)
		atomic.AddInt64(&this.connectLength, -1)
		this.conns.Delete(conn.id) //删除连接池的数据
	}()
	if data, err = conn.Read(); err != nil { //初次连接发送的头数据
		return
	}
	context := newContext(conn, data, nil)
	if !this.handler.OnConnect(context) { //关闭验证失败的连接
		return
	}
	this.conns.Store(conn.id, conn)
	for {
		if ndata, err := conn.Read(); err != nil {
			return
		} else {
			goTry(func() {
				context := newContext(conn, ndata, nil)
				if this.isHeartbeat(ndata) {
					conn.Write(ndata) //原样返回心跳数据
					this.handler.OnHeartbeat(context)
				} else {
					conn.triggerReadCall(ndata)
					this.handler.OnMessage(context)
				}
			}, func(e Error) {
				context := newContext(conn, nil, e)
				this.handler.OnError(context)
			})
		}
	}
}

//
func (this *Server) isHeartbeat(data []byte) bool {
	return len(data) == len(this.config.Heartbeat) && string(data) == string(this.config.Heartbeat)
}

//检查连接是否超时
func (this *Server) checktimeout() {
	defer recover()
	if this.config.Deadline <= 0 {
		return
	}
	wait := time.Second * time.Duration(this.config.Deadline)
	for {
		select {
		case <-time.After(wait):
			now := time.Now().Unix()
			this.EachConn(func(conn *Conn) bool {
				if (conn.LastAccessTime.Unix() + this.config.Deadline) < now {
					this.conns.Delete(conn.id)
					conn.Close()
				}
				return true
			})
		case <-this.stop:
			return
		}
	}
}

//
func (this *Server) getNewIndex() int64 {
	atomic.AddInt64(&this.index, 1)
	return this.index
}
