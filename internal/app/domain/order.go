package domain

import (
	"logistic-app/internal/common/errors"
	"net/http"
	"time"
)

type Order struct {
	ID               uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ProviderID       uint       `json:"provider_id" gorm:"index;not null"`
	Provider         *Provider  `json:"provider,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SenderID         uint       `json:"sender_id" gorm:"index;not null"`
	Sender           *Customer  `json:"sender,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ReceiverID       uint       `json:"receiver_id" gorm:"index;not null"`
	Receiver         *Customer  `json:"receiver,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Product          *string    `json:"product"`
	Status           string     `json:"status" gorm:"size:15;not null;default:'PROVIDER_SEEN'"`
	PickedUpDate     *time.Time `json:"picked_up_date" gorm:"type:date"`
	DeliveryDate     *time.Time `json:"delivery_date" gorm:"type:date"`
	NotifiedReceiver bool       `json:"notified_receiver" gorm:"not null;default:false"`
	CreatedAt        time.Time  `json:"created_at" gorm:"not null;index"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"not null"`
}

type OrderStatus struct {
	Pending      string
	InProgress   string
	ProviderSeen string
	PickedUp     string
	Delivered    string
}

func GetOrderStatus() *OrderStatus {
	return &OrderStatus{
		Pending:      "PENDING",
		InProgress:   "IN_PROGRESS",
		ProviderSeen: "PROVIDER_SEEN",
		PickedUp:     "PICKED_UP",
		Delivered:    "DELIVERED",
	}
}

func GetOngoingOrderStatus() []string {
	return []string{
		GetOrderStatus().InProgress,
		GetOrderStatus().ProviderSeen,
		GetOrderStatus().PickedUp,
	}
}

func ConvertOrderStatus(no string) string {
	switch no {
	case "1":
		return GetOrderStatus().PickedUp
	case "2":
		return GetOrderStatus().InProgress
	case "3":
		return GetOrderStatus().Delivered
	}
	return ""
}

func ConvertOrderStatusToNumber(status string) int {
	switch status {
	case GetOrderStatus().PickedUp:
		return 1
	case GetOrderStatus().InProgress:
		return 2
	case GetOrderStatus().Delivered:
		return 3
	}
	return 0
}

type OrderStatusInResponse struct {
	StatusNumber string `json:"status"`
	FaStatus     string `json:"fa_status"`
	CreatedAt    string `json:"created_at"`
}

type ProviderUrlResponse struct {
	Message string                   `json:"message"`
	Data    []*OrderStatusInResponse `json:"data"`
}

type OrderCreateRequest struct {
	noPathReq
	ProviderID uint    `json:"provider_id" required:"true"`
	ReceiverID uint    `json:"receiver_id" required:"true"`
	Product    *string `json:"product_name"`
}

func (or *OrderCreateRequest) UnmarshalBody(request *http.Request) *errors.AppError {
	return getBody(or, request)
}

type OrderGetRequest struct {
	noBodyReq
	OrderID uint `json:"order_id"`
}

func (or *OrderGetRequest) UnmarshalPathValue(request *http.Request) *errors.AppError {
	return getPathValues(or, request)
}
