package task

import (
	"testing"
	"time"
)

func TestT(t *testing.T) {
	NewTask("哈哈", "*/2 * * * * *", func() error {
		println(time.Now().Format("2006-01-02 15:04:05"))
		return nil
	})
	StartTask()
	time.Sleep(time.Hour)
}
