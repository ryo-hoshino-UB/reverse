package xerrors

import (
	"fmt"
	"net/http"

	"github.com/go-errors/errors"
)

func Wrap(errp *error, format string, args ...any) {
	if *errp != nil {
		*errp = errors.WrapPrefix(*errp, fmt.Sprintf(format, args...), 1)
	}
}

var (
	ErrNotFound = NewStatusCodeError(http.StatusNotFound, "not found")
	ErrInternal = NewStatusCodeError(http.StatusInternalServerError, "internal server error")
	ErrBadRequest = NewStatusCodeError(http.StatusBadRequest, "bad request")
)

type StatusCodeError struct {
	msg string
	code int
	cause error
}

func NewStatusCodeError(code int, msg string) *StatusCodeError {
	return &StatusCodeError{
		msg: msg,
		code: code,
	}
}

func (e *StatusCodeError) Error() string {
	return e.msg
}

func (e *StatusCodeError) Unwrap() error {
	return e.cause
}

func (e *StatusCodeError) StatusCode() int {
	return e.code
}