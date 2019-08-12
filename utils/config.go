package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//读取简单的ini配置

type EasyIni struct {
	ini string
	m   map[string]string
}

func NewEasyIni(str string) (*EasyIni, error) {
	object := new(EasyIni)
	object.m = make(map[string]string)
	str += "\n"
	//去除注释
	r, _ := regexp.Compile("//.*\n")
	b := r.ReplaceAll([]byte(str), []byte("\n"))
	str = string(b)
	//数据分离
	arr1 := strings.Split(strings.TrimSpace(str), "\n")
	for i1, v1 := range arr1 {
		arr2 := strings.Split(v1, "=")
		if len(arr2) != 2 {
			return nil, errors.New(fmt.Sprintf("第%d行错误", i1+1))
		}
		object.m[strings.TrimSpace(arr2[0])] = strings.TrimSpace(arr2[1])
	}
	return object, nil
}

//
func (this *EasyIni) Get(key string) string {
	return this.m[key]
}

//
func (this *EasyIni) GetInt(key string) int {
	if v, ok := this.m[key]; ok {
		r, _ := strconv.ParseInt(v, 10, 32)
		return int(r)
	}
	return 0
}

//
func (this *EasyIni) GetInt64(key string) int64 {
	if v, ok := this.m[key]; ok {
		r, _ := strconv.ParseInt(v, 10, 64)
		return int64(r)
	}
	return 0
}

//
func (this *EasyIni) GetFloat32(key string) float32 {
	if v, ok := this.m[key]; ok {
		r, _ := strconv.ParseFloat(v, 32)
		return float32(r)
	}
	return 0
}

//
func (this *EasyIni) GetFloat64(key string) float64 {
	if v, ok := this.m[key]; ok {
		r, _ := strconv.ParseFloat(v, 64)
		return float64(r)
	}
	return 0
}
