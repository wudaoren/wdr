package orm

import (
	"fmt"
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

type Model struct {
	Engine        //继承引擎基类
	session       *Session
	obj           ModelInterface //model对象
	usePrimarykey bool           //使用默认主键查询
	data          interface{}
}

//初始化数据模型，obj为指针
func (this *Model) InitModel(obj ModelInterface) *Model {
	this.obj = obj
	this.usePrimarykey = true
	this.session = this.Session()
	return this
}

//表别名
func (this *Model) Alias(alias string) string {
	return fmt.Sprintf("`%s`.`%s`", this.obj.TableName(), alias)
}

//必须有的字段，比如age=0会强制写入数据库
func (this *Model) Must(cols ...string) *Model {
	this.Session().MustCols(cols...)
	return this
}

//忽略添加修改的字段
func (this *Model) Omit(cols ...string) *Model {
	this.Session().Omit(cols...)
	return this
}

//指定添加修改的字段
func (this *Model) Cols(cols ...string) *Model {
	this.Session().Cols(cols...)
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
	b, _ := this.Session().Table(this.obj).Limit(1).NoAutoCondition().Exist(this.obj)
	return b
}

//查询数据信息（根据ID）
func (this *Model) Get() bool {
	sx := this.Session().Table(this.obj)
	obj := this.preprocessing(sx)
	b, _ := sx.NoAutoCondition().Get(obj)
	return b
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
	return this.IsOneChange(l, e)
}

//删除数据（只针对主键id删除）
func (this *Model) Delete() error {
	sx := this.Session().Limit(1)
	obj := this.preprocessing(sx)
	l, e := sx.NoAutoCondition().Delete(obj)
	return this.IsOneChange(l, e)
}

//sql预处理
func (this *Model) preprocessing(sx *Session) interface{} {
	if this.usePrimarykey {
		sx.Id(this.obj.PrimaryKey())
	}
	this.usePrimarykey = true
	if this.data != nil {
		data := this.data
		this.data = nil
		return data
	}
	return this.obj
}
