package service

import (
	"context"
	"logistic-app/internal/app/domain"
)

func (s *LogisticService) NotifyReceiver(order *domain.Order) {
	if order.NotifiedReceiver {
		return
	} else {
		s.repo.UpdateOrderNotification(context.Background(), order.ID)
		// todo: send a message to user
	}
}
