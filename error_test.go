package errw

import (
	"fmt"
	"testing"
)

type error2 = Error

func newRequest() error2 {
	return NewWithMsgInitCap("new request error", 5).WithMeta("2", 1)
}

func get() error2 {
	err := newRequest()
	return err.Wraps("new request error").WithMeta("1", "1").WithMeta("2", "3")
}

func doIt() error2 {
	err := get()
	return Wrap(err, "doit: get error").WithMeta("1", "5").Wraps("ggg")
}

func TestError(t *testing.T) {
	err := doIt()
	re := err.UnWrap()
	fmt.Println(err.Error(), "caller:\n", err.Callers(), "info", "raw error", re)
	fmt.Println(err.Error(), "caller:\n", err.Callers(), "info", "raw error", re)
	fmt.Println(err.GetMeta("code"), err.GetMeta("message"))

	fmt.Println(err.Metas())
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
