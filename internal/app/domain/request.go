package domain

import (
	"encoding/json"
	"fmt"
	"io"
	"logistic-app/internal/common/errors"
	"net/http"
	"reflect"
	"strconv"
)

type Request interface {
	UnmarshalBody(request *http.Request) *errors.AppError
	UnmarshalPathValue(request *http.Request) *errors.AppError
}

func getLimitNOffset(r *http.Request) (uint64, uint64) {
	return getCustomLimitNOffset(r, "limit", "offset", 10, 0)
}

func getCustomLimitNOffset(r *http.Request, limitKey, offsetKey string, defLimit, defOffset uint64) (uint64, uint64) {
	limitStr := r.URL.Query().Get(limitKey)
	offsetStr := r.URL.Query().Get(offsetKey)
	limit, errConv := strconv.ParseUint(limitStr, 10, 64)
	if errConv != nil {
		limit = defLimit
	}
	offset, errConv := strconv.ParseUint(offsetStr, 10, 64)
	if errConv != nil {
		offset = defOffset
	}
	return limit, offset
}

func getPathValues(r any, request *http.Request) *errors.AppError {
	instType := reflect.TypeOf(r).Elem()
	v := reflect.ValueOf(r).Elem()

	for i := 0; i < instType.NumField(); i++ {
		field := instType.Field(i)
		jsonTag := field.Tag.Get("json")

		value := request.PathValue(jsonTag)
		if fieldVal := v.FieldByName(field.Name); fieldVal.IsValid() && fieldVal.CanSet() {
			if fieldVal.Kind() == reflect.String {
				fieldVal.SetString(value)
			}
			if fieldVal.Kind() == reflect.Uint {
				n, e := strconv.ParseUint(value, 10, 64)
				if e != nil {
					return errors.InternalServerError(e)
				}
				fieldVal.SetUint(n)
			}
		}
	}
	return nil
}

func getBody(r any, request *http.Request) *errors.AppError {
	b, err := io.ReadAll(request.Body)
	if err != nil {
		return errors.InternalServerError(err)
	}
	if len(b) == 0 {
		errMsg := "body is missing"
		return errors.BadRequest(errMsg)
	}

	if err = json.Unmarshal(b, &r); err != nil {
		return errors.InternalServerError(err)
	}

	instType := reflect.TypeOf(r).Elem()
	v := reflect.ValueOf(r).Elem()

	for i := 0; i < instType.NumField(); i++ {
		field := instType.Field(i)
		fieldVal := v.FieldByName(field.Name)
		required := field.Tag.Get("required")

		if required == "true" {
			if fieldVal.Kind() == reflect.String && fieldVal.String() == "" {
				return errors.BadRequest(field.Name + " is required")
			} else if fieldVal.Kind() == reflect.Uint && fieldVal.Uint() == 0 {
				return errors.BadRequest(field.Name + " is required")
			}
		}
	}

	return nil
}

type noBodyReq struct{}

func (n *noBodyReq) UnmarshalBody(request *http.Request) *errors.AppError {
	err := fmt.Errorf("cannot accept body, %s", request.RequestURI)
	return errors.InternalServerError(err)
}

type noPathReq struct{}

func (n *noPathReq) UnmarshalPathValue(request *http.Request) *errors.AppError {
	err := fmt.Errorf("cannot accept path values, %s", request.RequestURI)
	return errors.InternalServerError(err)
}
