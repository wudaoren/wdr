package stcp

type Context struct {
	*Conn
	data []byte
	err  interface{}
}

func newContext(conn *Conn, data []byte, err interface{}) *Context {
	p := new(Context)
	p.Conn = conn
	p.data = data
	p.err = err
	return p
}

//
func newCloseContext(conn *Conn, e1 error, e2 interface{}) *Context {
	var err Error
	if e1 != nil {
		err = e1
	} else {
		err = e2
	}
	return newContext(conn, nil, err)
}

//
func (this *Context) Data() []byte {
	return this.data
}

//
func (this *Context) Error() Error {
	return this.err
}
