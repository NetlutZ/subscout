package repository

import "database/sql"

type userRepositoryDB struct {
	db *sql.DB
}

func NewUserRepositoryDB(db *sql.DB) UserRepository {
	return userRepositoryDB{db: db}
}

func (r userRepositoryDB) Create(name, email, password string) (*User, error) {
	var user User

	err := r.db.QueryRow(`
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, name, email
	`, name, email, password).
		Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r userRepositoryDB) GetByEmail(email string) (*User, error) {
	var user User

	err := r.db.QueryRow(`
		SELECT id, name, email, password
		FROM users
		WHERE email = $1
	`, email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
