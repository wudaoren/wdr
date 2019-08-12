package img

import (
	"math"
	"testing"
)

//水印图片平铺算法测试
func TestSuanfa(t *testing.T) {
	wx, wy := 20, 10
	ix, iy := 60, 30
	x_num := int(math.Ceil(float64(ix) / float64(wx)))
	y_num := int(math.Ceil(float64(iy) / float64(wy)))
	for n := 0; n < y_num; n++ {
		for r := 0; r < x_num; r++ {
			t.Log("坐标：", r*wx, n*wy, (r+1)*wx, (n+1)*wy)
		}
		t.Log("换行")
	}
}
