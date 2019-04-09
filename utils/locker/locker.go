package locker

import (
	"fmt"
	"sync"
)

//锁缓存
var lockMems = new(sync.Map)

//申请一个锁
func Use(keys ...interface{}) *sync.Mutex {
	key := fmt.Sprint(keys...)
	lock, ok := lockMems.Load(key)
	if !ok {
		lock = new(sync.Mutex)
		lockMems.Store(key, lock)
	}
	return lock.(*sync.Mutex)
}
