package handler

import (
	"strconv"

	"github.com/NetlutZ/subscout/internal/service"
	"github.com/gofiber/fiber/v2"
)

type subscriptionHandler struct {
	subService service.SubscriptionService
}

func NewSubscriptionHandler(subService service.SubscriptionService) subscriptionHandler {
	return subscriptionHandler{subService: subService}
}

func RegisterSubscriptionRoutes(app *fiber.App, subService service.SubscriptionService) {
	h := NewSubscriptionHandler(subService)

	api := app.Group("/api")
	subscriptions := api.Group("/subscriptions", Protected())

	subscriptions.Get("/", h.GetSubscriptions)
	subscriptions.Get("/:id", h.GetSubscription)
	subscriptions.Post("/", h.CreateSubscription)
	subscriptions.Delete("/:id", h.DeleteSubscription)
}

func getUserID(c *fiber.Ctx) (int, error) {
	userID, ok := c.Locals("user_id").(int)
	if !ok || userID <= 0 {
		return 0, fiber.ErrUnauthorized
	}
	return userID, nil
}

// GET /subscriptions
func (h subscriptionHandler) GetSubscriptions(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	subs, err := h.subService.GetSubscriptions(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(subs)
}

// GET /subscriptions/:id
func (h subscriptionHandler) GetSubscription(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid subscription id",
		})
	}

	sub, err := h.subService.GetSubscription(id, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if sub == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "subscription not found",
		})
	}

	return c.JSON(sub)
}

// POST /subscriptions
func (h subscriptionHandler) CreateSubscription(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	var req service.CreateSubscriptionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	sub, err := h.subService.CreateSubscription(req, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(sub)
}

// DELETE /subscriptions/:id
func (h subscriptionHandler) DeleteSubscription(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid subscription id",
		})
	}

	err = h.subService.DeleteSubscription(id, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON("message : delete success")
}
