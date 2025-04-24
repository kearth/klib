package main

import (
	"errors"
	"fmt"
)

// check
var _ IError = new(Error)

// IError
type IError interface {
	error
	Code() int
	Wrap(err error) IError
	Unwrap() error
}

// New
func New(c int, m string) *Error {
	return &Error{code: c, msg: m}
}

// Error
type Error struct {
	code int
	msg  string
}

// Error
func (e *Error) Error() string {
	return e.msg
}

// Code
func (e *Error) Code() int {
	return e.code
}

// Wrap
func (e *Error) Wrap(err error) IError {
	if err != nil {
		e.msg = fmt.Sprintf("%s:%s", e.msg, err.Error())
	}
	return e
}

// Unwrap
func (e *Error) Unwrap() error {
	return errors.New(e.Error())
}
