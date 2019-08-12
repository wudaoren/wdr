package cache

import (
	"encoding/json"
	"io/ioutil"
	"sync"
	"time"
)

//本缓存包，可以方便的使用缓存来保存热数据

//缓存数据结构
type Data struct {
	Key     string
	Index   int         //索引编号
	Value   interface{} //缓存数据的值
	Timeout int64       //到期时间(单位:秒）
}

//回调
type CallFunc func(k, v interface{}) error

//内存缓存
type Cache struct {
	data       map[string]*Data //数据结构
	sync.Mutex                  //互斥锁
	size       int              //初始内存大小
	maked      bool             //map是否已经初始化
	expire     int64            //延迟时间
	onDelCall  CallFunc         //当删除前
	onSetCall  CallFunc         //当设置前
	index      int              //当前索引
}

//创建一个内存缓存
//size=初始内存大小
func NewCache(size int) *Cache {
	obj := new(Cache)
	obj.size = size
	obj.index = 0
	obj.data = make(map[string]*Data, obj.size)
	return obj
}

//导入map数据到缓存数据结构中
//导入后不会覆盖原有数据
func (this *Cache) Import(data map[string]interface{}) {
	this.Lock()
	defer this.Unlock()
	for k, v := range data {
		vl := new(Data)
		vl.Value = v
		this.data[k] = vl
	}
}

//读取文件数据，并映射到缓存
//filename=文件名
func (this *Cache) ReadFile(filename string) error {
	this.Lock()
	defer this.Unlock()
	bt, e := ioutil.ReadFile(filename)
	if e != nil {
		return e
	}
	bt = deGzip(bt)
	bt = easyEncode(bt)
	readData := make([]*Data, 0)
	e2 := json.Unmarshal(bt, &readData)
	if e2 != nil {
		return e2
	}
	for _, d := range readData {
		this.data[d.Key] = d
	}
	return nil
}

//保存到外部文件
//file=文件名
func (this *Cache) SaveFile(filename string) error {
	this.Lock()
	defer this.Unlock()
	dumpData := make([]*Data, 0)
	for _, d := range this.data {
		dumpData = append(dumpData, d)
	}
	bt, e := json.Marshal(dumpData)
	if e != nil {
		return e
	}
	bt = easyEncode(bt)
	return ioutil.WriteFile(filename, enGzip(bt), 666)
}

//延时，（对所有的数据进行延时）
func (this *Cache) Expire(exp int64) {
	this.Lock()
	defer this.Unlock()
	for _, data := range this.data {
		data.Timeout = time.Now().Unix() + exp
	}
}

//添加数据前回调，如返回err则添加失败
func (this *Cache) OnSet(call CallFunc) {
	this.onSetCall = call
}

//保存缓存
//key=键名，value=值，expire=生存时间，生存时间=0时永久存储
//如果当前的大小已经大于初始化的空间，则自动释放部分空间
func (this *Cache) Set(key string, value interface{}, exp ...int64) error {
	if this.size > 0 && len(this.data) > this.size {
		this.checkExpired()
	}
	this.Lock()
	defer this.Unlock()
	var data *Data
	if d, ok := this.data[key]; ok {
		data = d
	} else {
		data = new(Data)
		this.index++
		data.Index = this.index
		data.Key = key
	}
	data.Value = value
	if len(exp) == 1 {
		data.Timeout = time.Now().Unix() + exp[0]
	}
	if this.onSetCall != nil {
		if err := this.onSetCall(key, value); err != nil {
			return err
		}
	}
	this.data[key] = data
	return nil
}

//读取缓存
//key=键名,返回值必须是未到期的，或者到期时间为0的
func (this *Cache) Get(key string) interface{} {
	this.Lock()
	defer this.Unlock()
	if data, ok := this.data[key]; ok && !this.isExpired(data) {
		return data.Value
	}
	delete(this.data, key)
	return nil
}

//获取值，并自动转换赋值给value
func (this *Cache) GetTo(key string, value interface{}) {
	v := this.Get(key)
	if v == nil {
		return
	}
	data, err := json.Marshal(v)
	if err != nil {
		return
	}
	json.Unmarshal(data, value)
}

//删除数据前回调,如果回调返回err,则不能删除
func (this *Cache) OnDel(call CallFunc) {
	this.onDelCall = call
}

//删除缓存,如果key存在
func (this *Cache) Del(key string) error {
	value := this.Get(key)
	this.Lock()
	defer this.Unlock()
	if this.onDelCall != nil {
		if err := this.onDelCall(key, value); err != nil {
			return err
		}
	}
	delete(this.data, key)
	return nil
}

//清空
func (this *Cache) Clean() {
	this.Lock()
	defer this.Unlock()
	this.data = make(map[string]*Data, this.size)
	this.maked = true
	this.index = 0
}

//判断是否到期，true到期，false未到期
func (this *Cache) isExpired(data *Data) bool {
	if data.Timeout == 0 || data.Timeout > time.Now().Unix() {
		return false
	}
	return true
}

//遍历所有数据，检查是否有到期的
func (this *Cache) checkExpired() {
	this.Lock()
	defer this.Unlock()
	for key, data := range this.data {
		if this.isExpired(data) {
			delete(this.data, key)
		}
	}
}

//缓存数据数量
func (this *Cache) Len() int {
	this.checkExpired()
	return len(this.data)
}

//有序的迭代<只迭代未到期的缓存数据>
func (this *Cache) Each(callback func(key, value interface{})) {
	this.Lock()
	indexArray := make([]*Data, this.index+1)
	for _, data := range this.data {
		if !this.isExpired(data) { //未到期
			indexArray[data.Index] = data
		} else {
			delete(this.data, data.Key)
		}
	}
	this.Unlock()
	for _, data := range indexArray {
		if data != nil {
			callback(data.Key, data.Value)
		}
	}
}
