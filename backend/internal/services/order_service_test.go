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

func TestCheckout_MultiItemTotals(t *testing.T) {
	db := testutil.OpenTestDB(t)
	migrate(t, db)

	prodRepo := repositories.NewProductRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	txRepo := repositories.NewTransactionRepository(db)

	// Seed two products
	p1 := &models.Product{Name: "P1", PriceCents: 1200, Stock: 10}
	p2 := &models.Product{Name: "P2", PriceCents: 350, Stock: 10}
	if err := prodRepo.Create(p1); err != nil { t.Fatalf("seed p1: %v", err) }
	if err := prodRepo.Create(p2); err != nil { t.Fatalf("seed p2: %v", err) }

	// Mock Mpesa
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]string{"transaction_id": "tx-multi", "status": "QUEUED"})
	}))
	defer ts.Close()

	svc := NewOrderService(prodRepo, orderRepo, txRepo, &MpesaClient{BaseURL: ts.URL, HTTP: ts.Client()})
	order, tx, err := svc.Checkout(42, []CheckoutItem{{ProductID: p1.ID, Qty: 3}, {ProductID: p2.ID, Qty: 5}})
	if err != nil { t.Fatalf("checkout: %v", err) }
	if order == nil || tx == nil { t.Fatalf("nil results") }

	// Expected total: 3*1200 + 5*350 = 3600 + 1750 = 5350
	if want := int64(5350); order.TotalCents != want { t.Fatalf("total want %d got %d", want, order.TotalCents) }

	// Ensure 2 order items saved
	var cnt int64
	if err := db.Model(&models.OrderItem{}).Where("order_id = ?", order.ID).Count(&cnt).Error; err != nil {
		t.Fatalf("count items: %v", err)
	}
	if cnt != 2 { t.Fatalf("want 2 items, got %d", cnt) }
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
