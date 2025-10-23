package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"mini-ecommerce/backend/internal/models"
	"mini-ecommerce/backend/internal/repositories"
)

type WebhookController struct {
	Orders       *repositories.OrderRepository
	Transactions *repositories.TransactionRepository
}

type mpesaWebhook struct {
	TransactionID string `json:"transaction_id"`
	OrderID       uint   `json:"order_id,string"`
	Status        string `json:"status"`
	AmountCents   int64  `json:"amount_cents"`
}

func NewWebhookController(o *repositories.OrderRepository, t *repositories.TransactionRepository) *WebhookController {
	return &WebhookController{Orders: o, Transactions: t}
}

func (h *WebhookController) Mpesa(c *fiber.Ctx) error {
	b, _ := io.ReadAll(c.Body())
	var w mpesaWebhook
	if err := json.Unmarshal(b, &w); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "bad payload"})
	}
	// Update transaction status by provider ref
	status := models.TxFailed
	if w.Status == "SUCCESS" { status = models.TxSuccess }
	if err := h.Transactions.UpdateStatusByProviderRef(w.TransactionID, status, string(b)); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "tx update failed"})
	}
	// Update order state
	orderStatus := models.OrderFailed
	if status == models.TxSuccess { orderStatus = models.OrderPaid }
	_ = h.Orders.UpdateStatus(w.OrderID, orderStatus)
	return c.SendStatus(http.StatusOK)
}
