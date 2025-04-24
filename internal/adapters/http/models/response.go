package models

import (
	"context"
	"encoding/json"
	"log"
	"logistic-app/internal/common/configs"
	"logistic-app/internal/common/errors"
	"net/http"
)

type Response struct {
	Result any `json:"result"`
	Code   int `json:"code"`
}

type ErrorFunc func() *errors.AppError

func ReturnResp(ctx context.Context, v any) *Response {
	return &Response{
		Result: v,
		Code:   http.StatusOK,
	}
}

func ReturnErrorResp(ctx context.Context, err *errors.AppError) *Response {
	if configs.LogError {
		log.Println(err.Err)
	}
	return &Response{Result: err.ApiErr, Code: err.Code}
}

func WriteJSON(w http.ResponseWriter, resp *Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Code)
	err := json.NewEncoder(w).Encode(resp.Result)
	if err != nil {
		_ = json.NewEncoder(w).Encode(errors.InternalServerError(err).ApiErr)
	}
}
