package web

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"
)

//gin context加强版

var (
	FORMAT_ERROR   = errors.New("数据格式错误")
	VALIDATE_ERROR = errors.New("数据验证错误")
)

//升级版context
type Context struct {
	isResponse bool
	*gin.Context
}

func NewContext(c *gin.Context) *Context {
	obj := new(Context)
	obj.Context = c
	return obj
}

//获取session对象
func (this *Context) Session() {

}

//json输出
func (this *Context) JSON(code int, res, msg interface{}) {
	if this.isResponse {
		return
	}
	type ret struct {
		Code   int         `json:"code"`
		Result interface{} `json:"result"`
		Msg    interface{} `json:"msg"`
	}
	this.Context.JSON(200, &ret{code, res, msg})
	this.isResponse = true
}

//只返回状态码
func (this *Context) Code(code int) {
	this.JSON(code, nil, nil)
}

//成功返回（只输出data）
func (this *Context) Res(code int, data interface{}) {
	this.JSON(code, data, nil)
}

//错误返回（只输出msg)
func (this *Context) Msg(code int, msg interface{}) {
	this.JSON(code, nil, msg)
}

//错误,并自动记录日志
func (this *Context) Err(code int, err error) {
	this.JSON(code, nil, err.Error())
}

//重写绑定方法
func (this *Context) ShouldBind(obj interface{}) error {
	e := this.Context.ShouldBind(obj)
	return this.BindError(e)
}

//重写绑定方法
func (this *Context) ShouldBindQuery(obj interface{}) error {
	e := this.Context.ShouldBindQuery(obj)
	return this.BindError(e)
}

//重写绑定方法
func (this *Context) ShouldBindJSON(obj interface{}) error {
	e := this.Context.ShouldBindJSON(obj)
	return this.BindError(e)
}

//数据绑定错误判断（如果某个字段验证失败，则返回该字段note标签作为提示符
//如果没有指定标签，则返回数据验证错误
func (this *Context) BindError(e error) error {
	if e == nil {
		return nil
	}
	switch errs := e.(type) {
	case validator.ValidationErrors:
		for _, field := range errs {
			t := field.RefField.Type().Elem()
			if f, ok := t.FieldByName(field.Field); ok {
				if note := f.Tag.Get("note"); note != "" {
					return errors.New(note)
				} else {
					return VALIDATE_ERROR
				}
			}
		}
	}
	return FORMAT_ERROR //数据格式错误
}

//获取jwt对象
func (this *Context) JWT() *JWT {
	jwt, ok := this.Get(KEY_JWT_OBJECT)
	if !ok {
		panic("jwt未启用.")
	}
	return jwt.(*JWT)
}
