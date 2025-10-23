package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"mini-ecommerce/backend/internal/models"
	"mini-ecommerce/backend/internal/repositories"
	"mini-ecommerce/backend/internal/testutil"
)

type badMpesaWebhook struct {
	TransactionID string `json:"transaction_id"`
	OrderID       uint   `json:"order_id,string"`
	Status        string `json:"status"`
	AmountCents   int64  `json:"amount_cents"`
}

func TestWebhook_UnknownTransactionReturns400(t *testing.T) {
	db := testutil.OpenTestDB(t)
	if err := db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.OrderItem{}, &models.Transaction{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	orderRepo := repositories.NewOrderRepository(db)
	txRepo := repositories.NewTransactionRepository(db)

	// Seed order only, no transaction for given provider ref
	o := &models.Order{UserID: 1, Status: models.OrderPending, TotalCents: 1000}
	if err := orderRepo.CreateWithItems(o, nil); err != nil { t.Fatalf("seed order: %v", err) }

	app := fiber.New()
	h := NewWebhookController(orderRepo, txRepo)
	app.Post("/api/webhooks/mpesa", h.Mpesa)

	payload := badMpesaWebhook{TransactionID: "tx-missing", OrderID: o.ID, Status: "SUCCESS", AmountCents: 1000}
	buf, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/webhooks/mpesa", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil { t.Fatalf("fiber test: %v", err) }
	if resp.StatusCode != http.StatusBadRequest { t.Fatalf("want 400 got %d", resp.StatusCode) }
}
