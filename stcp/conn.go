package stcp

import (
	"net"
	"sync"
	"time"
)

//-----------------------------------------------------------------
//socket连接
type Conn struct {
	net.Conn                    //net默认的连接对象
	id             int64        //当前连接的编号
	LastAccessTime time.Time    //最后访问时间
	cache          *sync.Map    //缓存
	onReadCall     func([]byte) //当收到消息时
	closed         bool         //是否关闭
	sync.Mutex                  //小锁
}

func newConn(c net.Conn) *Conn {
	p := new(Conn)
	p.Conn = c //保存连接对象
	p.cache = new(sync.Map)
	p.LastAccessTime = time.Now()
	return p
}

//连接编号
func (this *Conn) Id() int64 {
	return this.id
}

//当读取数据时的回调函数
func (this *Conn) OnRead(call func([]byte)) {
	this.onReadCall = call
}

//是否设置了数据读取时的回调函数
func (this *Conn) IsSetOnRead() bool {
	if this.onReadCall != nil {
		return true
	}
	return false
}

//保存临时数据
func (this *Conn) Set(k, v interface{}) {
	this.cache.Store(k, v)
}

//读取临时数据
func (this *Conn) Get(k interface{}) (interface{}, bool) {
	return this.cache.Load(k)
}

//是否已经关闭
func (this *Conn) IsClose() bool {
	return this.closed
}

//关闭
func (this *Conn) Close() error {
	this.closed = true
	return this.Conn.Close()
}

//读取通讯数据(加个锁，让它有点次序)
//输入：rl=读取初始数据长度
func (this *Conn) Read() ([]byte, error) {
	this.LastAccessTime = time.Now()
	buf := make([]byte, 1024*10)
	if l, e := this.Conn.Read(buf); e != nil {
		return nil, e
	} else {
		return buf[:l], nil
	}
}

//发送数据
func (this *Conn) Write(b []byte) (int, error) {
	this.LastAccessTime = time.Now()
	return this.Conn.Write(b)
}

//触发读取回调函数
func (this *Conn) triggerReadCall(data []byte) {
	if this.onReadCall != nil {
		this.onReadCall(data)
	}
}
