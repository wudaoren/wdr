package log

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Log struct {
	debug  bool
	logger *log.Logger
}

//
func NewLog(dir string, debug bool) *Log {
	object := new(Log)
	object.debug = debug
	if !debug {
		os.MkdirAll(dir, 0776)
		path := fmt.Sprintf("%s/%s.log", dir, time.Now().Format("2006.01.02"))
		file, _ := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0776)
		object.logger = log.New(file, "", log.LstdFlags|log.Llongfile)
	}
	return object
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
func (this *Log) Print(arg ...interface{}) {
	fmt.Println(arg...)
}

//
func (this *Log) Printf(fomrat string, arg ...interface{}) {
	fmt.Printf(fomrat+"\n", arg...)
}

//
func (this *Log) println(prefix string, arg ...interface{}) {
	prefix += "     "
	prefix = "[" + prefix[:5] + "]"
	if this.debug {
		fmt.Println(prefix, fmt.Sprint(arg...))
	} else {
		this.logger.SetPrefix(prefix)
		this.logger.Println(arg...)
	}
}
