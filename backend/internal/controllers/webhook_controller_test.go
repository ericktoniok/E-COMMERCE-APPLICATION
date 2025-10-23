package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"mini-ecommerce/backend/internal/models"
	"mini-ecommerce/backend/internal/repositories"
	"mini-ecommerce/backend/internal/testutil"
)

func migrate(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.OrderItem{}, &models.Transaction{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
}

type mpesaWebhook struct {
	TransactionID string `json:"transaction_id"`
	OrderID       uint   `json:"order_id,string"`
	Status        string `json:"status"`
	AmountCents   int64  `json:"amount_cents"`
}

func TestWebhook_Mpesa_UpdatesOrderAndTx(t *testing.T) {
	db := testutil.OpenTestDB(t)
	migrate(t, db)

	orderRepo := repositories.NewOrderRepository(db)
	txRepo := repositories.NewTransactionRepository(db)

	// Seed order + transaction (pending)
	o := &models.Order{UserID: 1, Status: models.OrderPending, TotalCents: 3000}
	if err := orderRepo.CreateWithItems(o, nil); err != nil { t.Fatalf("seed order: %v", err) }
	tx := &models.Transaction{OrderID: o.ID, ProviderRef: "tx-abc", Status: models.TxPending, AmountCents: 3000}
	if err := txRepo.Create(tx); err != nil { t.Fatalf("seed tx: %v", err) }

	app := fiber.New()
	h := NewWebhookController(orderRepo, txRepo)
	app.Post("/api/webhooks/mpesa", h.Mpesa)

	payload := mpesaWebhook{TransactionID: "tx-abc", OrderID: o.ID, Status: "SUCCESS", AmountCents: 3000}
	buf, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/webhooks/mpesa", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil { t.Fatalf("fiber test: %v", err) }
	if resp.StatusCode != http.StatusOK { t.Fatalf("status want 200 got %d", resp.StatusCode) }

	// Verify updates
	var gotTx models.Transaction
	if err := db.First(&gotTx, tx.ID).Error; err != nil { t.Fatalf("read tx: %v", err) }
	if gotTx.Status != models.TxSuccess { t.Fatalf("tx status want %s got %s", models.TxSuccess, gotTx.Status) }

	var gotOrder models.Order
	if err := db.First(&gotOrder, o.ID).Error; err != nil { t.Fatalf("read order: %v", err) }
	if gotOrder.Status != models.OrderPaid { t.Fatalf("order status want %s got %s", models.OrderPaid, gotOrder.Status) }
}
