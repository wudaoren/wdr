package orm

type Trans struct {
	Engine
}

//事务包裹
func TransWarp(fn func(*Trans) error) error {
	trans := new(Trans)
	return trans.Transaction(func(sx *Session) error {
		return fn(trans)
	})
}
