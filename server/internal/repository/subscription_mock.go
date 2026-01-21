package repository

import "github.com/stretchr/testify/mock"

type subscriptionRepositoryMock struct {
	mock.Mock
}

func NewSubscriptionRepositoryMock() *subscriptionRepositoryMock {
	return &subscriptionRepositoryMock{}
}

func (m *subscriptionRepositoryMock) GetAll(userID int) ([]Subscription, error) {
	args := m.Called(userID)
	return args.Get(0).([]Subscription), args.Error(1)
}

func (m *subscriptionRepositoryMock) GetById(id int, userID int) (*Subscription, error) {
	args := m.Called(id, userID)
	return args.Get(0).(*Subscription), args.Error(1)
}

func (m *subscriptionRepositoryMock) Create(sub *Subscription, userID int) (*Subscription, error) {
	args := m.Called(sub, userID)
	return args.Get(0).(*Subscription), args.Error(1)
}

func (m *subscriptionRepositoryMock) Delete(id int, userID int) error {
	args := m.Called(id, userID)
	return args.Error(0)
}
