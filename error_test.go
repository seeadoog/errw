package errw

import (
	"errors"
	"fmt"
	"testing"
)

type error2 = Error

func newRequest() error {
	return errors.New("new request error")
}

func get() error2 {
	err := newRequest()
	return Wrap(err, "get")
}

func doIt() error2 {
	err := get()
	return err.Wrap("doIt")
}

func TestError(t *testing.T) {
	err := doIt()
	fmt.Println(err)
	fmt.Println(err.UnWrap())
}
