package domain

import (
	"logistic-app/internal/common/errors"
	"net/http"
	"time"
)

type Customer struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PhoneNumber string    `json:"phone_number" gorm:"unique:not null"`
	Name        *string   `json:"name"`
	Address     string    `json:"address" gorm:"not null"`
	PostalCode  string    `json:"postal_code" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null"`
}

type CustomerCreateRequest struct {
	noPathReq
	PhoneNumber string  `json:"phone_number" required:"true"`
	Name        *string `json:"name"`
	Address     string  `json:"address" required:"true"`
	PostalCode  string  `json:"postal_code" required:"true"`
}

func (cr *CustomerCreateRequest) UnmarshalBody(request *http.Request) *errors.AppError {
	return getBody(cr, request)
}

type CustomerTokenRequest struct {
	noPathReq
	ID uint `json:"id" required:"true"`
}

func (cr *CustomerTokenRequest) UnmarshalBody(request *http.Request) *errors.AppError {
	return getBody(cr, request)
}
