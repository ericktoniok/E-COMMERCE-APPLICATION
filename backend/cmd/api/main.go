package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
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
