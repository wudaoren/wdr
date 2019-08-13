package utils

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

//sha1加密
func Sha1(data string) string {
	sha := sha1.New()
	sha.Write([]byte(data))
	return hex.EncodeToString(sha.Sum([]byte("")))
}

//sha1加密
func Sha256(data string) string {
	sha := sha256.New()
	sha.Write([]byte(data))
	return hex.EncodeToString(sha.Sum([]byte("")))
}

//md5加密
func Md5(str string) string {
	object := md5.New()
	object.Write([]byte(str))
	cipherStr := object.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

//crc校验码
func CRC16(data []byte) []byte {
	var crc16 uint16 = 0xffff
	l := len(data)
	for i := 0; i < l; i++ {
		crc16 ^= uint16(data[i])
		for j := 0; j < 8; j++ {
			if crc16&0x0001 > 0 {
				crc16 = (crc16 >> 1) ^ 0xA001
			} else {
				crc16 >>= 1
			}
		}
	}
	packet := make([]byte, 2)
	packet[0] = byte(crc16 & 0xff)
	packet[1] = byte(crc16 >> 8)
	return packet
}

//生成sn，sn长度为18个字符
func MakeSN(head string) string {
	str := fmt.Sprint(head, time.Now().Format("20060102150405"), rand.Int())
	return str[:18]
}

//
func CheckSN(sn string) bool {
	return len(sn) == 18
}

//将byte类型转换成字符串，例如[]byte{1,2,3,4,5}转换成"1,2,3,4,5"
//输入：bt=字节码,dilimiter=连接字符串
func Byte2Str(bt []byte, dilimiter string) string {
	str := ""
	l := len(bt)
	for k, v := range bt {
		if k == l-1 {
			dilimiter = ""
		}
		str = str + strconv.Itoa(int(v)) + dilimiter
	}
	return str
}

//将字符串转换为byte数据
func Str2Byte(str []string) (bt []byte, err error) {
	for _, b := range str {
		nb, e := strconv.Atoi(strings.TrimSpace(b))
		if e != nil {
			return nil, e
		}
		bt = append(bt, byte(nb))
	}
	return bt, nil
}

//16进制字符串转换为byte数据
func Hex2Byte(str []string) (bt []byte, err error) {
	for _, v := range str {
		b, e := strconv.ParseInt(v, 16, 10)
		if e != nil {
			return nil, err
		}
		bt = append(bt, byte(b))
	}
	return bt, nil
}

//将byte转换成数字类型
//输入：b=byte,data=要转换输出的变量'请输入指针
func Byte2Number(b []byte, data interface{}) {
	b_buf := bytes.NewBuffer(b)
	binary.Read(b_buf, binary.BigEndian, data)
	return
}

//数字转换成byte类型
func Number2Byte(data interface{}) []byte {
	b_buf := bytes.NewBuffer([]byte{})
	binary.Write(b_buf, binary.BigEndian, data)
	return b_buf.Bytes()
}

//快速新建一个错误对象，可以输入多个参数
func Error(fomart string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(fomart, args...))
}

//当前程序目录
func GetCurrentPath() string {
	s, err := exec.LookPath(os.Args[0])
	if err != nil {
		return ""
	}
	i := strings.LastIndex(s, "\\")
	path := string(s[0 : i+1])
	return path
}

//给字符串加校验值
//比如输入：ABC90 则返回ABC90P，P是算出来的校验值
func HexStrM256(str string) string {
	var c = []byte("123456789ABCDEF0")
	var l = len(c)
	var d = 0
	for _, i := range str {
		d += int(i)
	}
	return str + string(c[d%l])
}

// 产生随机字符串
//输入：size=字符串长度，kind=字符串类型 0数字，1小写字母，2大写字母，3数字大小写字母
func Rand(size int, kind int) string {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all {
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return string(result)
}

//根据ip获取地址
func IP2Address(ip string) string {
	url := fmt.Sprint("https://api.ttt.sh/ip/qqwry/", ip)
	req, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer req.Body.Close()
	out, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return ""
	}
	var response struct {
		Code    int    `json:"code"`
		Address string `json:"address"`
	}
	if err := json.Unmarshal(out, &response); err != nil {
		return ""
	}
	arr := strings.Split(response.Address, " ")
	if len(arr) > 0 {
		return arr[0]
	}
	return ""
}

func Print(str string) {
	println("ni shuru deshi:", str)
}
