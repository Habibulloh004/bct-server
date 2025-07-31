package main

import (
	"context"
	"log"
	"os"
	"time"

	"fiber-ecommerce/config"
	"fiber-ecommerce/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Initialize database
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func() {
		if err = db.Disconnect(context.TODO()); err != nil {
			log.Println("Failed to disconnect from database:", err)
		}
	}()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		BodyLimit:    50 * 1024 * 1024, // 50MB for file uploads
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	// Static files for uploads
	app.Static("/uploads", "./uploads")

	// API Routes
	api := app.Group("/api")

	// Authentication routes
	routes.AdminAuthRoutes(api, db)

	// Core models CRUD routes (updated models)
	routes.ClientRoutes(api, db)
	routes.TopCategoryRoutes(api, db)
	routes.CategoryRoutes(api, db)
	routes.ProductRoutes(api, db)
	routes.OrderRoutes(api, db)

	// Information pages (singleton models)
	routes.AboutRoutes(api, db)
	routes.LinksRoutes(api, db)

	// Media/Content routes
	routes.VendorRoutes(api, db)
	routes.ProjectRoutes(api, db)

	// Existing CRUD routes
	routes.ReviewRoutes(api, db)
	routes.SertificateRoutes(api, db)
	routes.LicenseRoutes(api, db)
	routes.NewsRoutes(api, db)
	routes.PartnerRoutes(api, db)
	routes.AdminRoutes(api, db)
	routes.CurrencyRoutes(api, db)
	routes.BannerRoutes(api, db)
	routes.SelectReviewRoutes(api, db)
	routes.BackgroundRoutes(api, db)
	routes.ContactsRoutes(api, db)
	routes.BannerSortRoutes(api, db)
	routes.TopCategorySortRoutes(api, db)
	routes.CategorySortRoutes(api, db)

	// File upload route
	routes.FileRoutes(api, db)

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Fiber E-commerce API",
			"version": "2.0.0",
			"status":  "running",
			"endpoints": fiber.Map{
				"health":        "/health",
				"api":           "/api/*",
				"uploads":       "/uploads/*",
				"admin_login":   "/api/admin/login",
				"documentation": "See README.md for full API documentation",
			},
		})
	})

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"message":   "Server is running",
			"timestamp": time.Now().UTC(),
			"database":  "connected",
		})
	})

	// 404 handler for debugging
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"error":   "Route not found",
			"path":    c.Path(),
			"method":  c.Method(),
			"message": "Please check the API documentation for available endpoints",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📚 API documentation: http://localhost:%s/", port)
	log.Printf("🔍 Health check: http://localhost:%s/health", port)
	log.Printf("🔐 Admin login: http://localhost:%s/api/admin/login", port)
	log.Printf("📁 File uploads: http://localhost:%s/uploads/", port)

	log.Fatal(app.Listen(":" + port))
}
