package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestCheckout_CreatesOrderAndTransaction(t *testing.T) {
	db := testutil.OpenTestDB(t)
	migrate(t, db)

	prodRepo := repositories.NewProductRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	txRepo := repositories.NewTransactionRepository(db)

	// Seed one product
	p := &models.Product{Name: "Test", PriceCents: 1500, Stock: 10}
	if err := prodRepo.Create(p); err != nil { t.Fatalf("seed product: %v", err) }

	// Mock Mpesa server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"transaction_id": "tx-123", "status": "QUEUED"})
	}))
	defer ts.Close()

	svc := NewOrderService(prodRepo, orderRepo, txRepo, &MpesaClient{BaseURL: ts.URL, HTTP: ts.Client()})
	order, tx, err := svc.Checkout(7, []CheckoutItem{{ProductID: p.ID, Qty: 2}})
	if err != nil { t.Fatalf("checkout error: %v", err) }
	if order == nil || tx == nil { t.Fatalf("expected order and tx, got nils") }
	if want := int64(3000); order.TotalCents != want { t.Fatalf("total cents want %d got %d", want, order.TotalCents) }

	// Verify transaction persisted
	var count int64
	if err := db.Model(&models.Transaction{}).Where("order_id = ? AND provider_ref = ?", order.ID, "tx-123").Count(&count).Error; err != nil {
		t.Fatalf("query tx: %v", err)
	}
	if count != 1 { t.Fatalf("expected 1 tx, got %d", count) }
}
