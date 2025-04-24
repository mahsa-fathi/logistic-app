package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"logistic-app/internal/app/domain"
	"net/http"
	"testing"
)

func setUpOrderForeignObjects(t *testing.T) (sender *domain.Customer, receiver *domain.Customer, provider *domain.Provider) {
	phone := "09"
	name := "name"
	phone2 := "20"
	address := "somewhere"
	postal := "some-code"
	test := "test-provider"
	sender, err := repo.CreateCustomer(context.Background(), &name, &phone, &address, &postal)
	if err != nil {
		t.Error(err.Err)
	}
	receiver, err = repo.CreateCustomer(context.Background(), &name, &phone2, &address, &postal)
	if err != nil {
		t.Error(err.Err)
	}
	provider, err = repo.CreateProvider(context.Background(), &test, &test)
	if err != nil {
		t.Error(err.Err)
	}
	return
}

func TestPostgres_CreateOrder(t *testing.T) {
	tearUpSuite := setupSuite()
	defer tearUpSuite()

	sender, receiver, provider := setUpOrderForeignObjects(t)

	t.Run("successful create", func(t *testing.T) {
		product := "phone"
		order, err := repo.CreateOrder(context.Background(), sender.ID, receiver.ID, provider.ID, &product)
		assert.Empty(t, err)
		assert.Equal(t, product, *order.Product)
		assert.Equal(t, sender.ID, order.SenderID)
		assert.Equal(t, receiver.ID, order.ReceiverID)
		assert.Equal(t, provider.ID, order.ProviderID)
	})

	t.Run("successful create with no name", func(t *testing.T) {
		order, err := repo.CreateOrder(context.Background(), sender.ID, receiver.ID, provider.ID, nil)
		assert.Empty(t, err)
		assert.Empty(t, order.Product)
		assert.Equal(t, sender.ID, order.SenderID)
		assert.Equal(t, receiver.ID, order.ReceiverID)
		assert.Equal(t, provider.ID, order.ProviderID)
	})
}

func TestPostgres_GetOrder(t *testing.T) {
	tearUpSuite := setupSuite()
	defer tearUpSuite()

	sender, receiver, provider := setUpOrderForeignObjects(t)

	t.Run("successful get", func(t *testing.T) {
		product := "phone"
		act, err := repo.CreateOrder(context.Background(), sender.ID, receiver.ID, provider.ID, &product)
		assert.Empty(t, err)

		order, err := repo.GetOrder(context.Background(), act.ID, sender.ID)
		assert.Empty(t, err)
		assert.Equal(t, product, *order.Product)
		assert.Equal(t, sender.ID, order.SenderID)
		assert.Equal(t, receiver.ID, order.ReceiverID)
		assert.Equal(t, provider.ID, order.ProviderID)

		order, err = repo.GetOrder(context.Background(), act.ID, receiver.ID)
		assert.Empty(t, err)
		assert.Equal(t, act.ID, order.ID)
	})

	t.Run("unsuccessful get", func(t *testing.T) {
		_, err := repo.GetOrder(context.Background(), 1000, 1000)
		assert.NotEmpty(t, err)
		assert.Equal(t, http.StatusNotFound, err.Code)
	})
}

func TestPostgres_UpdateOrderStatus(t *testing.T) {
	tearUpSuite := setupSuite()
	defer tearUpSuite()

	sender, receiver, provider := setUpOrderForeignObjects(t)
	order, err := repo.CreateOrder(context.Background(), sender.ID, receiver.ID, provider.ID, nil)
	if err != nil {
		t.Error(err.Err)
	}
	assert.Empty(t, order.PickedUpDate)
	assert.Empty(t, order.DeliveryDate)

	t.Run("successful update to PROVIDER_SEEN", func(t *testing.T) {
		order, err = repo.UpdateOrderStatus(context.Background(), order.ID, domain.GetOrderStatus().ProviderSeen)
		assert.Empty(t, err)
		order, err = repo.GetOrder(context.Background(), order.ID, order.ReceiverID)
		assert.Equal(t, domain.GetOrderStatus().ProviderSeen, order.Status)
	})

	t.Run("successful update to PICKED_UP", func(t *testing.T) {
		order, err = repo.UpdateOrderStatus(context.Background(), order.ID, domain.GetOrderStatus().PickedUp)
		assert.Empty(t, err)
		order, err = repo.GetOrder(context.Background(), order.ID, order.ReceiverID)
		assert.Equal(t, domain.GetOrderStatus().PickedUp, order.Status)
		assert.NotEmpty(t, order.PickedUpDate)
	})

	t.Run("successful update to DELIVERED", func(t *testing.T) {
		order, err = repo.UpdateOrderStatus(context.Background(), order.ID, domain.GetOrderStatus().Delivered)
		assert.Empty(t, err)
		order, err = repo.GetOrder(context.Background(), order.ID, order.ReceiverID)
		assert.Equal(t, domain.GetOrderStatus().Delivered, order.Status)
		assert.NotEmpty(t, order.DeliveryDate)
	})
}

func TestPostgres_GetOngoingOrders(t *testing.T) {
	tearUpSuite := setupSuite()
	defer tearUpSuite()

	sender, receiver, provider := setUpOrderForeignObjects(t)

	t.Run("successful get", func(t *testing.T) {
		order1, err := repo.CreateOrder(context.Background(), sender.ID, receiver.ID, provider.ID, nil)
		assert.Empty(t, err)
		order1, err = repo.UpdateOrderStatus(context.Background(), order1.ID, domain.GetOrderStatus().InProgress)
		assert.Empty(t, err)
		_, err = repo.CreateOrder(context.Background(), sender.ID, receiver.ID, provider.ID, nil)
		assert.Empty(t, err)

		orders, err := repo.GetOngoingOrders(context.Background())
		assert.Empty(t, err)
		assert.Equal(t, 1, len(orders))
		assert.Equal(t, order1.ID, orders[0].ID)
	})
}

func TestPostgres_UpdateOrderNotification(t *testing.T) {
	tearUpSuite := setupSuite()
	defer tearUpSuite()

	sender, receiver, provider := setUpOrderForeignObjects(t)

	t.Run("successful update notified_user", func(t *testing.T) {
		order, err := repo.CreateOrder(context.Background(), sender.ID, receiver.ID, provider.ID, nil)
		assert.Empty(t, err)
		assert.Equal(t, false, order.NotifiedReceiver)

		err = repo.UpdateOrderNotification(context.Background(), order.ID)
		assert.Empty(t, err)

		order, err = repo.GetOrder(context.Background(), order.ID, order.ReceiverID)
		assert.Empty(t, err)
		assert.Equal(t, true, order.NotifiedReceiver)
	})
}

func TestPostgres_GetProvidersMeanDeliveryTime(t *testing.T) {
	tearUpSuite := setupSuite()
	defer tearUpSuite()

	sender, receiver, provider := setUpOrderForeignObjects(t)
	test2 := "test-2"
	provider2, err := repo.CreateProvider(context.Background(), &test2, &test2)
	if err != nil {
		t.Error(err)
	}

	t.Run("successful get mean delivery time", func(t *testing.T) {
		order1, err := repo.CreateOrder(context.Background(), sender.ID, receiver.ID, provider.ID, nil)
		assert.Empty(t, err)
		order1, err = repo.UpdateOrderStatus(context.Background(), order1.ID, domain.GetOrderStatus().PickedUp)
		assert.Empty(t, err)
		order1, err = repo.UpdateOrderStatus(context.Background(), order1.ID, domain.GetOrderStatus().Delivered)
		assert.Empty(t, err)
		order1, err = repo.GetOrder(context.Background(), order1.ID, order1.ReceiverID)
		assert.Empty(t, err)
		order1DelTime := float32(order1.DeliveryDate.Sub(*order1.PickedUpDate).Hours() / 24)

		order2, err := repo.CreateOrder(context.Background(), sender.ID, receiver.ID, provider2.ID, nil)
		assert.Empty(t, err)
		order2, err = repo.UpdateOrderStatus(context.Background(), order2.ID, domain.GetOrderStatus().PickedUp)
		assert.Empty(t, err)
		order2, err = repo.UpdateOrderStatus(context.Background(), order2.ID, domain.GetOrderStatus().Delivered)
		assert.Empty(t, err)
		order2, err = repo.GetOrder(context.Background(), order2.ID, order2.ReceiverID)
		assert.Empty(t, err)
		order2DelTime := float32(order2.DeliveryDate.Sub(*order2.PickedUpDate).Hours() / 24)

		result, err := repo.GetProvidersMeanDeliveryTime(context.Background())
		assert.Empty(t, err)
		assert.Equal(t, 2, len(result))

		for _, res := range result {
			if res.ProviderID == order1.ProviderID {
				assert.Equal(t, order1DelTime, res.MeanDeliveryTimeInDays)
			} else {
				assert.Equal(t, order2DelTime, res.MeanDeliveryTimeInDays)
			}
		}
		if order1DelTime > order2DelTime {
			assert.Equal(t, order1DelTime, result[0].MeanDeliveryTimeInDays)
			assert.Equal(t, order2DelTime, result[1].MeanDeliveryTimeInDays)
		} else {
			assert.Equal(t, order2DelTime, result[0].MeanDeliveryTimeInDays)
			assert.Equal(t, order1DelTime, result[1].MeanDeliveryTimeInDays)
		}
	})
}
