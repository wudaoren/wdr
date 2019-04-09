package orm

//模型接口
type ModelInterface interface {
	PrimaryKey() interface{} //主键id
	TableName() string       //表名
}

//引擎接口
type Enginer interface {
	GetEngine() *Engine
}

//模块接口
type ModuleInterface interface {
	Install() error
}
