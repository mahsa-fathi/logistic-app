package http

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"logistic-app/internal/app/service"
	"logistic-app/internal/common/configs"
	"net/http"
	"net/http/httptest"
	"testing"
)

var logisticService *service.LogisticService
var server *Server

func setupSuite() {
	logisticService = service.NewLogisticService(nil)
	server = NewServer(logisticService)
}

func setupTest() (func(), *http.Request, *httptest.ResponseRecorder) {
	leagueId, _ := uuid.NewRandom()
	userId, _ := uuid.NewRandom()
	limit := 1
	offset := 0

	ctx := context.Background()
	ctx = context.WithValue(ctx, configs.UserIDKey, &userId)

	targetUrl := fmt.Sprintf("/league/%s/leaderboard/?limit=%v&offset=%v", leagueId.String(), limit, offset)
	req := httptest.NewRequest(http.MethodGet, targetUrl, nil).WithContext(ctx)
	req.SetPathValue("league_id", leagueId.String())

	w := httptest.NewRecorder()

	return func() {
		logisticService = service.NewLogisticService(nil)
		server = NewServer(logisticService)
	}, req, w
}

func TestHome(t *testing.T) {
	setupSuite()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	ctx := req.Context()

	makeHTTPHandleFunc(perform(server.service.HealthCheck))(w, req.WithContext(ctx))

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	assert.Equal(t, nil, err, "response body cannot be read")
	assert.Equal(t, nil, err, "response body cannot be unmarshalled")
	assert.Equal(t, "{\"authorization\":null,\"communication\":{\"kaka\":true,\"kompany\":true,\"nesta\":true,\"pirlo\":true},\"status\":\"healthy\"}\n", string(body))
}

func TestLeaderboard(t *testing.T) {
	setupSuite()
	//id, _ := uuid.NewRandom()
	//weekId, _ := uuid.NewRandom()
	//userDN := "my_display_name"

	t.Run("successful_test", func(t *testing.T) {

	})

	t.Cleanup(func() {
		mock.AssertExpectationsForObjects(t)
	})
}
