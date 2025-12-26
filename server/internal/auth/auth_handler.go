package auth

import (
	"database/sql"
	"time"

	"github.com/NetlutZ/subscout/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("SUPER_SECRET_KEY") // move to env later

type Service struct {
	DB *sql.DB
}

func NewService() *Service {
	return &Service{DB: db.DB}
}

func (s *Service) RegisterRoute(r fiber.Router) {
	r.Post("/register", s.register)
	r.Post("/login", s.login)
	r.Post("/logout", Protected(), s.logout)
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (s *Service) register(c *fiber.Ctx) error {
	type request struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	if body.Name == "" || body.Email == "" || body.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name, email and password required"})
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	var user User
	err := s.DB.QueryRow(`
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, name, email
	`, body.Name, body.Email, string(hash)).
		Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		return c.Status(409).JSON(fiber.Map{"error": "email already exists"})
	}

	return c.Status(201).JSON(user)
}

func (s *Service) login(c *fiber.Ctx) error {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	var user User
	err := s.DB.QueryRow(`
		SELECT id, name, email, password
		FROM users
		WHERE email=$1
	`, body.Email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)) != nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"name":    user.Name,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString(jwtSecret)

	return c.JSON(fiber.Map{
		"token": signedToken,
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func (s *Service) logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "logged out successfully",
	})
}
