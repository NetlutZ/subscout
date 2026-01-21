package service

type SubscriptionResponse struct {
	SubscriptionID int     `json:"id"`
	Name           string  `json:"name"`
	Category       string  `json:"category"`
	Amount         float32 `json:"amount"`
	Currency       string  `json:"currency"`
	BillingCycle   string  `json:"billing_cycle"`
	BillingDate    string  `json:"billing_date"`
	Status         string  `json:"status"`
	Trial          bool    `json:"is_trial"`
}

type CreateSubscriptionRequest struct {
	Name         string  `json:"name"`
	Category     string  `json:"category"`
	Amount       float32 `json:"amount"`
	Currency     string  `json:"currency"`
	BillingCycle string  `json:"billing_cycle"`
	BillingDate  string  `json:"billing_date"`
	Status       string  `json:"status"`
	Trial        bool    `json:"is_trial"`
}

type SubscriptionService interface {
	GetSubscriptions(userID int) ([]SubscriptionResponse, error)
	GetSubscription(id int, userID int) (*SubscriptionResponse, error)
	CreateSubscription(req CreateSubscriptionRequest, userID int) (*SubscriptionResponse, error)
	DeleteSubscription(id int, userID int) error
}
