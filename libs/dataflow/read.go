package dataflow

import (
	"sync"
)

const (
	GOON = 1 //继续
	OVER = 0 //结束
)

//字节流
type Flow struct {
	packet *Packet      //
	buffer []byte       //数据缓冲区(接收时用)
	onSucc func([]byte) //解包成功后回调
	head   []byte
	sync.Mutex
}

//新建数据流
func NewFlow(head [4]byte) *Flow {
	o := new(Flow)
	o.buffer = make([]byte, 0)
	o.head = head[:]
	return o
}

//解包成功后的回调
func (this *Flow) Succ(fn func([]byte)) *Flow {
	this.onSucc = fn
	return this
}

//输入字节流
func (this *Flow) Read(bt []byte) {
	this.Lock()
	defer this.Unlock()
	this.buffer = append(this.buffer, bt...)
	for {
		if this.parseBytes() == OVER {
			//fmt.Println("解码：", this.buffer)
			//fmt.Println("读取", code)
			return
		}
	}
}

//处理字节流,返回0=无需处理，1=需要再次处理
func (this *Flow) parseBytes() int8 {
	bufLen := len(this.buffer) //字节流长度
	if bufLen < 8 {
		return OVER
	}
	head := this.head
	//无字节流时进行头校验
	if this.packet == nil {
		checkHead := this.buffer[:4]
		if string(checkHead) != string(head[:]) {
			this.buffer = this.buffer[1:]
			return GOON
		}
		this.packet = &Packet{}
		this.packet.Head = checkHead
		this.buffer = this.buffer[4:]
		var dataLen int32
		byte2Number(this.buffer[:4], &dataLen)
		this.packet.BodyLen = dataLen
		this.buffer = this.buffer[4:]
	}
	bufLen = len(this.buffer)
	if bufLen >= int(this.packet.BodyLen+2) {
		stop := this.packet.BodyLen
		this.packet.Body = this.buffer[:stop]
		this.buffer = this.buffer[stop:]
		this.packet.Check = this.buffer[0]
		this.packet.Tail = this.buffer[1]
		if this.packet.Check == m256(int(this.packet.BodyLen)) && this.packet.Tail == tail(this.packet.Check) {
			body := gzipDecode(this.packet.Body)
			if this.onSucc != nil {
				this.onSucc(body)
			}
			this.packet = nil
			if len(this.buffer) > 0 {
				return GOON
			}
		}
	}
	return OVER
}
