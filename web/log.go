package web

import (
	"fmt"
	"log"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	KEY_LOG_OBJECT = "key_log_object"
)

type Log struct {
	path   string
	date   string
	logger *log.Logger
}

//logo处理中间件
func LogHandler(path string, errorHandler func(*gin.Context)) gin.HandlerFunc {
	l := NewLog(path)
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				l.Fatal(string(httprequest), err, string(stack))
				errorHandler(c)
			}
		}()
		c.Set(KEY_LOG_OBJECT, l)
		c.Next()
	}
}

//
func NewLog(path string) *Log {
	obj := new(Log)
	obj.path = path
	os.MkdirAll(path, 0776)
	return obj
}

//
func (this *Log) Warn(arg ...interface{}) {
	this.println("warn", arg...)
}

//
func (this *Log) Warnf(fomrat string, arg ...interface{}) {
	this.Warn(fmt.Sprintf(fomrat, arg...))
}

//
func (this *Log) Error(arg ...interface{}) {
	this.println("error", arg...)
}

//
func (this *Log) Errorf(fomrat string, arg ...interface{}) {
	this.Error(fmt.Sprintf(fomrat, arg...))
}

//
func (this *Log) Info(arg ...interface{}) {
	this.println("info", arg...)
}

//
func (this *Log) Infof(fomrat string, arg ...interface{}) {
	this.Info(fmt.Sprintf(fomrat, arg...))
}

//
func (this *Log) Fatal(arg ...interface{}) {
	this.println("fatal", arg...)
}

//
func (this *Log) Fatalf(fomrat string, arg ...interface{}) {
	this.Warn(fmt.Sprintf(fomrat, arg...))
}

//
func (this *Log) Debug(arg ...interface{}) {
	fmt.Println(arg...)
}

//
func (this *Log) Debugf(fomrat string, arg ...interface{}) {
	fmt.Printf(fomrat, arg...)
}

//
func (this *Log) println(prefix string, arg ...interface{}) {
	prefix += "     "
	prefix = "[" + prefix[:5] + "]"
	now := time.Now()
	if this.path == "" {
		fmt.Print(prefix, now.Format(" 2006-01-02 15:04:05 "))
		fmt.Println(arg...)
		return
	}
	date := now.Format("2006-01-02")
	if this.date != date {
		fmt.Println(fmt.Sprintf("%s/%s.log", this.path, date))
		w, _ := os.OpenFile(fmt.Sprintf("%s/%s.log", this.path, date), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0776)
		this.logger = log.New(w, "", log.LstdFlags|log.Llongfile)
		this.date = date
	}
	this.logger.SetPrefix(prefix)
	this.logger.Println(arg...)
}
