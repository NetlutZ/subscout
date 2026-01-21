package service_test

import (
	"errors"
	"testing"

	"github.com/NetlutZ/subscout/internal/repository"
	"github.com/NetlutZ/subscout/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	t.Run("Register Success", func(t *testing.T) {
		// arrange
		userRepo := repository.NewUserRepositoryMock()

		userRepo.
			On("Create", "John", "john@test.com", mock.AnythingOfType("string")).
			Return(&repository.User{
				ID:    1,
				Name:  "John",
				Email: "john@test.com",
			}, nil)

		svc := service.NewAuthService(userRepo)

		// act
		user, err := svc.Register("John", "john@test.com", "password123")

		// assert
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "John", user.Name)
		assert.Equal(t, "john@test.com", user.Email)

		userRepo.AssertExpectations(t)
	})

	t.Run("Register Error", func(t *testing.T) {
		// arrange
		userRepo := repository.NewUserRepositoryMock()
		expectedErr := errors.New("insert failed")

		userRepo.
			On("Create", mock.Anything, mock.Anything, mock.Anything).
			Return((*repository.User)(nil), expectedErr)

		svc := service.NewAuthService(userRepo)

		// act
		user, err := svc.Register("John", "john@test.com", "password123")

		// assert
		assert.Nil(t, user)
		assert.EqualError(t, err, "insert failed")

		userRepo.AssertExpectations(t)
	})

}

func TestLogin(t *testing.T) {
	t.Run("Login Success", func(t *testing.T) {
		// arrange
		userRepo := repository.NewUserRepositoryMock()

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10)

		userRepo.
			On("GetByEmail", "john@test.com").
			Return(&repository.User{
				ID:       1,
				Name:     "John",
				Email:    "john@test.com",
				Password: string(hashedPassword),
			}, nil)

		svc := service.NewAuthService(userRepo)

		// act
		token, user, err := svc.Login("john@test.com", "password123")

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotNil(t, user)
		assert.Equal(t, "John", user.Name)
		assert.Empty(t, user.Password) // password must be cleared

		userRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		// arrange
		userRepo := repository.NewUserRepositoryMock()

		userRepo.
			On("GetByEmail", "john@test.com").
			Return((*repository.User)(nil), nil)

		svc := service.NewAuthService(userRepo)

		// act
		token, user, err := svc.Login("john@test.com", "password123")

		// assert
		assert.Empty(t, token)
		assert.Nil(t, user)
		assert.EqualError(t, err, "invalid credentials")

		userRepo.AssertExpectations(t)
	})

	t.Run("Invalid Password", func(t *testing.T) {
		// arrange
		userRepo := repository.NewUserRepositoryMock()

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), 10)

		userRepo.
			On("GetByEmail", "john@test.com").
			Return(&repository.User{
				ID:       1,
				Name:     "John",
				Email:    "john@test.com",
				Password: string(hashedPassword),
			}, nil)

		svc := service.NewAuthService(userRepo)

		// act
		token, user, err := svc.Login("john@test.com", "wrongpassword")

		// assert
		assert.Empty(t, token)
		assert.Nil(t, user)
		assert.EqualError(t, err, "invalid credentials")

		userRepo.AssertExpectations(t)
	})

	t.Run("Repo Error", func(t *testing.T) {
		// arrange
		userRepo := repository.NewUserRepositoryMock()
		expectedErr := errors.New("db error")

		userRepo.
			On("GetByEmail", "john@test.com").
			Return((*repository.User)(nil), expectedErr)

		svc := service.NewAuthService(userRepo)

		// act
		token, user, err := svc.Login("john@test.com", "password123")

		// assert
		assert.Empty(t, token)
		assert.Nil(t, user)
		assert.EqualError(t, err, "invalid credentials")

		userRepo.AssertExpectations(t)
	})

}
