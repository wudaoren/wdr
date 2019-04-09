package errs

import (
	"fmt"
)

type FatalError struct {
	err string
}

func (this *FatalError) Error() string {
	return this.err
}

func Fatal(args ...interface{}) error {
	return &FatalError{fmt.Sprint(args...)}
}

func Fatalf(fomart string, args ...interface{}) error {
	return &FatalError{fmt.Sprintf(fomart, args...)}
}
