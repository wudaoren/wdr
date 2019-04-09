package dataflow

//打包(并进行gzip压缩）
func Pack(head [4]byte, body []byte) []byte {
	body = gzipEncode(body)
	dataLen := len(body)
	s := new(Packet)
	//组成数据流
	s.Head = head[:]           //
	s.BodyLen = int32(dataLen) //
	s.Body = body              //
	s.Check = m256(dataLen)    //
	s.Tail = tail(s.Check)     //
	//组成字节码
	packBytes := make([]byte, 0)
	packBytes = append(packBytes, s.Head...)
	packBytes = append(packBytes, number2byte(s.BodyLen)...)
	packBytes = append(packBytes, s.Body...)
	packBytes = append(packBytes, s.Check)
	packBytes = append(packBytes, s.Tail)
	return packBytes
}
