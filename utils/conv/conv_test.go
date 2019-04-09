package conv

import (
	"fmt"
	"testing"
)

type Test struct {
	Age  int `m:"pk"`
	Name string
}

func TestAny(t *testing.T) {
	test := &Test{Name: "爱爱爱"}
	var pk *int
	EachStruct(test, func(v *Val) {
		if v.Tag("m") == "pk" {
			pk = v.V.Addr().Interface().(*int)
		}
		fmt.Println("s", v.structName)
	})
	test.Age = 999
	fmt.Println("改变", test)
	fmt.Println("原值", *pk)
}
