package array

type Int []int

//求数组最大值
func (list Int) Max() int {
	var max = list[0]
	for _, v := range list {
		if v > max {
			max = v
		}
	}
	return max
}

//求数组最小值
func (list Int) Min() int {
	var min = list[0]
	for _, v := range list {
		if v < min {
			min = v
		}
	}
	return min
}

//求数组和
func (list Int) Sum() int {
	var sum int
	for _, v := range list {
		sum = sum + v
	}
	return sum
}

func (list Int) In(data int) bool {
	for _, v := range list {
		if data == v {
			return true
		}
	}
	return false
}
