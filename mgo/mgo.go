package mongo

import (
	"math"

	"github.com/globalsign/mgo"
)

//mongodb分页查询
//分页对象
type Page struct {
	Total     int
	Limit     int
	Page      int
	DataTotal int
}

//分页查询
//参数：query=查询参数，obj=查询对象，current=当前页码，limit=每一页数量
//返回：p=分页数据，start=sql查询limit起点

func QueryPage(query *mgo.Query, objList interface{}, page, limit int) (p *Page, err error) {
	if limit == 0 {
		limit = 20
	}
	totle, err := query.Count()
	if err != nil {
		return
	}
	p = new(Page)
	p.Total = int(math.Ceil(float64(totle) / float64(limit))) //计算出页码
	if page >= p.Total {                                      //页末不越界
		page = p.Total
	}
	if page <= 0 {
		page = 1
	}
	p.DataTotal = int(totle)
	p.Limit = limit
	p.Page = page
	start := (page - 1) * limit
	err = query.Limit(limit).Skip(start).All(objList)
	return p, err
}
