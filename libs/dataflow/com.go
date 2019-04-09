package dataflow

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io/ioutil"
)

//包
type Packet struct {
	Head    []byte //头部：4个字节 = 必须自定义
	BodyLen int32  //长度：4个字节 = len(Body)
	Body    []byte //数据：n个字节 = json序列化，并gzip
	Check   byte   //校验：1个字节 = m256
	Tail    byte   //尾部：1个字节 = 255-Check
}

//gzip压缩
func gzipEncode(bt []byte) []byte {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()
	w.Write(bt)
	w.Flush()
	return b.Bytes()
}

//gzip解压
func gzipDecode(bt []byte) []byte {
	var rd bytes.Buffer
	rd.Write(bt)
	r, e := gzip.NewReader(&rd)
	defer r.Close()
	if e != nil {
		return []byte{}
	}
	data, _ := ioutil.ReadAll(r)
	return data

}

func number2byte(data interface{}) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, data)
	return buf.Bytes()
}

func byte2Number(b []byte, data interface{}) {
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.BigEndian, data)
	return
}

//尾
func tail(l byte) byte {
	return 255 - l
}

//模256
func m256(l int) byte {
	return byte(l % 256)
}
