package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"mini-ecommerce/backend/internal/config"
	"mini-ecommerce/backend/internal/db"
	"mini-ecommerce/backend/internal/models"
	"mini-ecommerce/backend/internal/appctx"
	"mini-ecommerce/backend/internal/auth"
	"mini-ecommerce/backend/internal/controllers"
	"mini-ecommerce/backend/internal/middlewares"
	"mini-ecommerce/backend/internal/repositories"
	"mini-ecommerce/backend/internal/services"
	"mini-ecommerce/backend/internal/utils"
)

func main() {
	cfg := config.Load()
	port := cfg.Port

	app := fiber.New()

	// CORS for local dev (frontend on :5173)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders:  "*",
		AllowMethods:  "*",
	}))

	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Serve images from storage path
	app.Static("/images", cfg.ImageStoragePath)

	// DB connection and migrations
	dbConn := db.Connect(cfg.DatabaseURL)
	if err := dbConn.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.OrderItem{}, &models.Transaction{}); err != nil {
		log.Fatalf("migration error: %v", err)
	}

	// Seed admin user
	userRepo := repositories.NewUserRepository(dbConn)
	if cfg.AdminEmail != "" && cfg.AdminPassword != "" {
		if hash, err := utils.HashPassword(cfg.AdminPassword); err == nil {
			if _, err := userRepo.EnsureAdmin(cfg.AdminEmail, hash); err != nil {
				log.Printf("admin seed error: %v", err)
			}
		}
	}

	// Idempotent image enrichment for known products (only if missing)
	imgMap := map[string]string{
		"Sample Phone":       "https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?w=1200&q=80&auto=format&fit=crop",
		"Wireless Earbuds":   "https://images.unsplash.com/photo-1585386959984-a41552231658?w=1200&q=80&auto=format&fit=crop",
		"Laptop Sleeve":      "https://images.unsplash.com/photo-1527430253228-e93688616381?w=1200&q=80&auto=format&fit=crop",
		"Mechanical Keyboard": "https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=1200&q=80&auto=format&fit=crop",
		"USB‑C Hub":          "https://images.unsplash.com/photo-1612815154858-60aa4c59eaa0?w=1200&q=80&auto=format&fit=crop",
		"4K Monitor":         "https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=1200&q=80&auto=format&fit=crop",
		"Gaming Mouse":       "https://images.unsplash.com/photo-1555617117-08a1b4d8f8b9?w=1200&q=80&auto=format&fit=crop",
		"Portable SSD 1TB":   "https://images.unsplash.com/photo-1580910051074-3eb694886505?w=1200&q=80&auto=format&fit=crop",
		"Smartwatch":         "https://images.unsplash.com/photo-1519241047957-be31d7379a5d?w=1200&q=80&auto=format&fit=crop",
		"Desk Lamp":          "https://images.unsplash.com/photo-1493666438817-866a91353ca9?w=1200&q=80&auto=format&fit=crop",
	}
	for name, url := range imgMap {
		dbConn.Model(&models.Product{}).
			Where("name = ? AND (image_url = '' OR image_url IS NULL)", name).
			Update("image_url", url)
	}

	// App context and auth wiring
	jwtMgr := auth.NewJWTManager(cfg.JWTSecret, 72*time.Hour)
	ax := &appctx.Context{DB: dbConn, JWT: jwtMgr, Cfg: cfg}
	authSvc := services.NewAuthService(userRepo)
	authCtrl := controllers.NewAuthController(ax, authSvc)

	api := app.Group("/api")
	authGroup := api.Group("/auth")
	authGroup.Post("/register", authCtrl.Register)
	authGroup.Post("/login", authCtrl.Login)

	// Products
	prodRepo := repositories.NewProductRepository(dbConn)
	prodSvc := services.NewProductService(prodRepo)
	prodCtrl := controllers.NewProductController(ax, prodSvc)

	api.Get("/products", prodCtrl.List)
	admin := api.Group("/products", middlewares.JWTProtect(ax), middlewares.RequireRole(string(models.RoleAdmin)))
	admin.Post("/", prodCtrl.Create)
	admin.Put(":id", prodCtrl.Update)
	admin.Delete(":id", prodCtrl.Delete)
	admin.Post(":id/image", prodCtrl.UploadImage)

	// Seed sample products when requested or when table is empty
	var pc int64
	if err := dbConn.Model(&models.Product{}).Count(&pc).Error; err == nil {
		if os.Getenv("SEED") == "true" || pc == 0 {
			seed := []models.Product{
				{Name: "Sample Phone", Description: "6.1\" OLED, 128GB, great camera.", PriceCents: 59999, Stock: 12, Category: "Electronics", SKU: "PHN-001", Rating: 4.6},
				{Name: "Wireless Earbuds", Description: "Bluetooth 5.3, ANC, 24h battery.", PriceCents: 19999, Stock: 40, Category: "Audio", SKU: "EAR-201", Rating: 4.4},
				{Name: "Laptop Sleeve", Description: "Neoprene, fits 13\" laptops.", PriceCents: 2999, Stock: 50, Category: "Accessories", SKU: "SLV-133", Rating: 4.1},
				{Name: "Mechanical Keyboard", Description: "Hot‑swappable, RGB, brown switches.", PriceCents: 8999, Stock: 20, Category: "Peripherals", SKU: "KEY-882", Rating: 4.7},
				{Name: "USB‑C Hub", Description: "7‑in‑1: HDMI, SD, USB‑A, PD pass‑through.", PriceCents: 4599, Stock: 35, Category: "Peripherals", SKU: "HUB-701", Rating: 4.3},
				{Name: "4K Monitor", Description: "27\" IPS, 60Hz, thin bezels.", PriceCents: 229999, Stock: 8, Category: "Displays", SKU: "MON-270", Rating: 4.5},
				{Name: "Gaming Mouse", Description: "Ultra‑light, 6 programmable buttons.", PriceCents: 4999, Stock: 30, Category: "Peripherals", SKU: "MOU-510", Rating: 4.2},
				{Name: "Portable SSD 1TB", Description: "USB‑C, 1000MB/s, metal body.", PriceCents: 119999, Stock: 15, Category: "Storage", SKU: "SSD-1TB", Rating: 4.8},
				{Name: "Smartwatch", Description: "AMOLED, GPS, heart‑rate, 7‑day battery.", PriceCents: 89999, Stock: 18, Category: "Wearables", SKU: "SWT-009", Rating: 4.0},
				{Name: "Desk Lamp", Description: "LED, adjustable color temperature.", PriceCents: 2599, Stock: 60, Category: "Home Office", SKU: "LMP-120", Rating: 4.1},
			}
			if err := dbConn.Create(&seed).Error; err != nil {
				log.Printf("seed products error: %v", err)
			}
		}
	}

	// Orders & Checkout
	orderRepo := repositories.NewOrderRepository(dbConn)
	txRepo := repositories.NewTransactionRepository(dbConn)
	mpesaClient := &services.MpesaClient{BaseURL: cfg.MockMpesaURL, BackendWebhook: cfg.Port + "/api/webhooks/mpesa", HTTP: &http.Client{}}
	orderSvc := services.NewOrderService(prodRepo, orderRepo, txRepo, mpesaClient)
	orderCtrl := controllers.NewOrderController(ax, orderSvc)

	api.Post("/cart/checkout", middlewares.JWTProtect(ax), orderCtrl.Checkout)
	api.Get("/orders/me", middlewares.JWTProtect(ax), orderCtrl.ListMine)

	adminOrders := api.Group("/admin/orders", middlewares.JWTProtect(ax), middlewares.RequireRole(string(models.RoleAdmin)))
	adminOrders.Get("/", orderCtrl.AdminList)

	// Webhooks
	webhookCtrl := controllers.NewWebhookController(orderRepo, txRepo)
	api.Post("/webhooks/mpesa", webhookCtrl.Mpesa)

	log.Printf("API listening on :%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
