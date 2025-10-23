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

type mpesaWebhook struct {
	TransactionID string `json:"transaction_id"`
	OrderID       uint   `json:"order_id,string"`
	Status        string `json:"status"`
	AmountCents   int64  `json:"amount_cents"`
}

func TestWebhook_IdempotentOnDuplicateSuccess(t *testing.T) {
	db := testutil.OpenTestDB(t)
	if err := db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.OrderItem{}, &models.Transaction{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	orderRepo := repositories.NewOrderRepository(db)
	txRepo := repositories.NewTransactionRepository(db)

	// Seed order + tx pending
	o := &models.Order{UserID: 1, Status: models.OrderPending, TotalCents: 4500}
	if err := orderRepo.CreateWithItems(o, nil); err != nil { t.Fatalf("seed order: %v", err) }
	tx := &models.Transaction{OrderID: o.ID, ProviderRef: "tx-dupe", Status: models.TxPending, AmountCents: 4500}
	if err := txRepo.Create(tx); err != nil { t.Fatalf("seed tx: %v", err) }

	app := fiber.New()
	h := NewWebhookController(orderRepo, txRepo)
	app.Post("/api/webhooks/mpesa", h.Mpesa)

	payload := mpesaWebhook{TransactionID: "tx-dupe", OrderID: o.ID, Status: "SUCCESS", AmountCents: 4500}
	buf, _ := json.Marshal(payload)

	// First time
	req1 := httptest.NewRequest(http.MethodPost, "/api/webhooks/mpesa", bytes.NewReader(buf))
	req1.Header.Set("Content-Type", "application/json")
	resp1, err := app.Test(req1)
	if err != nil { t.Fatalf("fiber test 1: %v", err) }
	if resp1.StatusCode != http.StatusOK { t.Fatalf("want 200 got %d", resp1.StatusCode) }

	// Second time (duplicate)
	req2 := httptest.NewRequest(http.MethodPost, "/api/webhooks/mpesa", bytes.NewReader(buf))
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := app.Test(req2)
	if err != nil { t.Fatalf("fiber test 2: %v", err) }
	if resp2.StatusCode != http.StatusOK { t.Fatalf("want 200 got %d", resp2.StatusCode) }

	// Verify final statuses
	var gotTx models.Transaction
	if err := db.First(&gotTx, tx.ID).Error; err != nil { t.Fatalf("read tx: %v", err) }
	if gotTx.Status != models.TxSuccess { t.Fatalf("tx status want %s got %s", models.TxSuccess, gotTx.Status) }
	var gotOrder models.Order
	if err := db.First(&gotOrder, o.ID).Error; err != nil { t.Fatalf("read order: %v", err) }
	if gotOrder.Status != models.OrderPaid { t.Fatalf("order status want %s got %s", models.OrderPaid, gotOrder.Status) }
}
