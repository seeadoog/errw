package errw

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

type Error interface {
	error
	Caller
	cause
	Wrap(args ...interface{}) Error
	Wrapf(f string, args ...interface{}) Error
	Wraps(s string) Error
	UnWrap() error
	WithCaller() Error
	WithCallers(callerSkip int) Error
	WithMeta(key string, val any) Error
	GetMeta(key string) (val any)
	Metas() map[string]any
}

type Caller interface {
	Callers() Callers
}

type Callers []string

func (c Callers) String() string {
	return strings.Join([]string(c), "\n")
}

type cause interface {
	cause() error
}

type werror struct {
	err      error
	msg      []string
	stack    Callers
	metas    map[string]any
	reversed bool
}

func (w *werror) Wrap(args ...interface{}) Error {
	w.msg = append(w.msg, fmt.Sprint(args...))
	return w
	// return Wrap(w, args...)
}

func (w *werror) Wrapf(f string, args ...interface{}) Error {
	w.msg = append(w.msg, fmt.Sprintf(f, args...))
	return w
	// return Wrapf(w, f, args...)

}

func (w *werror) Wraps(s string) Error {
	w.msg = append(w.msg, s)
	return w
}

func (w *werror) UnWrap() error {
	return Unwrap(w)
}

func (w *werror) Error() string {
	if len(w.msg) == 0 {
		return w.err.Error()
	}
	if !w.reversed {
		reverse(w.msg)
		w.reversed = true
	}
	return strings.Join(append(w.msg, w.err.Error()), " -> ")
}

func (w *werror) cause() error {
	return w.err
}

func (w *werror) WithCallers(callerSkip int) Error {
	w.stack = append(w.stack, getCall(callerSkip+1))
	return w
}
func (w *werror) WithCaller() Error {
	return w.WithCallers(1)
}

func (w *werror) WithMeta(key string, val any) Error {
	if w.metas == nil {
		w.metas = map[string]any{}
	}
	w.metas[key] = val

	return w
}

func (w *werror) GetMeta(key string) (val any) {
	if w.metas != nil {
		val = w.metas[key]
		if val != nil {
			return val
		}
	}
	e, ok := w.err.(*werror)
	if ok {
		return e.GetMeta(key)
	}
	return nil
}

func (w *werror) Metas() map[string]any {
	m := w.metas
	pe, ok := w.err.(*werror)
	if ok {
		pm := pe.Metas()
		if pm == nil {
			pm = map[string]any{}
		}
		for k, v := range m {
			pm[k] = v
		}
		return pm
	}
	return m
}

func getCall(callSkip int) string {
	pc, file, line, ok := runtime.Caller(callSkip + 1)
	if ok {
		fc := runtime.FuncForPC(pc)
		funName := ""
		if fc != nil {
			funName = fc.Name()
		}
		return file + ":" + strconv.Itoa(line) + " " + funName
	}
	return ""
}

func (w *werror) Callers() Callers {
	c, ok := w.err.(Caller)
	if ok {
		return append(c.Callers(), w.stack...)
	}
	return w.stack
}

func Wrap(err error, args ...interface{}) Error {

	return &werror{
		err: err,
		msg: []string{fmt.Sprint(args...)},
	}
}

func New(s string) Error {
	return &werror{
		err: errors.New(s),
		msg: make([]string, 0, 3),
	}
}

func NewWithMsgInitCap(s string, cap int) Error {
	return &werror{
		err: errors.New(s),
		msg: make([]string, 0, cap),
	}
}

func Wrapf(err error, f string, args ...interface{}) Error {
	return &werror{
		err: err,
		msg: []string{fmt.Sprintf(f, args...)},
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

func CallersOf(err error) Callers {
	e, ok := err.(Caller)
	if ok {
		return e.Callers()
	}
	return nil
}

func reverse(s []string) {
	i := 0
	j := len(s) - 1

	for i < j {
		s[i], s[j] = s[j], s[i]
		i++
		j--
	}
}
