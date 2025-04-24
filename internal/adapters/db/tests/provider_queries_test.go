package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log"
	"logistic-app/internal/adapters/db"
	"logistic-app/internal/app/ports"
	"net/http"
	"testing"
)

var repo ports.Repo
var e error

func setupSuite() func() {
	repo, e = db.NewMockPostgresDB()
	if e != nil {
		log.Fatal(e)
	}
	return func() {
		repo.Close()
	}
}

func TestPostgres_CreateProvider(t *testing.T) {
	tearUpSuite := setupSuite()
	defer tearUpSuite()

	t.Run("successful create", func(t *testing.T) {
		name := "test"
		url := "test"

		actProv, err := repo.CreateProvider(context.Background(), &name, &url)
		assert.Empty(t, err)
		assert.Equal(t, name, actProv.Name)
		assert.Equal(t, url, actProv.Url)
	})

	t.Run("duplicate create", func(t *testing.T) {
		name := "test-2"
		url := "test-2"

		_, err := repo.CreateProvider(context.Background(), &name, &url)
		assert.Empty(t, err)

		_, err = repo.CreateProvider(context.Background(), &name, &url)
		assert.NotEmpty(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
	})
}

func TestPostgres_GetProvider(t *testing.T) {
	tearUpSuite := setupSuite()
	defer tearUpSuite()

	t.Run("successful get", func(t *testing.T) {
		test := "test"
		provider, err := repo.CreateProvider(context.Background(), &test, &test)
		assert.Empty(t, err)

		actProv, err := repo.GetProvider(context.Background(), provider.ID)
		assert.Empty(t, err)
		assert.Equal(t, provider.ID, actProv.ID)
		assert.Equal(t, provider.Name, actProv.Name)
		assert.Equal(t, provider.Url, actProv.Url)
	})

	t.Run("unsuccessful get", func(t *testing.T) {
		_, err := repo.GetProvider(context.Background(), 5)
		assert.Equal(t, http.StatusNotFound, err.Code)
	})
}

func TestPostgres_GetAllProviders(t *testing.T) {
	tearUpSuite := setupSuite()
	defer tearUpSuite()

	t.Run("successful get", func(t *testing.T) {
		test1 := "test-1"
		test2 := "test-2"
		_, err := repo.CreateProvider(context.Background(), &test1, &test1)
		assert.Empty(t, err)
		_, err = repo.CreateProvider(context.Background(), &test2, &test2)
		assert.Empty(t, err)

		providers, err := repo.GetAllProviders(context.Background())
		assert.Empty(t, err)
		assert.Equal(t, 2, len(providers))
		for _, prov := range providers {
			assert.True(t, prov.Name == test1 || prov.Name == test2)
		}
	})
}
