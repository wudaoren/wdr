package errs

import (
	"fmt"
)

type WarnError struct {
	err string
}

func (this *WarnError) Error() string {
	return this.err
}

func Warn(args ...interface{}) error {
	return &WarnError{fmt.Sprint(args...)}
}

func Warnf(fomart string, args ...interface{}) error {
	return &WarnError{fmt.Sprintf(fomart, args...)}
}
