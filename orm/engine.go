package orm

import (
	"errors"
	"reflect"

	"github.com/go-xorm/xorm"
)

type Engine struct {
	useTransaction bool
	session        *xorm.Session
	engine         *xorm.Engine
}

//设置引擎
func (this *Engine) SetXORM(eg *xorm.Engine) *Engine {
	this.engine = eg
	return this
}

//获取引擎
func (this *Engine) GetXORM() *xorm.Engine {
	if this.engine == nil {
		if engine == nil {
			panic("未初始化orm数据库引擎。")
		}
		this.engine = engine
	}
	return this.engine
}

//获取模型数据库会话，如果会话不存在则创建新会话
func (this *Engine) Session() *xorm.Session {
	if this.session == nil {
		this.session = this.GetXORM().NewSession()
	}
	return this.session
}

//获取引擎对象
func (this *Engine) GetEngine() *Engine {
	return this
}

//继承引擎,保证所有的事务都在同一条线上
func (this *Engine) ExtendEngine(eg Enginer) *Engine {
	this = eg.GetEngine()
	return this
}

//开启事务,如果继承了引擎，那么当前事务将使用上级的事务
func (this *Engine) Transaction(fn func(*Session) error) (err error) {
	sx := this.Session()
	if this.useTransaction {
		return fn(sx)
	}
	this.useTransaction = true
	defer func() {
		this.useTransaction = false
		if err != nil {
			sx.Rollback()
		} else {
			sx.Commit()
		}
	}()
	sx.Begin()
	err = fn(sx)
	return
}

//不使用继承方式直接创建一个模型
func (this *Engine) Model(obj ModelInterface) *Model {
	var model = new(Model)
	model.engine = this.engine
	return model.InitModel(obj)
}

//当更新的数据不是1条时返回错误
func (this *Engine) IsOneChange(changes int64, e error) error {
	if e != nil {
		return e
	} else if changes != 1 {
		return errors.New("存储失败")
	}
	return nil
}

//查询异常判断
func (this *Engine) IsGetOK(ok bool, e error) error {
	if e != nil {
		return e
	} else if !ok {
		return errors.New("查询错误")
	}
	return nil
}

//分页查询
//listPtr	= 查询列表指针
//page	= 页码
//limit	= 每页查询数量
func (this *Engine) FindPage(session *xorm.Session, listPtr interface{}, page, limit int) (total int64, err error) {
	if page <= 0 {
		page = 1
	}
	if limit == 0 {
		limit = 20
	}
	total, err = session.Clone().Select("count(*)").Count()
	if err != nil {
		return 0, err
	}
	start := (page - 1) * limit
	err = session.Limit(limit, start).Find(listPtr)
	return total, err
}

//通过反射获取切片结构体id切片
func (this *Engine) GetFieldList(list interface{}, idFieldName string) ([]interface{}, error) {
	v := reflect.ValueOf(list).Elem()
	if v.Kind() != reflect.Slice {
		return nil, errors.New("必须传入一个slice类型数据")
	}
	fieldMaps := make(map[interface{}]bool) //id集合
	for i := 0; i < v.Len(); i++ {
		ptr := v.Index(i)
		if reflect.ValueOf(ptr).Kind() != reflect.Struct {
			return nil, errors.New("slice的元素必须是struct")
		}
		if ptr.Kind() != reflect.Ptr {
			ptr = ptr.Addr()
		}
		ptr = ptr.Elem()
		val := ptr.FieldByName(idFieldName).Interface()
		fieldMaps[val] = true
	}
	filedSlice := make([]interface{}, 0)
	for field, _ := range fieldMaps {
		filedSlice = append(filedSlice, field)
	}
	return filedSlice, nil
}
