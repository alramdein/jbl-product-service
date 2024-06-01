package handler

import (
	"net/http"
	"referral-system/usecase"

	"github.com/labstack/echo/v4"
)

type CustomError struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
}

var ErrorMap = map[error]int{
	usecase.ErrEmailRequired:           http.StatusBadRequest,
	usecase.ErrPasswordRequired:        http.StatusBadRequest,
	usecase.ErrRoleRequired:            http.StatusBadRequest,
	usecase.ErrReferalCodeRequired:     http.StatusBadRequest,
	usecase.ErrEmailAlreadyExist:       http.StatusConflict,
	usecase.ErrReferralCodeIsNotExisit: http.StatusBadRequest,
	usecase.ErrInvalidRole:             http.StatusBadRequest,
	usecase.ErrInvalidEmail:            http.StatusBadRequest,
	usecase.ErrRoleNotFound:            http.StatusNotFound,
}

func MapErrorToHTTPResponse(err error) *echo.HTTPError {
	httpError, ok := ErrorMap[err]
	if ok && httpError != http.StatusInternalServerError {
		return echo.NewHTTPError(httpError, CustomError{
			StatusCode: httpError,
			Message:    err.Error(),
		})
	}
	return echo.NewHTTPError(http.StatusInternalServerError, "something went wrong")
}
