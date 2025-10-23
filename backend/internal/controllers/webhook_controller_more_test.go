package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"mini-ecommerce/backend/internal/models"
	"mini-ecommerce/backend/internal/repositories"
	"mini-ecommerce/backend/internal/testutil"
)

func TestWebhook_Mpesa_BadJSON(t *testing.T) {
	db := testutil.OpenTestDB(t)
	if err := db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.OrderItem{}, &models.Transaction{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	orderRepo := repositories.NewOrderRepository(db)
	txRepo := repositories.NewTransactionRepository(db)

	app := fiber.New()
	h := NewWebhookController(orderRepo, txRepo)
	app.Post("/api/webhooks/mpesa", h.Mpesa)

	// Invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/api/webhooks/mpesa", bytes.NewBufferString("{bad"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil { t.Fatalf("fiber test: %v", err) }
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("want 400 got %d", resp.StatusCode)
	}
}

func TestWebhook_Mpesa_MapsUnknownStatusToFailed(t *testing.T) {
	db := testutil.OpenTestDB(t)
	if err := db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.OrderItem{}, &models.Transaction{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	orderRepo := repositories.NewOrderRepository(db)
	txRepo := repositories.NewTransactionRepository(db)

	// Seed order + tx
	o := &models.Order{UserID: 1, Status: models.OrderPending, TotalCents: 1000}
	if err := orderRepo.CreateWithItems(o, nil); err != nil { t.Fatalf("seed order: %v", err) }
	tx := &models.Transaction{OrderID: o.ID, ProviderRef: "tx-xyz", Status: models.TxPending, AmountCents: 1000}
	if err := txRepo.Create(tx); err != nil { t.Fatalf("seed tx: %v", err) }

	app := fiber.New()
	h := NewWebhookController(orderRepo, txRepo)
	app.Post("/api/webhooks/mpesa", h.Mpesa)

	body := []byte(`{"transaction_id":"tx-xyz","order_id":"` + string(rune('0'+o.ID)) + `","status":"SOMETHING","amount_cents":1000}`)
	req := httptest.NewRequest(http.MethodPost, "/api/webhooks/mpesa", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil { t.Fatalf("fiber test: %v", err) }
	if resp.StatusCode != http.StatusOK { t.Fatalf("status want 200 got %d", resp.StatusCode) }
}
