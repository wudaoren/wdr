package orm

import (
	"fmt"
	"reflect"
	"strings"

	//	"reflect"
	"sync"
)

/*------------------  使用示例  -----------------

type User struct {
	Id    int    `xorm:"int(11) autoincr pk" model:""`
	Name  string `xorm:"varchar(32)"`
	Age   int8   `xorm:"int(11)"`
	Model `xorm:"-"`
}

func NewUser() *User {
	u := new(User)
	u.InitModel(u)
	return u
}
注：一定要在结构体的主键字段的tag里声明model,默认主键名与字段名相同，也可任意设置

*/
var tableCache = new(sync.Map)

//清除所有缓存
func ClearCache() {
	tableCache = new(sync.Map)
}

type Model struct {
	Engine                       //继承引擎基类
	obj           ModelInterface //model对象
	usePrimarykey bool           //使用默认主键查询
	omit          []string
	cols          []string
	must          []string
	data          interface{}
	cache         bool
}

//初始化数据模型，obj为指针
func (this *Model) InitModel(obj ModelInterface) *Model {
	this.obj = obj
	this.omit = make([]string, 0)
	this.cols = make([]string, 0)
	this.must = make([]string, 0)
	this.usePrimarykey = true
	return this
}

//表别名
func (this *Model) Alias(alias string) string {
	return fmt.Sprintf("`%s` as `%s`", this.obj.TableName(), alias)
}

//必须有的字段，比如age=0会强制写入数据库
func (this *Model) Must(cols ...string) *Model {
	this.must = append(this.must, cols...)
	return this
}

//忽略添加修改的字段
func (this *Model) Omit(cols ...string) *Model {
	this.omit = append(this.omit, cols...)
	return this
}

//指定添加修改的字段
func (this *Model) Cols(cols ...string) *Model {
	this.cols = append(this.cols, cols...)
	return this
}

//使用条件查询
func (this *Model) Where(query string, args ...interface{}) *Model {
	this.Session().Where(query, args...)
	this.usePrimarykey = false
	return this
}

//输入输出数据
func (this *Model) Data(data interface{}) *Model {
	this.data = data
	return this
}

//根据指定字段获取对象
func (this *Model) Match(column string, data interface{}) *Model {
	this.Session().And("`"+column+"`=?", data)
	this.usePrimarykey = false
	return this
}

//是否存在
func (this *Model) Exists() bool {
	sx := this.Session()
	this.preprocessing(sx)
	b, _ := sx.Table(this.obj).Limit(1).NoAutoCondition().Exist(this.obj)
	return b
}

//查询数据信息（根据ID）
func (this *Model) Get() bool {
	sx := this.Session().Table(this.obj)
	obj := this.preprocessing(sx)
	b, _ := sx.NoAutoCondition().Get(obj)
	return b
}

//查询对象，如果有缓存，则使用缓存
func (this *Model) GetById(id interface{}) error {
	key := this.obj.TableName()
	table, ok := tableCache.Load(key)
	if !ok {
		table = new(sync.Map)
		tableCache.Store(key, table)
	}
	tableMap := table.(*sync.Map)
	data, ok := tableMap.Load(id)
	if ok {
		val := reflect.ValueOf(this.obj).Elem()
		dta := reflect.ValueOf(data).Elem()
		typ := reflect.TypeOf(this.obj)
		for i := 0; i < val.NumField(); i++ {
			tag := typ.Elem().Field(i).Tag
			if field := val.Field(i); field.CanSet() && strings.TrimSpace(tag.Get("xorm")) != "-" {
				field.Set(dta.Field(i))
			}
		}
		return nil
	}
	sx := this.Session().Table(this.obj).ID(id)
	if err := this.IsGetOK(sx.Get(this.obj)); err != nil {
		return err
	}
	tableMap.Store(id, this.obj)
	return nil
}

//清理缓存
func (this *Model) clearCacheById(id interface{}) {
	key := this.obj.TableName()
	table, ok := tableCache.Load(key)
	if !ok {
		return
	}
	table.(*sync.Map).Delete(id)
}

//添加数据
func (this *Model) Insert() error {
	sx := this.Session().Table(this.obj)
	obj := this.preprocessing(sx)
	l, e := sx.InsertOne(obj)
	return this.IsOneChange(l, e)
}

//修改数据
func (this *Model) Update() error {
	sx := this.Session().Limit(1)
	l, e := sx.Update(this.preprocessing(sx))
	this.clearCacheById(this.obj.PrimaryKey())
	return this.IsOneChange(l, e)
}

//删除数据（只针对主键id删除）
func (this *Model) Delete() error {
	sx := this.Session().Limit(1)
	obj := this.preprocessing(sx)
	l, e := sx.NoAutoCondition().Delete(obj)
	this.clearCacheById(this.obj.PrimaryKey())
	return this.IsOneChange(l, e)
}

//sql预处理
func (this *Model) preprocessing(sx *Session) interface{} {
	if this.usePrimarykey {
		sx.Id(this.obj.PrimaryKey())
	}
	if len(this.cols) > 0 {
		sx.Cols(this.cols...)
	}
	if len(this.omit) > 0 {
		sx.Omit(this.omit...)
	}
	if len(this.must) > 0 {
		sx.MustCols(this.must...)
	}
	this.usePrimarykey = true
	this.omit = make([]string, 0)
	this.cols = make([]string, 0)
	this.must = make([]string, 0)
	if this.data != nil {
		data := this.data
		this.data = nil
		return data
	}
	return this.obj
}
