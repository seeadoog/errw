package errw

import "fmt"

type Error interface {
	error
	cause
	Wrap(args ...interface{}) Error
	Wrapf(f string, args ...interface{}) Error
	UnWrap() error
}

type cause interface {
	cause() error
}

type werror struct {
	err error
	msg string
}

func (w *werror) Wrap(args ...interface{}) Error {
	return Wrap(w, args...)
}

func (w *werror) Wrapf(f string, args ...interface{}) Error {
	return Wrapf(w, f, args...)
}

func (w *werror) UnWrap() error {
	return Unwrap(w)
}

func (w *werror) Error() string {
	return w.msg + " -> " + w.err.Error()
}

func (w *werror) cause() error {
	return w.err
}

func Wrap(err error, args ...interface{}) Error {
	return &werror{
		err: err,
		msg: fmt.Sprint(args...),
	}
}

func Wrapf(err error, f string, args ...interface{}) Error {
	return &werror{
		err: err,
		msg: fmt.Sprintf(f, args...),
	}
}

func Unwrap(err error) error {
	var e cause
	var ok bool
	for {
		e, ok = err.(cause)
		if ok {
			err = e.cause()
			continue
		}
		return err
	}
}
