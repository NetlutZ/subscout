package subscriptions

import (
	"database/sql"
	"strconv"

	"github.com/NetlutZ/subscout/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

type Subscriptions struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Category     string  `json:"category"`
	Amount       float32 `json:"amount"`
	Currency     string  `json:"currency"`
	BillingCycle string  `json:"billing_cycle"`
	BillingDate  string  `json:"billing_date"`
	Status       string  `json:"status"`
	Trial        bool    `json:"is_trial"`
}

type Service struct {
	DB *sql.DB
}

func NewService() *Service {
	return &Service{DB: db.DB}
}

func (s *Service) RegisterRoute(r fiber.Router) {
	r.Get("/subscriptions", s.getAllSubscriptions)
	r.Get("/subscriptions/:id", s.getSubscriptions)
	r.Post("/subscriptions", s.createSubscription)
	r.Delete("/subscriptions/:id", s.deleteSubscription)
}

func (s *Service) getAllSubscriptions(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	rows, err := s.DB.Query(`
		SELECT id, name, category, amount, currency, billing_cycle, billing_date, status, is_trial
		FROM subscriptions
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var subscriptions []Subscriptions
	for rows.Next() {
		var sub Subscriptions
		if err := rows.Scan(&sub.ID, &sub.Name, &sub.Category, &sub.Amount, &sub.Currency, &sub.BillingCycle, &sub.BillingDate, &sub.Status, &sub.Trial); err != nil {
			return err
		}
		subscriptions = append(subscriptions, sub)
	}
	return c.JSON(subscriptions)
}

func (s *Service) getSubscriptions(c *fiber.Ctx) error {
	subId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid subscription id",
		})
	}

	userID := c.Locals("user_id").(int)
	query := `
		SELECT id, name, category, amount, currency, billing_cycle,billing_date, status, is_trial
		FROM subscriptions
		WHERE id = $1 AND user_id = $2
	`

	sub := new(Subscriptions)
	err = s.DB.QueryRow(
		query,
		subId,
		userID,
	).Scan(&sub.ID, &sub.Name, &sub.Category, &sub.Amount, &sub.Currency, &sub.BillingCycle, &sub.BillingDate, &sub.Status, &sub.Trial)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "subscription not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(sub)
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

func handlePostgresError(c *fiber.Ctx, err error) error {
	if pgErr, ok := err.(*pq.Error); ok {
		switch pgErr.Code {
		case "23505": // unique_violation
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "subscription name already exists",
			})
		}
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": err.Error(),
	})
}

func (s *Service) createSubscription(c *fiber.Ctx) error {
	req := new(CreateSubscriptionRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.Name == "" ||
		req.Amount <= 0 ||
		req.BillingCycle == "" ||
		req.BillingDate == "" {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id, name, amount, billing_cycle, and billing_date are required",
		})
	}

	userID := c.Locals("user_id").(int)

	query := `
		INSERT INTO subscriptions (user_id, name, category, amount, currency, billing_cycle, billing_date, status, is_trial)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, name, category, amount, currency, billing_cycle, billing_date, status, is_trial
	`

	sub := new(Subscriptions)
	err := s.DB.QueryRowContext(
		c.Context(),
		query,
		userID,
		req.Name,
		req.Category,
		req.Amount,
		req.Currency,
		req.BillingCycle,
		req.BillingDate,
		req.Status,
		req.Trial,
	).Scan(&sub.ID, &sub.Name, &sub.Category, &sub.Amount, &sub.Currency, &sub.BillingCycle, &sub.BillingDate, &sub.Status, &sub.Trial)

	if err != nil {
		return handlePostgresError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(sub)
}

func (s *Service) deleteSubscription(c *fiber.Ctx) error {
	subId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	userID := c.Locals("user_id").(int)
	query := `DELETE FROM subscriptions WHERE id = $1 AND user_id = $2`

	result, err := s.DB.Exec(query, subId, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to determine delete result",
		})
	}

	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "subscription not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "subscription deleted successfully",
	})
}
