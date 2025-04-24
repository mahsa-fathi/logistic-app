package middlewares

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"logistic-app/internal/common/configs"
	"net/http"
	"strings"
)

var byteSecKey = []byte(configs.SecretKey)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := authenticate(r)
		var ctx context.Context
		if err != nil {
			ctx = context.WithValue(r.Context(), configs.AuthStatusKey, configs.AuthStatusValUnauthorized)
		} else {
			ctx = context.WithValue(r.Context(), configs.AuthStatusKey, configs.AuthStatusValAuthorized)
			ctx = context.WithValue(ctx, configs.UserIDKey, userID)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func authenticate(r *http.Request) (uint, error) {
	header := getHeader(r)
	if header == "" {
		return 0, fmt.Errorf("auth header name not found in headers")
	}

	rawToken, err := getRawToken(header)
	if rawToken == "" {
		return 0, err
	}

	valToken, err := getValidatedToken(rawToken)
	if valToken == nil {
		return 0, err
	}
	return getUser(valToken)
}

func getHeader(r *http.Request) string {
	header := r.Header.Get(configs.JWTDefaults["AUTH_HEADER_NAME"].(string))
	return header
}

func getRawToken(h string) (string, error) {
	arg := strings.Split(h, " ")
	if arg[0] != configs.JWTDefaults["AUTH_HEADER_TYPES"].([]string)[0] {
		return "", fmt.Errorf("auth header type not found in header")
	}

	if len(arg) != 2 {
		return "", fmt.Errorf("auth header does not have the right length")
	}

	return arg[1], nil
}

func getValidatedToken(r string) (map[string]any, error) {
	token, err := jwt.Parse(r, func(t *jwt.Token) (interface{}, error) {
		return byteSecKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, fmt.Errorf("error in parsing token claims")
}

func getUser(t map[string]any) (uint, error) {
	userID, ok := t[configs.JWTDefaults["USER_ID_CLAIM"].(string)].(float64)
	if !ok {
		return 0, fmt.Errorf("could not map to float")
	}
	return uint(userID), nil
}
