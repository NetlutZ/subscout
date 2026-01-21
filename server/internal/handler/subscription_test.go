package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NetlutZ/subscout/internal/handler"
	"github.com/NetlutZ/subscout/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func mockAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("user_id", 10)
		return c.Next()
	}
}

func setupApp(mockSvc *service.SubscriptionServiceMock) *fiber.App {
	app := fiber.New()

	h := handler.NewSubscriptionHandler(mockSvc)

	api := app.Group("/api")
	subscriptions := api.Group("/subscriptions", mockAuth())

	subscriptions.Get("/", h.GetSubscriptions)
	subscriptions.Get("/:id", h.GetSubscription)
	subscriptions.Post("/", h.CreateSubscription)
	subscriptions.Delete("/:id", h.DeleteSubscription)

	return app
}

func TestGetSubscriptions(t *testing.T) {
	tests := []struct {
		name       string
		mockReturn []service.SubscriptionResponse
		mockErr    error
		status     int
	}{
		{
			name: "success",
			mockReturn: []service.SubscriptionResponse{
				{SubscriptionID: 1, Name: "Netflix"},
			},
			status: fiber.StatusOK,
		},
		{
			name:    "service error",
			mockErr: errors.New("db error"),
			status:  fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewSubscriptionServiceMock()

			svc.On("GetSubscriptions", 10).
				Return(tt.mockReturn, tt.mockErr)

			app := setupApp(svc)

			req := httptest.NewRequest(http.MethodGet, "/api/subscriptions", nil)
			resp, _ := app.Test(req)

			assert.Equal(t, tt.status, resp.StatusCode)
			svc.AssertExpectations(t)
		})
	}
}

func TestGetSubscription(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		mockReturn *service.SubscriptionResponse
		mockErr    error
		status     int
	}{
		{
			name: "success",
			id:   "1",
			mockReturn: &service.SubscriptionResponse{
				SubscriptionID: 1,
				Name:           "Netflix",
			},
			status: fiber.StatusOK,
		},
		{
			name:   "invalid id",
			id:     "abc",
			status: fiber.StatusBadRequest,
		},
		{
			name:    "service error",
			id:      "1",
			mockErr: errors.New("db error"),
			status:  fiber.StatusInternalServerError,
		},
		{
			name:   "not found",
			id:     "1",
			status: fiber.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewSubscriptionServiceMock()

			if tt.id == "1" {
				svc.On("GetSubscription", 1, 10).
					Return(tt.mockReturn, tt.mockErr)
			}

			app := setupApp(svc)

			req := httptest.NewRequest(http.MethodGet, "/api/subscriptions/"+tt.id, nil)
			resp, _ := app.Test(req)

			assert.Equal(t, tt.status, resp.StatusCode)
			svc.AssertExpectations(t)
		})
	}
}

func TestCreateSubscription(t *testing.T) {
	body := service.CreateSubscriptionRequest{
		Name: "Netflix",
	}

	jsonBody, _ := json.Marshal(body)

	tests := []struct {
		name    string
		body    []byte
		mockErr error
		status  int
	}{
		{
			name:   "success",
			status: fiber.StatusCreated,
		},
		{
			name:    "service error",
			mockErr: errors.New("insert failed"),
			status:  fiber.StatusInternalServerError,
		},
		{
			name:   "invalid body",
			body:   []byte("{invalid"),
			status: fiber.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewSubscriptionServiceMock()

			if tt.body == nil || string(tt.body) == string(jsonBody) {
				svc.On("CreateSubscription", mock.Anything, 10).
					Return(&service.SubscriptionResponse{SubscriptionID: 1}, tt.mockErr)
			}

			app := setupApp(svc)

			reqBody := tt.body
			if reqBody == nil {
				reqBody = jsonBody
			}

			req := httptest.NewRequest(http.MethodPost, "/api/subscriptions", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			assert.Equal(t, tt.status, resp.StatusCode)
			svc.AssertExpectations(t)
		})
	}
}

func TestDeleteSubscription(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mockErr error
		status  int
	}{
		{
			name:   "success",
			id:     "1",
			status: fiber.StatusOK,
		},
		{
			name:    "service error",
			id:      "1",
			mockErr: errors.New("delete failed"),
			status:  fiber.StatusInternalServerError,
		},
		{
			name:   "invalid id",
			id:     "abc",
			status: fiber.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewSubscriptionServiceMock()

			if tt.id == "1" {
				svc.On("DeleteSubscription", 1, 10).
					Return(tt.mockErr)
			}

			app := setupApp(svc)

			req := httptest.NewRequest(http.MethodDelete, "/api/subscriptions/"+tt.id, nil)
			resp, _ := app.Test(req)

			assert.Equal(t, tt.status, resp.StatusCode)
			svc.AssertExpectations(t)
		})
	}
}
