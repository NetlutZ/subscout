package db

import (
	"fmt"
	"log"
)

func Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT now()
	);
	
	ALTER TABLE users
	ADD COLUMN IF NOT EXISTS name TEXT NOT NULL DEFAULT '';
	
	CREATE TABLE IF NOT EXISTS subscriptions (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,

		name VARCHAR(100) NOT NULL,       -- Netflix, Spotify
		category VARCHAR(50),             -- Entertainment, Fitness
		amount DECIMAL(10,2) NOT NULL,
		currency VARCHAR(10) DEFAULT 'THB',

		billing_cycle VARCHAR(20) NOT NULL, -- monthly, yearly
		billing_date DATE NOT NULL,         -- next renewal date

		status VARCHAR(20) DEFAULT 'active', -- active, canceled
		is_trial BOOLEAN DEFAULT false,

		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	DO $$
	BEGIN
		IF NOT EXISTS (
			SELECT 1 FROM pg_constraint WHERE conname = 'unique_user_subscription'
		) THEN
			ALTER TABLE subscriptions
			ADD CONSTRAINT unique_user_subscription
			UNIQUE (user_id, name);
		END IF;
	END$$;
		
	CREATE TABLE IF NOT EXISTS notifications (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,

		type VARCHAR(50),		-- renewal_reminder, trial_ending
		title VARCHAR(100),
		message TEXT,

		is_read BOOLEAN DEFAULT false,

		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := DB.Exec(query)
	if err != nil {
		log.Println("error while creating/migrate the database: ", err)
		return err
	}
	fmt.Println("Creation/Migration completed")
	return nil
}
