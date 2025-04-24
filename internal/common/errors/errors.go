package errors

import (
	"fmt"
	"gorm.io/gorm"
	"net/http"
)

type apiError struct {
	Msg string `json:"errors"`
}

type AppError struct {
	ApiErr *apiError
	Err    error
	Code   int
}

func Unauthorized() *AppError {
	return &AppError{
		ApiErr: &apiError{Msg: "Unauthorized"},
		Err:    fmt.Errorf("unauthorized"),
		Code:   http.StatusUnauthorized,
	}
}

func BadRequest(msg string) *AppError {
	return &AppError{
		ApiErr: &apiError{Msg: msg},
		Err:    fmt.Errorf(msg),
		Code:   http.StatusBadRequest,
	}
}

func InternalServerError(err error) *AppError {
	return &AppError{
		ApiErr: &apiError{Msg: err.Error()},
		Err:    err,
		Code:   http.StatusInternalServerError,
	}
}

func NotFoundError(err error) *AppError {
	return &AppError{
		ApiErr: &apiError{Msg: "object not found"},
		Err:    err,
		Code:   http.StatusNotFound,
	}
}

func ConvertGormErrors(err error) *AppError {
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return NotFoundError(err)
	default:
		return InternalServerError(err)
	}
}
