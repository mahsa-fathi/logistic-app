package middlewares

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"logistic-app/internal/common/configs"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthenticate(t *testing.T) {
	userID := uint(6)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(5 * time.Minute).Unix(),
		configs.JWTDefaults["USER_ID_CLAIM"].(string): userID,
	})
	tokenStr, err := token.SignedString(byteSecKey)
	if err != nil {
		t.Error(err)
	}

	t.Run("Successful Authentication", func(t *testing.T) {
		request := httptest.NewRequest("GET", "http://localhost:8080/test/", http.NoBody)
		request.Header.Set(configs.JWTDefaults["AUTH_HEADER_NAME"].(string), "Bearer "+tokenStr)

		user, _ := authenticate(request)
		assert.Equal(t, userID, user)
	})

	t.Run("Without Authentication", func(t *testing.T) {
		request := httptest.NewRequest("GET", "http://localhost:8080/test/", http.NoBody)
		user, _ := authenticate(request)
		assert.Empty(t, user)
	})

	t.Run("Without User Id", func(t *testing.T) {
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"exp": time.Now().Add(5 * time.Minute).Unix(),
		})
		tokenStr, err = token.SignedString(byteSecKey)
		if err != nil {
			t.Error(err)
		}

		request := httptest.NewRequest("GET", "http://localhost:8080/test/", http.NoBody)
		request.Header.Set(configs.JWTDefaults["AUTH_HEADER_NAME"].(string), "Bearer "+tokenStr)

		user, _ := authenticate(request)
		assert.Empty(t, user)
	})
}
