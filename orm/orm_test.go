package orm

import (
	"database/sql"
	"fmt"
	"testing"

	  "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

type User struct {
	Id      int    `xorm:"int(11) autoincr pk" model:""`
	Name    string `xorm:"varchar(32)" put:"user"`
	Age     int64  `xorm:"int(11)" put:"user"`
	Version int    `xorm:"version" put:"sss"`
	Model   `xorm:"-"`
}

func NewUser() *User {
	u := new(User)
	u.InitModel(u)
	return u
}

func (this *User) TableName() string {
	return "User"
}
func (this *User) PrimaryKey() interface{} {
	return this.Id
}

//乐观锁测试
func initdata() {
	mysql.Config
	eg, _ := xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", "root", "", "127.0.0.1:3306", "test"))
	eg.ShowSQL(true)
	eg.SetMapper(new(core.SameMapper))
	InitEngine(eg)
}

func TestTrans(t *testing.T) {
	initdata()
	a := NewUser()
	_, e := a.Session().Where("id=?", "' show databases;").Get(a)
	fmt.Println(e)
	new(sql.DB).Exec()
}
