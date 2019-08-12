package web

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/url"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	SESSION = "SESSION"
)

func SessionHandler(name string, maxLifeTime int64) gin.HandlerFunc {
	session := &SessionEngine{
		seesionKey:  name,
		maxLifeTime: maxLifeTime,
		data:        new(sync.Map),
	}
	go func() {
		for {
			time.Sleep(time.Second * 30)
			session.gc()
		}
	}()
	return func(c *gin.Context) {
		sessionId, _ := c.Cookie(name)
		if sessionId == "" {
			key := fmt.Sprintf("%d-%x", time.Now().Unix(), rand.Int63())
			md5Ctx := md5.New()
			md5Ctx.Write([]byte(key))
			cipherStr := md5Ctx.Sum(nil)
			sessionId = hex.EncodeToString(cipherStr)
			c.SetCookie(name, url.QueryEscape(sessionId), 3600*24*365, "/", "", false, false)
		}
		c.Set("SESSION", session.getSession(sessionId))
		c.Next()
	}
}

//默认session engine
type SessionEngine struct {
	seesionKey  string
	maxLifeTime int64 //最长保存时间（秒）
	visited     int   //访问次数
	data        *sync.Map
}

type sessionData struct {
	sess *MemSession
	time int64
}

func (this *SessionEngine) getSession(sessionId string) *MemSession {
	data := new(sessionData)
	if res, ok := this.data.Load(sessionId); ok {
		data = res.(*sessionData)
	}
	overTime := this.isTimeout(data)
	if overTime > this.maxLifeTime {
		data.sess = new(MemSession)
		this.data.Store(sessionId, data)
	} else if overTime > this.maxLifeTime/2 && overTime < this.maxLifeTime {
		this.data.Store(sessionId, data)
	}
	data.sess.id = sessionId
	return data.sess
}

//是否超时
func (this *SessionEngine) isTimeout(data *sessionData) int64 {
	nowTimestemp := time.Now().Unix()
	overTime := nowTimestemp - data.time //小于0表示正常，大于0表示超时
	data.time = nowTimestemp
	return overTime
}

//
func (this *SessionEngine) gc() {
	this.data.Range(func(k, v interface{}) bool {
		if data, ok := v.(*sessionData); ok && this.isTimeout(data) >= 0 {
			this.data.Delete(k)
		}
		return true
	})
}

//session对象
type MemSession struct {
	id string
	sync.Map
}

func (this *MemSession) Id() string {
	return this.id
}

//
func (this *MemSession) Set(key string, value interface{}) {
	this.Store(key, value)
}

func (this *MemSession) Get(key string) interface{} {
	res, _ := this.Load(key)
	return res
}

func (this *MemSession) Del(key string) {
	this.Delete(key)
}

func (this *MemSession) Clear() {
	this.Map = sync.Map{}
}
