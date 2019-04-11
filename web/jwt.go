package web

import (
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
)

/*
jwt中间价例子

type Handler struct {
}
//发行者
func (this *Handler) Issuer() string {
	return "test"
}

//设置密匙
func (this *Handler) Secret() string {
	return "abcdefaasdfasdft"
}

//获取token的方式
func (this *Handler) Token(w *Wontext) string {
	var token = ""
	if w.Request.Header["Authorization"] != nil {
		token = w.Request.Header["Authorization"][0]
	}
	return token
}

//当错误时
func (this *Handler) Error(w *Wontext) {
	token, e := w.GenerateToken("hait", "hait", 323)
	fmt.Println("jwt 加密结果：", token, e)
	w.Succ(100, token)
}

//当超时时
func (this *Handler) Timeout(w *Wontext) {

}
*/
const (
	KEY_JWT_OBJECT = "key_jwt_object"
)

//
type Claims struct {
	UserId   int64  `json:"userid"`
	Username string `json:"username"`
	Password string `json:"password"`
	Outtime  int64  `json:"outTime"`
	Group    string `json:"group"`
	jwt.StandardClaims
}

//jwt配置接口
type JWTConfiger interface {
	CacheFile() string     //数据库文件
	Issuer() string        //发行人
	Secret() string        //密匙
	Token(*Context) string //口令获取方式
	Error(*Context)        //当出现错误时
	Timeout(*Context)      //当超时时
}

//
type JWT struct {
	cache  *cache
	Token  string
	Issuer string
	Secret []byte
	Claims *Claims
	conf   JWTConfiger
}

//jwt处理器，初始化
func JWTHandler(conf JWTConfiger) gin.HandlerFunc {
	db, err := buntdb.Open(conf.CacheFile())
	if err != nil {
		log.Fatal(err)
	}
	return func(c *gin.Context) {
		var jwt = &JWT{
			Secret: []byte(conf.Secret()),
			Issuer: conf.Issuer(),
			conf:   conf,
			cache:  newCache(db),
		}
		c.Set(KEY_JWT_OBJECT, jwt)
	}
}

//jwt处理器，拦截器
func JWTintercept(group string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var w = NewContext(c)
		var jwt = w.JWT()
		jwt.Token = jwt.conf.Token(w)
		if jwt.Token == "" {
			jwt.conf.Error(w)
			c.Abort()
			return
		}
		claims, err := jwt.ParseToken(jwt.Token)
		if err != nil {
			jwt.conf.Error(w)
			c.Abort()
			return
		}
		if claims.Group != group {
			jwt.conf.Error(w)
			c.Abort()
			return
		}
		jwt.Claims = claims
		jwt.cache.setGroup(fmt.Sprint(group, claims.UserId))
		key := fmt.Sprint(claims.UserId)
		if jwt.cache.Get(key) == "" {
			jwt.conf.Error(w)
			c.Abort()
			return
		}
		var now = time.Now().Unix()
		if now > claims.Outtime && now < claims.ExpiresAt {
			jwt.conf.Timeout(w)
			c.Abort()
			return
		}
		c.Next()
	}
}

//生成token
//username 用户名
//password 密码
//id 用户id
//group 用户组
func (this *JWT) GenerateToken(username, password string, id int64, group string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(180 * time.Minute) //失效时间
	outTime := nowTime.Add(150 * time.Minute)    //刷新token时间
	claims := Claims{
		UserId:   id,
		Username: username,
		Password: password,
		Outtime:  outTime.Unix(),
		Group:    group,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    this.Issuer,
		},
	}
	//crypto.Hash加密方案
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//该方法内部生成签名字符串，再用于获取完整、已签名的token
	token, err := tokenClaims.SignedString(this.Secret)
	if err != nil {
		return "", err
	}
	this.cache.setGroup(fmt.Sprint(group, claims.UserId))
	key := fmt.Sprint(claims.UserId)
	if err := this.cache.Set(key, claims.Username); err != nil {
		return "", err
	}
	return token, nil
}

//解析
func (this *JWT) ParseToken(token string) (*Claims, error) {
	//用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return this.Secret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

//清除session
func (this *JWT) ClearSession() error {
	return this.cache.Clear()
}

//jwt使用的session
func (this *JWT) Session() *cache {
	return this.cache
}

type cache struct {
	db    *buntdb.DB
	group string
}

func newCache(db *buntdb.DB) *cache {
	object := new(cache)
	object.db = db
	object.group = "public"
	return object
}

//
func (this *cache) setGroup(name string) {
	this.group = name
}

//
func (this *cache) key(key string) string {
	return this.group + ":" + key
}

//
func (this *cache) Set(key string, value string) error {
	return this.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(this.key(key), value, nil)
		return err
	})
}

//
func (this *cache) Get(key string) (value string) {
	this.db.View(func(tx *buntdb.Tx) error {
		value, _ = tx.Get(this.key(key))
		return nil
	})
	return
}

//
func (this *cache) Del(key string) error {
	return this.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(this.key(key))
		return err
	})
}

func (this *cache) Clear() error {
	return this.db.Update(func(tx *buntdb.Tx) error {
		var delkeys []string
		var err error
		tx.AscendKeys(this.group+":*", func(k, v string) bool {
			delkeys = append(delkeys, k)
			return true // continue
		})
		for _, k := range delkeys {
			if _, err = tx.Delete(k); err != nil {
				return err
			}
		}
		return nil
	})
}
