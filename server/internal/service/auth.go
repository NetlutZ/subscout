package service

import "github.com/NetlutZ/subscout/internal/repository"

type AuthService interface {
	Register(name, email, password string) (*repository.User, error)
	Login(email, password string) (string, *repository.User, error)
}
