package domain

import (
	"errors"
	"net/http"
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

func MapServiceErrHTTPResponse(err error) (int, map[string]any) {
	var svcErr Error
	if errors.As(err, &svcErr) {
		resp := map[string]any{
			"error": svcErr.Message,
		}
		switch svcErr.ErrCode {
		case BadRequest:
			return http.StatusBadRequest, resp
		case Unauthorized:
			return http.StatusUnauthorized, resp
		case Forbidden:
			return http.StatusForbidden, resp
		case NotFound:
			return http.StatusNotFound, resp
		case Timeout:
			return http.StatusBadGateway, resp
		default:
			return http.StatusInternalServerError, resp
		}
	}
	return http.StatusInternalServerError, map[string]any{
		"error": "internal server error",
	}
}
