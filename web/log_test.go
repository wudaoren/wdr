package web

import (
	"fmt"
	"testing"
)

func Test_log(t *testing.T) {
	l := NewLog("logs")
	l.Warn("哈哈哈")
	l.Info("哈哈哈")
	l.Fatal("哈哈哈")
	l.Error("哈哈哈")
	l.Errorf("12312%s", "gogog")
	var s = "as1 aaaa"
	fmt.Println(s[2:])
}
