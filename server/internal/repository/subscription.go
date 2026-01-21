package repository

type Subscription struct {
	SubscriptionID int     `db:"id"`
	Name           string  `db:"name"`
	Category       string  `db:"category"`
	Amount         float32 `db:"amount"`
	Currency       string  `db:"currency"`
	BillingCycle   string  `db:"billing_cycle"`
	BillingDate    string  `db:"billing_date"`
	Status         string  `db:"status"`
	Trial          bool    `db:"is_trial"`
}

type SubscriptionRepository interface {
	GetAll(userID int) ([]Subscription, error)
	GetById(id int, userID int) (*Subscription, error)
	Create(sub *Subscription, userID int) (*Subscription, error)
	Delete(id int, userID int) error
}
