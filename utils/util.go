package utils

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
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

//读取配置文件(自动替换配置文件里面的/**/注释)
//输入：file=配置文件名，structPtr=配置对象
func ReadConfig(file string, structPtr interface{}) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	compile, _ := regexp.Compile(`\/\*[^\*]+\*\/`)
	data = compile.ReplaceAll(data, []byte(""))
	return json.Unmarshal(data, structPtr)
}
