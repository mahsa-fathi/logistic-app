package http

import (
	"context"
	"log"
	"logistic-app/internal/adapters/http/middlewares"
	"logistic-app/internal/adapters/http/models"
	"logistic-app/internal/app/domain"
	"logistic-app/internal/app/ports"
	"logistic-app/internal/common/configs"
	"logistic-app/internal/common/errors"
	"net/http"
	"reflect"
)

type responseFunc func(request *http.Request) *models.Response

type Server struct {
	listenAddr string
	service    ports.Service
}

func NewServer(service ports.Service) *Server {
	return &Server{
		listenAddr: configs.ServerURL,
		service:    service,
	}
}

func (s *Server) Run() {
	router := http.NewServeMux()
	stack := middlewares.MiddlewareStack(
		middlewares.Logging,
		middlewares.JWTMiddleware,
	)

	router.HandleFunc("GET /api/health/", makeHTTPHandleFunc(perform(s.service.HealthCheck)))

	router.HandleFunc("GET /api/providers/", makeHTTPHandleFunc(perform(s.service.GetProviders)))
	router.HandleFunc("GET /api/providers/report/", makeHTTPHandleFunc(perform(s.service.GetProvidersMeanDelTime)))
	router.HandleFunc("POST /api/provider/", makeHTTPHandleFunc(performWith(s.service.CreateProvider)))

	router.HandleFunc("POST /api/customer/", makeHTTPHandleFunc(performWith(s.service.CreateCustomer)))
	router.HandleFunc("POST /api/customer/token/", makeHTTPHandleFunc(performWith(s.service.GetCustomerToken)))

	router.HandleFunc("POST /api/order/", makeHTTPHandleFuncWithAuth(performWith(s.service.CreateOrder)))
	router.HandleFunc("GET /api/order/{order_id}/", makeHTTPHandleFuncWithAuth(performWith(s.service.GetOrder)))

	server := http.Server{
		Addr:    s.listenAddr,
		Handler: stack(router),
	}

	log.Println("API Server Running on:", s.listenAddr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func performWith[T any, S domain.Request](f func(ctx context.Context, body S) (T, *errors.AppError)) responseFunc {
	return func(request *http.Request) *models.Response {
		var b S
		bodyType := reflect.TypeOf(b)
		body := reflect.New(bodyType.Elem()).Interface().(S)

		var err *errors.AppError
		if request.Method == "GET" {
			err = body.UnmarshalPathValue(request)
		} else {
			err = body.UnmarshalBody(request)
		}
		if err != nil {
			return models.ReturnErrorResp(request.Context(), err)
		}
		result, err := f(request.Context(), body)
		resp := models.ReturnResp(request.Context(), result)
		if err != nil {
			resp = models.ReturnErrorResp(request.Context(), err)
		}
		return resp
	}
}

func perform[T any](f func(ctx context.Context) (T, *errors.AppError)) responseFunc {
	return func(request *http.Request) *models.Response {
		result, err := f(request.Context())
		resp := models.ReturnResp(request.Context(), result)
		if err != nil {
			resp = models.ReturnErrorResp(request.Context(), err)
		}
		return resp
	}
}

func makeHTTPHandleFuncWithAuth(f responseFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Context().Value(configs.AuthStatusKey) == configs.AuthStatusValUnauthorized {
			resp := models.ReturnErrorResp(request.Context(), errors.Unauthorized())
			models.WriteJSON(writer, resp)
		} else {
			request = request.WithContext(request.Context())
			resp := f(request)
			models.WriteJSON(writer, resp)
		}
	}
}

func makeHTTPHandleFunc(f responseFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		request = request.WithContext(request.Context())
		resp := f(request)
		models.WriteJSON(writer, resp)
	}
}
