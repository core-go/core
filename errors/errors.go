package errors

import (
	"fmt"
	"runtime"
)

type Error struct {
	Message string      `json:"message,omitempty"`
	Code    int         `json:"-"`
	caller  string      `json:"-"`
	payload interface{} `json:"payload"`
	error   `json:"-`
}

func New(msg string) *Error {
	_, f, l, _ := runtime.Caller(1)
	caller := fmt.Sprintf("%v:%v", f, l)
	return &Error{
		Message: msg,
		caller:  caller,
	}
}

func (e *Error) Error() string {
	return e.Message
}

func W(err error) *Error {
	if err == nil {
		return nil
	}
	_, f, l, _ := runtime.Caller(1)
	caller := fmt.Sprintf("%v:%v", f, l)
	return &Error{
		Message: err.Error(),
		caller:  caller,
		error:   err,
	}
}

func (e *Error) WithCode(code int) *Error {
	e.Code = code
	return e
}

func (e *Error) GetCaller() string {
	return e.caller
}

func (e *Error) Cause() error {
	return e.error
}
