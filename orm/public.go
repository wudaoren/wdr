package orm

import (
	"github.com/go-xorm/xorm"
)

//session别名
type Session = xorm.Session

//数据库引擎
var engine *xorm.Engine

//初始化引擎
func InitEngine(eg *xorm.Engine) {
	engine = eg
}
