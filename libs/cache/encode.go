package cache

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

//加密解密都用此方法
func easyEncode(data []byte) []byte {
	l := len(data)
	ret := make([]byte, l)
	for k, b := range data {
		ret[l-k-1] = 255 - b
	}
	return ret
}

//gzip压缩
func enGzip(bt []byte) []byte {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()
	w.Write(bt)
	w.Flush()
	return b.Bytes()
}

//gzip解压
func deGzip(bt []byte) []byte {
	var rd bytes.Buffer
	rd.Write(bt)
	r, e := gzip.NewReader(&rd)
	defer r.Close()
	if e != nil {
		return []byte{}
	}
	dt, _ := ioutil.ReadAll(r)
	return dt

}
