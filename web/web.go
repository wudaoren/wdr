package web

import (
	"github.com/gin-gonic/gin"
)

//wab中间价
type Config struct {
	JWT struct {
		Issuser string
		Scrate  string
	}
	Log struct {
		Path string
	}
	Session struct {
		MaxAge  int
		KeyName string
	}
}

func Web() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
