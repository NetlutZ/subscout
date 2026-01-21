package repository

import "github.com/stretchr/testify/mock"

type userRepositoryMock struct {
	mock.Mock
}

func NewUserRepositoryMock() *userRepositoryMock {
	return &userRepositoryMock{}
}

func (m *userRepositoryMock) Create(name, email, password string) (*User, error) {
	args := m.Called(name, email, password)
	return args.Get(0).(*User), args.Error(1)
}

func (m *userRepositoryMock) GetByEmail(email string) (*User, error) {
	args := m.Called(email)
	return args.Get(0).(*User), args.Error(1)
}
