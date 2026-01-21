package repository

import (
	"database/sql"
)

type subscriptionRepositoryDB struct {
	db *sql.DB
}

func NewSubscriptionRepositoryDB(db *sql.DB) SubscriptionRepository {
	return subscriptionRepositoryDB{db: db}
}

func (r subscriptionRepositoryDB) GetAll(userID int) ([]Subscription, error) {
	query := `
		SELECT id, name, category, amount, currency,
		       billing_cycle, billing_date, status, is_trial
		FROM subscriptions
		WHERE user_id = $1
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []Subscription
	for rows.Next() {
		var sub Subscription
		if err := rows.Scan(
			&sub.SubscriptionID,
			&sub.Name,
			&sub.Category,
			&sub.Amount,
			&sub.Currency,
			&sub.BillingCycle,
			&sub.BillingDate,
			&sub.Status,
			&sub.Trial,
		); err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}

	return subs, nil
}

func (r subscriptionRepositoryDB) GetById(id int, userID int) (*Subscription, error) {
	query := `
		SELECT id, name, category, amount, currency,
		       billing_cycle, billing_date, status, is_trial
		FROM subscriptions
		WHERE id = $1 AND user_id = $2
	`

	var sub Subscription
	err := r.db.QueryRow(query, id, userID).Scan(
		&sub.SubscriptionID,
		&sub.Name,
		&sub.Category,
		&sub.Amount,
		&sub.Currency,
		&sub.BillingCycle,
		&sub.BillingDate,
		&sub.Status,
		&sub.Trial,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (r subscriptionRepositoryDB) Create(sub *Subscription, userID int) (*Subscription, error) {
	query := `
		INSERT INTO subscriptions
		(name, category, amount, currency, billing_cycle, billing_date, status, is_trial, user_id)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		sub.Name,
		sub.Category,
		sub.Amount,
		sub.Currency,
		sub.BillingCycle,
		sub.BillingDate,
		sub.Status,
		sub.Trial,
		userID,
	).Scan(&sub.SubscriptionID)

	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (r subscriptionRepositoryDB) Delete(id int, userID int) error {
	query := `
		DELETE FROM subscriptions
		WHERE id = $1 AND user_id = $2
	`

	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
