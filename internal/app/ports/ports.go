package ports

import (
	"context"
	"logistic-app/internal/app/domain"
	"logistic-app/internal/common/errors"
)

type Service interface {
	HealthCheck(ctx context.Context) (any, *errors.AppError)

	GetProviders(ctx context.Context) ([]*domain.Provider, *errors.AppError)
	GetProvidersMeanDelTime(ctx context.Context) ([]*domain.ProviderByDeliveryTime, *errors.AppError)
	CreateProvider(ctx context.Context, request *domain.ProviderCreateRequest) (*domain.Provider, *errors.AppError)

	CreateCustomer(ctx context.Context, request *domain.CustomerCreateRequest) (*domain.Customer, *errors.AppError)
	GetCustomerToken(ctx context.Context, request *domain.CustomerTokenRequest) (any, *errors.AppError)

	CreateOrder(ctx context.Context, request *domain.OrderCreateRequest) (*domain.Order, *errors.AppError)
	GetOrder(ctx context.Context, request *domain.OrderGetRequest) (*domain.Order, *errors.AppError)

	ScheduleUpdateOrderStatus()
}

type Repo interface {
	Ready() bool
	Close()

	GetOrder(ctx context.Context, orderID, senderID uint) (*domain.Order, *errors.AppError)
	GetOrderWithForeignObjects(ctx context.Context, orderID, senderID uint) (*domain.Order, *errors.AppError)
	CreateOrder(ctx context.Context, userID, receiverID, providerID uint, product *string) (*domain.Order, *errors.AppError)
	GetOngoingOrders(ctx context.Context) ([]*domain.Order, *errors.AppError)
	UpdateOrderStatus(ctx context.Context, orderID uint, status string) (*domain.Order, *errors.AppError)
	UpdateOrderNotification(ctx context.Context, orderID uint) *errors.AppError
	GetProvidersMeanDeliveryTime(ctx context.Context) ([]*domain.ProviderByDeliveryTime, *errors.AppError)

	GetCustomer(ctx context.Context, userID uint) (*domain.Customer, *errors.AppError)
	CreateCustomer(ctx context.Context, name, phone, addr, postalCode *string) (*domain.Customer, *errors.AppError)

	GetProvider(ctx context.Context, providerID uint) (*domain.Provider, *errors.AppError)
	GetAllProviders(ctx context.Context) ([]*domain.Provider, *errors.AppError)
	CreateProvider(ctx context.Context, name, url *string) (*domain.Provider, *errors.AppError)

	CreateOrUpdatePeriodicTask(ctx context.Context, name string, interval int, failed bool, e *string) (*domain.PeriodicTask, *errors.AppError)
	GetOrCreatePeriodicTask(ctx context.Context, name string, interval int) (*domain.PeriodicTask, *errors.AppError)
}
