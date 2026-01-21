package handler

import (
	"log"
	"os"
	"strings"

	"github.com/NetlutZ/subscout/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return AuthHandler{authService: authService}
}

func (h AuthHandler) Register(c *fiber.Ctx) error {
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	if body.Name == "" || body.Email == "" || body.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "name, email and password required",
		})
	}

	user, err := h.authService.Register(body.Name, body.Email, body.Password)
	if err != nil {
		return c.Status(409).JSON(fiber.Map{"error": "email already exists"})
	}

	return c.Status(201).JSON(user)
}

func (h AuthHandler) Login(c *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	token, user, err := h.authService.Login(body.Email, body.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user":  user,
	})
}

func (h AuthHandler) Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "logged out successfully",
	})
}

func Protected() fiber.Handler {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error Loading .env File : ", err)
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	var jwtSecret = []byte(secret)

	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return c.Status(401).JSON(fiber.Map{"error": "missing token"})
		}

		tokenString := strings.Replace(header, "Bearer ", "", 1)

		// Verify Token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
		}

		claims := token.Claims.(jwt.MapClaims)

		c.Locals("user_id", int(claims["user_id"].(float64)))

		return c.Next()
	}
}

func RegisterAuthRoutes(app *fiber.App, authService service.AuthService) {
	h := NewAuthHandler(authService)

	auth := app.Group("/auth")
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
	auth.Post("/logout", Protected(), h.Logout)
}
