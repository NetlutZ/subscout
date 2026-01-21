package service

import "github.com/stretchr/testify/mock"

type SubscriptionServiceMock struct {
	mock.Mock
}

func NewSubscriptionServiceMock() *SubscriptionServiceMock {
	return &SubscriptionServiceMock{}
}

func (m *SubscriptionServiceMock) GetSubscriptions(userID int) ([]SubscriptionResponse, error) {
	args := m.Called(userID)
	return args.Get(0).([]SubscriptionResponse), args.Error(1)
}

func (m *SubscriptionServiceMock) GetSubscription(id int, userID int) (*SubscriptionResponse, error) {
	args := m.Called(id, userID)
	return args.Get(0).(*SubscriptionResponse), args.Error(1)
}

func (m *SubscriptionServiceMock) CreateSubscription(req CreateSubscriptionRequest, userID int) (*SubscriptionResponse, error) {
	args := m.Called(req, userID)
	return args.Get(0).(*SubscriptionResponse), args.Error(1)
}

func (m *SubscriptionServiceMock) DeleteSubscription(id int, userID int) error {
	args := m.Called(id, userID)
	return args.Error(0)
}
