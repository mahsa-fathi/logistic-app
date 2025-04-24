package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestPostgres_CreateCustomer(t *testing.T) {
	tearUpSuite := setupSuite()
	defer tearUpSuite()

	name := "name"
	address := "somewhere"
	postal := "some-code"

	t.Run("successful create", func(t *testing.T) {
		phone := "09"
		customer, err := repo.CreateCustomer(context.Background(), &name, &phone, &address, &postal)
		assert.Empty(t, err)
		assert.Equal(t, name, *customer.Name)
		assert.Equal(t, phone, customer.PhoneNumber)
		assert.Equal(t, address, customer.Address)
		assert.Equal(t, postal, customer.PostalCode)
	})

	t.Run("successful create with no name", func(t *testing.T) {
		phone := "08"
		customer, err := repo.CreateCustomer(context.Background(), nil, &phone, &address, &postal)
		assert.Empty(t, err)
		assert.Empty(t, customer.Name)
		assert.Equal(t, phone, customer.PhoneNumber)
	})
}

func TestPostgres_GetCustomer(t *testing.T) {
	tearUpSuite := setupSuite()
	defer tearUpSuite()

	name := "name"
	address := "somewhere"
	postal := "some-code"

	t.Run("successful get", func(t *testing.T) {
		phone := "09"
		user, err := repo.CreateCustomer(context.Background(), &name, &phone, &address, &postal)
		assert.Empty(t, err)

		customer, err := repo.GetCustomer(context.Background(), user.ID)
		assert.Empty(t, err)
		assert.Equal(t, user.ID, customer.ID)
		assert.Equal(t, phone, customer.PhoneNumber)
	})

	t.Run("unsuccessful get", func(t *testing.T) {
		_, err := repo.GetCustomer(context.Background(), 1000)
		assert.NotEmpty(t, err)
		assert.Equal(t, http.StatusNotFound, err.Code)
	})
}
