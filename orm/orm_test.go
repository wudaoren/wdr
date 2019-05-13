package orm

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
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
func init() {
	eg, _ := xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", "root", "", "127.0.0.1:3306", "test"))
	eg.ShowSQL(true)
	eg.SetMapper(new(core.SameMapper))
}

func TestTrans(t *testing.T) {
	a := NewUser()
	b := NewUser()
	c := New
}
