package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"logistic-app/internal/app/domain"
	"logistic-app/internal/app/ports"
	"logistic-app/internal/common/configs"
	"logistic-app/internal/common/errors"
	"time"
)

type LogisticService struct {
	repo ports.Repo
}

func NewLogisticService(repo ports.Repo) *LogisticService {
	return &LogisticService{repo: repo}
}

func (s *LogisticService) HealthCheck(ctx context.Context) (any, *errors.AppError) {
	return map[string]bool{
		"db_healthy": s.repo.Ready(),
	}, nil
}

func (s *LogisticService) GetProviders(ctx context.Context) ([]*domain.Provider, *errors.AppError) {
	return s.repo.GetAllProviders(ctx)
}

func (s *LogisticService) GetProvidersMeanDelTime(ctx context.Context) ([]*domain.ProviderByDeliveryTime, *errors.AppError) {
	return s.repo.GetProvidersMeanDeliveryTime(ctx)
}

func (s *LogisticService) CreateProvider(ctx context.Context, request *domain.ProviderCreateRequest) (*domain.Provider, *errors.AppError) {
	return s.repo.CreateProvider(ctx, &request.Name, &request.Url)
}

func (s *LogisticService) CreateCustomer(ctx context.Context, request *domain.CustomerCreateRequest) (*domain.Customer, *errors.AppError) {
	return s.repo.CreateCustomer(ctx, request.Name, &request.PhoneNumber, &request.Address, &request.PostalCode)
}

func (s *LogisticService) GetCustomerToken(ctx context.Context, request *domain.CustomerTokenRequest) (any, *errors.AppError) {
	var byteSecKey = []byte(configs.SecretKey)

	claims := jwt.MapClaims{
		configs.JWTDefaults["USER_ID_CLAIM"].(string): request.ID,
		"exp": time.Now().Add(configs.TokenExpiration).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, e := token.SignedString(byteSecKey)
	if e != nil {
		return nil, errors.InternalServerError(e)
	}
	return map[string]string{"token": signed}, nil
}

func (s *LogisticService) CreateOrder(ctx context.Context, request *domain.OrderCreateRequest) (*domain.Order, *errors.AppError) {
	userID, ok := ctx.Value(configs.UserIDKey).(uint)
	if !ok {
		return nil, errors.NotFoundError(fmt.Errorf("user uuid not found in context"))
	}
	customer, err := s.repo.GetCustomer(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.repo.CreateOrder(ctx, customer.ID, request.ReceiverID, request.ProviderID, request.Product)
}

func (s *LogisticService) GetOrder(ctx context.Context, request *domain.OrderGetRequest) (*domain.Order, *errors.AppError) {
	userID, ok := ctx.Value(configs.UserIDKey).(uint)
	if !ok {
		return nil, errors.NotFoundError(fmt.Errorf("user uuid not found in context"))
	}
	return s.repo.GetOrderWithForeignObjects(ctx, request.OrderID, userID)
}
