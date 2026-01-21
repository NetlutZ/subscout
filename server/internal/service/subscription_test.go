package service_test

import (
	"errors"
	"testing"

	"github.com/NetlutZ/subscout/internal/repository"
	"github.com/NetlutZ/subscout/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetSubscriptions(t *testing.T) {

	t.Run("Get Subscriptions Success", func(t *testing.T) {
		// arrange
		subscriptionRepo := repository.NewSubscriptionRepositoryMock()
		subscriptionRepo.
			On("GetAll", 1).
			Return([]repository.Subscription{
				{
					SubscriptionID: 1,
					Name:           "Netflix",
					Category:       "Entertain",
					Amount:         20,
					Currency:       "THB",
					BillingCycle:   "Monthly",
					BillingDate:    "30/01/2568",
					Status:         "active",
					Trial:          false,
				},
			}, nil)

		subscriptionService := service.NewSubscriptionService(subscriptionRepo)

		// act
		subs, err := subscriptionService.GetSubscriptions(1)

		// assert
		assert.NoError(t, err)

		expected := []service.SubscriptionResponse{
			{
				SubscriptionID: 1,
				Name:           "Netflix",
				Category:       "Entertain",
				Amount:         20,
				Currency:       "THB",
				BillingCycle:   "Monthly",
				BillingDate:    "30/01/2568",
				Status:         "active",
				Trial:          false,
			},
		}

		assert.Equal(t, expected, subs)
		subscriptionRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		// arrange
		subscriptionRepo := repository.NewSubscriptionRepositoryMock()
		expectedErr := errors.New("database error")

		subscriptionRepo.
			On("GetAll", 1).
			Return([]repository.Subscription(nil), expectedErr)

		subscriptionService := service.NewSubscriptionService(subscriptionRepo)

		// act
		subs, err := subscriptionService.GetSubscriptions(1)

		// assert
		assert.Nil(t, subs)
		assert.EqualError(t, err, expectedErr.Error())
		subscriptionRepo.AssertExpectations(t)
	})
}

func TestGetSubscription(t *testing.T) {
	t.Run("Get Subscription Success", func(t *testing.T) {
		// arrange
		subscriptionRepo := repository.NewSubscriptionRepositoryMock()

		subscriptionRepo.
			On("GetById", 1, 10).
			Return(&repository.Subscription{
				SubscriptionID: 1,
				Name:           "Netflix",
				Category:       "Entertain",
				Amount:         20,
				Currency:       "THB",
				BillingCycle:   "Monthly",
				BillingDate:    "30/01/2568",
				Status:         "active",
				Trial:          false,
			}, nil)

		subService := service.NewSubscriptionService(subscriptionRepo)

		// act
		res, err := subService.GetSubscription(1, 10)

		// assert
		assert.NoError(t, err)
		assert.NotNil(t, res)

		expected := &service.SubscriptionResponse{
			SubscriptionID: 1,
			Name:           "Netflix",
			Category:       "Entertain",
			Amount:         20,
			Currency:       "THB",
			BillingCycle:   "Monthly",
			BillingDate:    "30/01/2568",
			Status:         "active",
			Trial:          false,
		}

		assert.Equal(t, expected, res)
		subscriptionRepo.AssertExpectations(t)
	})

	t.Run("Subscription Not Found", func(t *testing.T) {
		// arrange
		subscriptionRepo := repository.NewSubscriptionRepositoryMock()
		expectedErr := errors.New("not found")

		subscriptionRepo.
			On("GetById", 1, 10).
			Return((*repository.Subscription)(nil), expectedErr)

		subService := service.NewSubscriptionService(subscriptionRepo)

		// act
		res, err := subService.GetSubscription(1, 10)

		// assert
		assert.Nil(t, res)
		assert.EqualError(t, err, "not found")
		subscriptionRepo.AssertExpectations(t)
	})

}

func TestCreateSubscription(t *testing.T) {
	t.Run("Create Subscription Success", func(t *testing.T) {
		// arrange
		subscriptionRepo := repository.NewSubscriptionRepositoryMock()

		req := service.CreateSubscriptionRequest{
			Name:         "Netflix",
			Category:     "Entertain",
			Amount:       20,
			Currency:     "THB",
			BillingCycle: "Monthly",
			BillingDate:  "30/01/2568",
			Status:       "active",
			Trial:        false,
		}

		subscriptionRepo.
			On("Create", mock.AnythingOfType("*repository.Subscription"), 10).
			Return(&repository.Subscription{
				SubscriptionID: 1,
				Name:           "Netflix",
				Category:       "Entertain",
				Amount:         20,
				Currency:       "THB",
				BillingCycle:   "Monthly",
				BillingDate:    "30/01/2568",
				Status:         "active",
				Trial:          false,
			}, nil)

		subService := service.NewSubscriptionService(subscriptionRepo)

		// act
		res, err := subService.CreateSubscription(req, 10)

		// assert
		assert.NoError(t, err)
		assert.NotNil(t, res)

		assert.Equal(t, "Netflix", res.Name)
		assert.Equal(t, float32(20), res.Amount)
		assert.Equal(t, "THB", res.Currency)

		subscriptionRepo.AssertExpectations(t)
	})

	t.Run("Create Subscription Error", func(t *testing.T) {
		// arrange
		subscriptionRepo := repository.NewSubscriptionRepositoryMock()
		expectedErr := errors.New("insert failed")

		req := service.CreateSubscriptionRequest{
			Name: "Netflix",
		}

		subscriptionRepo.
			On("Create", mock.Anything, 10).
			Return((*repository.Subscription)(nil), expectedErr)

		subService := service.NewSubscriptionService(subscriptionRepo)

		// act
		res, err := subService.CreateSubscription(req, 10)

		// assert
		assert.Nil(t, res)
		assert.EqualError(t, err, "insert failed")
		subscriptionRepo.AssertExpectations(t)
	})

}

func TestDeleteSubscription(t *testing.T) {
	t.Run("Delete Subscription Success", func(t *testing.T) {
		// arrange
		subscriptionRepo := repository.NewSubscriptionRepositoryMock()

		subscriptionRepo.
			On("Delete", 1, 10).
			Return(nil)

		subService := service.NewSubscriptionService(subscriptionRepo)

		// act
		err := subService.DeleteSubscription(1, 10)

		// assert
		assert.NoError(t, err)
		subscriptionRepo.AssertExpectations(t)
	})

	t.Run("Delete Subscription Error", func(t *testing.T) {
		// arrange
		subscriptionRepo := repository.NewSubscriptionRepositoryMock()
		expectedErr := errors.New("delete failed")

		subscriptionRepo.
			On("Delete", 1, 10).
			Return(expectedErr)

		subService := service.NewSubscriptionService(subscriptionRepo)

		// act
		err := subService.DeleteSubscription(1, 10)

		// assert
		assert.EqualError(t, err, "delete failed")
		subscriptionRepo.AssertExpectations(t)
	})

}
