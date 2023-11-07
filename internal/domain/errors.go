package domain

import (
	"errors"
)

const (
	// Service error
	BadRequest   = "bad_request"
	Unauthorized = "unauthorized"
	Forbidden    = "forbidden"
	NotFound     = "not_found"
	ServerErr    = "server_err"
	Timeout      = "bad_gateway"

	// PQ Error
	PQErrDuplicate = "23505"
)

type Error struct {
	ErrCode string
	Message string
}

func NewError(errCode, msg string) error {
	return Error{
		ErrCode: errCode,
		Message: msg,
	}
}

func (e Error) Error() string {
	err := errors.New(e.Message)
	return err.Error()
}
