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
func (this *User) PrimaryKey() interface{} {
	return this.Id
}

//乐观锁测试
func TestIncr(t *testing.T) {
	eg, _ := xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", "root", "", "127.0.0.1:3306", "test"))
	eg.ShowSQL(true)
	eg.SetMapper(new(core.SameMapper))
	var user = &User{
		Age: 22,
	}
	eg.Cols("Age").Incr("Age", 1).ID(3).Update(user)
}

func aTestModel(t *testing.T) {
	eg, _ := xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", "root", "", "127.0.0.1:3306", "test"))
	InitEngine(eg)

	cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	eg.SetDefaultCacher(cacher)

	eg.ShowSQL(true)
	//	eg.Sync2(new(User))

	var data = make([]User, 0)
	var ndata = make([]User, 0)
	data = append(data, User{Id: 2, Age: 33})
	data = append(data, User{Id: 2})
	//创建数据
	u := NewUser()
	err := u.PutOnData(u, "Id", "Id", "user", &data, &ndata)
	fmt.Println("--", data, err)
	for _, v := range data {
		fmt.Println("输出；", v)
	}

	//list := make([]User, 0)
	//sx := u.Session().Table(u)
	//u.FindPage(sx, &list, 1, 10)
	fmt.Println(data)
	return
	//创建数据
	u.Name = "test"
	u.Age = 58
	e := u.Insert()
	fmt.Println("创建结果：", e, u.Id)
	//******
	u2 := NewUser()
	u2.Id = u.Id
	u2.Get()
	fmt.Println("查询结果2:", u2)

	//*********
	u.Name = "haitao"
	e1 := u.Update()
	fmt.Println("修改结果：", e1)
	//*********
	u4 := NewUser()
	u4.Id = u.Id
	u4.Get()
	fmt.Println("查询结果4:", u4)
	//*********
	u5 := NewUser()
	u5.Id = u.Id
	u5.Get()
	fmt.Println("查询结果5:", u5)
	//*********
	u6 := NewUser()
	u6.Id = u.Id
	u6.Get()
	fmt.Println("查询结果6:", u6)
	//***********
	e2 := u.Delete()
	fmt.Println("删除结果:", e2)
}
