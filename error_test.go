package errw

import (
	"fmt"
	"testing"
)

type error2 = Error

func newRequest() error2 {
	return NewWithMsgInitCap("new request error", 5)
}

func get() error2 {
	err := newRequest()
	return err.Wraps("new request error")
}

func doIt() error2 {
	err := get()
	return err.Wraps("doit: get error").Wraps("ggg")
}

func TestError(t *testing.T) {
	err := doIt()
	re := err.UnWrap()
	fmt.Println(err.Error(), "caller:\n", err.Callers(), "info", "raw error", re)
	fmt.Println(err.Error(), "caller:\n", err.Callers(), "info", "raw error", re)
	fmt.Println(err.GetMeta("code"), err.GetMeta("message"))
}

func stackOf() {
	defer func() {
		recover()

	}()

	panic("gg")
}

func BenchmarkName(b *testing.B) {
	// TODO: Initialize
	for i := 0; i < b.N; i++ {
		// TODO: Your Code Here
		// stackOf()
		doIt().Error()

	}
}

func BenchmarkPa(b *testing.B) {
	// TODO: Initialize
	e := 1
	for i := 0; i < b.N; i++ {
		assert(e)

	}
}

func assert(v any) {
	// _, _ = v.(Error)
}
