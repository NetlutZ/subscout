package service

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/NetlutZ/subscout/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return authService{userRepo: userRepo}
}

func (s authService) Register(name, email, password string) (*repository.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.Create(name, email, string(hash))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s authService) Login(email, password string) (string, *repository.User, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error Loading .env File : ", err)
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	var jwtSecret = []byte(secret)

	user, err := s.userRepo.GetByEmail(email)
	if err != nil || user == nil {
		return "", nil, errors.New("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	) != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Create Token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"name":    user.Name,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", nil, err
	}

	user.Password = "" // never expose hash
	return signed, user, nil
}
