package domain

import (
	"logistic-app/internal/common/errors"
	"net/http"
	"time"
)

type Provider struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"not null;unique"`
	Url       string    `json:"url" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

type ProviderCreateRequest struct {
	noPathReq
	Name string `json:"name" required:"true"`
	Url  string `json:"url" required:"true"`
}

func (pr *ProviderCreateRequest) UnmarshalBody(request *http.Request) *errors.AppError {
	return getBody(pr, request)
}

type ProviderByDeliveryTime struct {
	ProviderID             uint    `json:"provider_id"`
	MeanDeliveryTimeInDays float32 `json:"mean_delivery_time_in_days"`
}
