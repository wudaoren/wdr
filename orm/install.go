package orm

import (
	"github.com/go-xorm/xorm"
)

var models_list = make([]interface{}, 0)

func RegModels(models ...interface{}) {
	models_list = append(models_list, models...)
}

//安装
func Install(engine *xorm.Engine) error {
	return engine.Sync2(models_list...)
}
