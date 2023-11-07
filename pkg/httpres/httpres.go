package httpres

import (
	"errors"
	"net/http"

	"github.com/ryanadiputraa/spendr-backend/internal/domain"
)

func MapServiceErrHTTPResponse(err error) (int, map[string]any) {
	var svcErr domain.Error
	if errors.As(err, &svcErr) {
		resp := map[string]any{
			"error": svcErr.Message,
		}
		switch svcErr.ErrCode {
		case domain.BadRequest:
			return http.StatusBadRequest, resp
		case domain.Unauthorized:
			return http.StatusUnauthorized, resp
		case domain.Forbidden:
			return http.StatusForbidden, resp
		case domain.NotFound:
			return http.StatusNotFound, resp
		case domain.Timeout:
			return http.StatusBadGateway, resp
		default:
			return http.StatusInternalServerError, resp
		}
	}
	return http.StatusInternalServerError, map[string]any{
		"error": "internal server error",
	}
}
