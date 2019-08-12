package utils

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"time"
)

//执行函数，捕获错误
func Try(fn func(), errFn ...func(interface{})) {
	defer func() {
		if err := recover(); err != nil && len(errFn) == 1 {
			errFn[0](err)
		}
	}()
	fn()
}

//执行协程，并捕获错误
func Go(fn func(), errFn ...func(interface{})) {
	go Try(fn, errFn...)
}

//读取配置文件,(自动替换配置文件里面的/**/注释)
//输入：file=配置文件名，structPtrs=配置对象，可多个
func ReadConfig(file string, structPtrs ...interface{}) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	compile, _ := regexp.Compile(`\/\*[^\*]+\*\/`)
	data = compile.ReplaceAll(data, []byte(""))
	for _, ptr := range structPtrs {
		if err := json.Unmarshal(data, ptr); err != nil {
			return err
		}
	}
	return nil
}

//当前日期
func Date() string {
	return time.Now().Format(FOMART_DATE)
}

//当前时间
func DateTime() string {
	return time.Now().Format(FORMAT_DATETIME)
}
