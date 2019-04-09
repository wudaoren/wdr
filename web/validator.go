package web

import (
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
)

//验证函数
type ValidFunc func(reflect.Value, reflect.Type) bool

//注册自定义验证器
func RegisterValidation(param string, fun ValidFunc) {
	if valid, ok := binding.Validator.Engine().(*validator.Validate); ok {
		valid.RegisterValidation(param, func(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {
			return fun(field, fieldType)
		})
	}
}

//注册几个常用的验证器
func init() {
	RegisterValidation("password", func(v reflect.Value, t reflect.Type) bool {
		return true
	})
}
