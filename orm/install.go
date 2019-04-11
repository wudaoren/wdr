package orm

import (
	"github.com/go-xorm/xorm"
)

var models_list = make([]interface{}, 0)

//
func RegModels(models ...interface{}) {
	models_list = append(models_list, models...)
}

type Moduler interface {
	InitModels()
}

//安装
func Install(engine *xorm.Engine, modules ...Moduler) error {
	for _, module := range modules {
		module.InitModels()
	}
	return engine.Sync2(models_list...)
}
