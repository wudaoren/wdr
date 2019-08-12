package utils

import (
	"testing"
)

func TestConfig(t *testing.T) {
	ini := `
	王海涛=32s
	流行派=123.3 //阿斯顿发
	
	`
	c, e := NewEasyIni(ini)
	if e != nil {
		t.Log(e.Error())
		return
	}
	t.Log(c.Get("王海涛"))
	t.Log(c.GetInt("王海涛"))
	t.Log("流行病", c.GetFloat32("流行派"))
}
