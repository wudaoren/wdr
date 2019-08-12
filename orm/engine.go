package orm

import (
	"errors"
	"reflect"
	"runtime/debug"
	"wdr/errs"

	"github.com/go-xorm/xorm"
)

type storage struct {
	openTrans bool
	session   *xorm.Session
	engine    *xorm.Engine
}

type Engine struct {
	storage *storage
	sons    []*Engine
	extend  bool
}

func (this *Engine) init() {
	if this.storage == nil {
		this.storage = new(storage)
		if engine == nil {
			panic("未初始化orm数据库引擎。")
		}
		this.storage.engine = engine
		this.sons = make([]*Engine, 0)
	}
}

//设置引擎
func (this *Engine) SetEngine(eg *xorm.Engine) {
	this.init()
	this.storage.engine = eg
	return
}

//获取引擎
func (this *Engine) GetEngine() *xorm.Engine {
	this.init()
	return this.storage.engine
}

//获取模型数据库会话，如果会话不存在则创建新会话
func (this *Engine) Session() *xorm.Session {
	this.init()
	if this.storage.session == nil {
		this.storage.session = this.GetEngine().NewSession()
	}
	return this.storage.session
}

type Enginer interface {
	grantSon(*Engine)
}

//被继承
//给继承者赋予同样session
func (this *Engine) grantSon(son *Engine) {
	this.Session()
	son.storage = this.storage
	if !son.extend {
		this.sons = append(this.sons, son)
		son.extend = true
	}
}

//继承引擎,如果有子继承者，则将子继承者全部统一
func (this *Engine) ExtendEngine(parent Enginer) {
	parent.grantSon(this)
	if this.sons != nil {
		for _, eg := range this.sons {
			eg.ExtendEngine(this)
		}
	}
}

//开启事务,如果继承了引擎，那么当前事务将使用上级的事务
func (this *Engine) Transaction(fn func(*Session) error) (err error) {
	sx := this.Session()
	if this.storage.openTrans {
		return fn(sx)
	}
	this.storage.openTrans = true
	defer func() {
		this.storage.openTrans = false
		if e := recover(); e != nil {
			err = errs.Fatal(e, string(debug.Stack()))
			sx.Rollback()
			return
		}
		if err != nil {
			sx.Rollback()
			return
		}
		sx.Commit()
	}()
	sx.Begin()
	err = fn(sx)
	return
}

//不使用继承方式直接创建一个模型
func (this *Engine) Model(obj ModelInterface) *Model {
	var model = new(Model)
	model.ExtendEngine(this)
	return model.InitModel(obj)
}

//当更新的数据不是1条时返回错误
func (this *Engine) IsOneChange(changes int64, e error) error {
	if e != nil {
		return errs.Fatal(e.Error())
	} else if changes != 1 {
		return errors.New("存储失败")
	}
	return nil
}

//查询异常判断
func (this *Engine) IsGetOK(ok bool, e error) error {
	if e != nil {
		return errs.Fatal(e.Error())
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
	if err != nil {
		return 0, errs.Fatal(err.Error())
	}
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
