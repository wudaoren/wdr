package array

type Float []float64

//求数组最大值
func (list Float) Max() float64 {
	var max = list[0]
	for _, v := range list {
		if v > max {
			max = v
		}
	}
	return max
}

//求数组最小值
func (list Float) Min() float64 {
	var min = list[0]
	for _, v := range list {
		if v < min {
			min = v
		}
	}
	return min
}

//求数组和
func (list Float) Sum() float64 {
	var sum float64
	for _, v := range list {
		sum = sum + v
	}
	return sum
}

func (list Float) In(data float64) bool {
	for _, v := range list {
		if data == v {
			return true
		}
	}
	return false
}
