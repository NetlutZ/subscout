package service

import (
	"github.com/NetlutZ/subscout/internal/repository"
)

type subscriptionService struct {
	subRepo repository.SubscriptionRepository
}

func NewSubscriptionService(subRepo repository.SubscriptionRepository) SubscriptionService {
	return subscriptionService{subRepo: subRepo}
}

func toResponse(sub repository.Subscription) SubscriptionResponse {
	return SubscriptionResponse{
		SubscriptionID: sub.SubscriptionID,
		Name:           sub.Name,
		Category:       sub.Category,
		Amount:         sub.Amount,
		Currency:       sub.Currency,
		BillingCycle:   sub.BillingCycle,
		BillingDate:    sub.BillingDate,
		Status:         sub.Status,
		Trial:          sub.Trial,
	}
}

func (s subscriptionService) GetSubscriptions(userID int) ([]SubscriptionResponse, error) {
	subs, err := s.subRepo.GetAll(userID)
	if err != nil {
		return nil, err
	}

	var res []SubscriptionResponse
	for _, sub := range subs {
		res = append(res, toResponse(sub))
	}

	return res, nil
}

func (s subscriptionService) GetSubscription(id int, userID int) (*SubscriptionResponse, error) {
	sub, err := s.subRepo.GetById(id, userID)
	if err != nil || sub == nil {
		return nil, err
	}

	res := toResponse(*sub)
	return &res, nil
}

func (s subscriptionService) CreateSubscription(
	req CreateSubscriptionRequest,
	userID int,
) (*SubscriptionResponse, error) {

	sub := &repository.Subscription{
		Name:         req.Name,
		Category:     req.Category,
		Amount:       req.Amount,
		Currency:     req.Currency,
		BillingCycle: req.BillingCycle,
		BillingDate:  req.BillingDate,
		Status:       req.Status,
		Trial:        req.Trial,
	}

	created, err := s.subRepo.Create(sub, userID)
	if err != nil {
		return nil, err
	}

	res := toResponse(*created)
	return &res, nil
}

func (s subscriptionService) DeleteSubscription(id int, userID int) error {
	err := s.subRepo.Delete(id, userID)
	if err != nil {
		return err
	}
	return nil
}
