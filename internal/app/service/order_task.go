package service

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"logistic-app/internal/app/domain"
	"logistic-app/internal/common/configs"
	"logistic-app/internal/common/errors"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var orderTaskRetries = 3
var orderTaskJobName = "update_orders_status"
var orderTaskIntervalInMinutes = int(configs.OrderUpdatePeriod.Minutes())

func (s *LogisticService) ScheduleUpdateOrderStatus() {
	ctx := context.Background()
	task, err := s.repo.GetOrCreatePeriodicTask(
		ctx, orderTaskJobName, orderTaskIntervalInMinutes)
	if err != nil {
		log.Fatal(err.Err)
	}

	var lastRun time.Time
	for {
		if task.LastRunTime == nil {
			go s.updateOrderWithRetries(orderTaskRetries)
			lastRun = time.Now()
		} else {
			lastRun = *task.LastRunTime
		}
		nextRun := lastRun.Add(configs.OrderUpdatePeriod)
		sleepDuration := nextRun.Sub(time.Now())
		time.Sleep(sleepDuration)
		if sleepDuration < 0 {
			nextRun = time.Now()
		}
		go s.updateOrderWithRetries(orderTaskRetries)
		task.LastRunTime = &nextRun
	}
}

func (s *LogisticService) updateOrderWithRetries(retries int) {
	ctx := context.Background()
	var (
		failed bool
		orders []*domain.Order
		err    *errors.AppError
		es     []error
		eStr   string
	)
	defer func() {
		log.Printf("Updated orders -> failed: %v, errors: %s", failed, eStr)
		s.repo.CreateOrUpdatePeriodicTask(ctx, orderTaskJobName, orderTaskIntervalInMinutes, failed, &eStr)
	}()
	orders, err = s.repo.GetOngoingOrders(ctx)
	if err != nil {
		failed = true
		eStr = err.Err.Error()
		return
	}

	for i := 0; i < retries; i++ {
		orders, es = s.updateOrdersStatusTask(orders)
		if len(orders) == 0 {
			break
		}
	}
	if len(es) > 0 {
		failed = true
		for _, e := range es {
			eStr += ";" + e.Error()
		}
	}
}

func (s *LogisticService) updateOrdersStatusTask(orders []*domain.Order) ([]*domain.Order, []error) {
	ctx := context.Background()
	sem := make(chan struct{}, configs.PeriodicTaskMaxConcurrency)
	var wg sync.WaitGroup
	var failedOrders []*domain.Order
	var errs []error
	var mu sync.Mutex

	for _, order := range orders {
		wg.Add(1)
		sem <- struct{}{}

		go func(o *domain.Order) {
			defer wg.Done()
			c, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			e := s.orderUpdateWorker(c, o)
			if e != nil {
				mu.Lock()
				failedOrders = append(failedOrders, o)
				errs = append(errs, e)
				mu.Unlock()
			}
			<-sem
		}(order)
	}

	wg.Wait()
	return failedOrders, errs
}

func (s *LogisticService) orderUpdateWorker(ctx context.Context, order *domain.Order) error {
	provider, err := s.repo.GetProvider(ctx, order.ProviderID)
	if provider == nil {
		return err.Err
	}

	resp, e := http.Get(provider.Url)
	if e != nil {
		return e
	}
	defer resp.Body.Close()

	dataBytes, e := io.ReadAll(resp.Body)
	if e != nil {
		return e
	}

	var data *domain.ProviderUrlResponse
	if e = json.Unmarshal(dataBytes, &data); e != nil {
		return e
	}

	// simulating real scenario
	source := domain.ConvertOrderStatusToNumber(order.Status)
	choice := rand.Intn(2)
	log.Println(source, choice)
	if source == 0 && choice == 0 {
		return nil
	}

	newStatus := domain.ConvertOrderStatus(data.Data[source-1+choice].StatusNumber)
	if newStatus == domain.GetOrderStatus().PickedUp {
		s.NotifyReceiver(order)
	}
	_, err = s.repo.UpdateOrderStatus(ctx, order.ID, newStatus)
	if err != nil {
		return err.Err
	}
	return nil
}
