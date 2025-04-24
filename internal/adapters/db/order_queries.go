package db

import (
	"context"
	"fmt"
	"logistic-app/internal/app/domain"
	"logistic-app/internal/common/errors"
	"sync"
	"time"
)

func (p *Postgres) GetOrder(ctx context.Context, orderID, userID uint) (*domain.Order, *errors.AppError) {
	var order *domain.Order
	result := p.db.WithContext(ctx).First(&order, orderID)
	if result.Error == nil && (order.SenderID != userID && order.ReceiverID != userID) {
		return nil, nil
	}
	return order, errors.ConvertGormErrors(result.Error)
}

func (p *Postgres) GetOrderWithForeignObjects(ctx context.Context, orderID, senderID uint) (*domain.Order, *errors.AppError) {
	order, err := p.GetOrder(ctx, orderID, senderID)
	if err != nil {
		return nil, err
	}
	errChan := make(chan *errors.AppError, 3)
	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		provider, e := p.GetProvider(ctx, order.ProviderID)
		if e != nil {
			errChan <- e
			return
		}
		order.Provider = provider
	}()

	go func() {
		defer wg.Done()
		sender, e := p.GetCustomer(ctx, order.SenderID)
		if e != nil {
			errChan <- e
			return
		}
		order.Sender = sender
	}()

	go func() {
		defer wg.Done()
		receiver, e := p.GetCustomer(ctx, order.ReceiverID)
		if e != nil {
			errChan <- e
			return
		}
		order.Receiver = receiver
	}()

	wg.Wait()
	close(errChan)

	for e := range errChan {
		if e != nil {
			return nil, e
		}
	}
	return order, nil
}

func (p *Postgres) CreateOrder(ctx context.Context, userID, receiverID, providerID uint, product *string) (*domain.Order, *errors.AppError) {
	order := &domain.Order{
		ProviderID: providerID,
		SenderID:   userID,
		ReceiverID: receiverID,
		Product:    product,
	}
	result := p.db.WithContext(ctx).Create(&order)
	return order, errors.ConvertGormErrors(result.Error)
}

func (p *Postgres) GetOngoingOrders(ctx context.Context) ([]*domain.Order, *errors.AppError) {
	var orders []*domain.Order
	result := p.db.WithContext(ctx).Where("status IN ?", domain.GetOngoingOrderStatus()).Find(&orders)
	return orders, errors.ConvertGormErrors(result.Error)
}

func (p *Postgres) UpdateOrderStatus(ctx context.Context, orderID uint, status string) (*domain.Order, *errors.AppError) {
	var order *domain.Order
	result := p.db.WithContext(ctx).First(&order, orderID)
	if result.Error != nil {
		return nil, errors.ConvertGormErrors(result.Error)
	}
	if order.Status == status {
		return order, nil
	}

	result = p.db.WithContext(ctx).WithContext(ctx).Model(&order).
		Where(domain.Order{ID: orderID}).
		Updates(domain.Order{Status: status})

	if status == domain.GetOrderStatus().PickedUp {
		// used these two lines to simulate
		//pastDays := 2 + rand.Intn(3)
		//pickedUpDate := time.Now().AddDate(0, 0, -1*pastDays)
		pickedUpDate := time.Now()
		result = result.Updates(domain.Order{PickedUpDate: &pickedUpDate})
		order.PickedUpDate = &pickedUpDate
	} else if status == domain.GetOrderStatus().Delivered {
		// used these two lines to simulate
		//pastDays := rand.Intn(3)
		//deliveryDate := time.Now().AddDate(0, 0, -1*pastDays)
		deliveryDate := time.Now()
		result = result.Updates(domain.Order{DeliveryDate: &deliveryDate})
		order.DeliveryDate = &deliveryDate
	}

	order.Status = status
	return order, errors.ConvertGormErrors(result.Error)
}

func (p *Postgres) UpdateOrderNotification(ctx context.Context, orderID uint) *errors.AppError {
	var order *domain.Order
	result := p.db.WithContext(ctx).WithContext(ctx).Model(&order).
		Where(domain.Order{ID: orderID}).
		Updates(domain.Order{NotifiedReceiver: true})
	return errors.ConvertGormErrors(result.Error)
}

func (p *Postgres) GetProvidersMeanDeliveryTime(ctx context.Context) ([]*domain.ProviderByDeliveryTime, *errors.AppError) {
	var data []*domain.ProviderByDeliveryTime

	result := p.db.WithContext(ctx).Raw(fmt.Sprintf(`
		select avg(delivery_date - picked_up_date) as mean_delivery_time_in_days, provider_id  
			from orders 
			where delivery_date is not null 
			and created_at >= NOW() - INTERVAL '7 days' 
			group by provider_id 
			order by mean_delivery_time_in_days desc
	`)).Scan(&data)

	return data, errors.ConvertGormErrors(result.Error)
}
