package repository

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRepository interface {
	Create(name, email, password string) (*User, error)
	GetByEmail(email string) (*User, error)
}
