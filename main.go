package main

import (
	"context"
	"log"
	"os"
	"time"

	"fiber-ecommerce/config"
	// "fiber-ecommerce/middleware"
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

	// Routes
	api := app.Group("/api")

	// Auth routes
	routes.AuthRoutes(api, db)

	// Protected routes
	// protected := api.Group("/", middleware.JWTMiddleware())

	// All CRUD routes
	routes.ReviewRoutes(api.Group("/"), db)
	routes.TopCategoryRoutes(api.Group("/"), db)
	routes.CategoryRoutes(api.Group("/"), db)
	routes.ProductRoutes(api.Group("/"), db)
	routes.SertificateRoutes(api.Group("/"), db)
	routes.LicenseRoutes(api.Group("/"), db)
	routes.NewsRoutes(api.Group("/"), db)
	routes.PartnerRoutes(api.Group("/"), db)
	routes.AdminRoutes(api.Group("/"), db)
	routes.CurrencyRoutes(api.Group("/"), db)
	routes.BannerRoutes(api.Group("/"), db)
	routes.SelectReviewRoutes(api.Group("/"), db)
	routes.BackgroundRoutes(api.Group("/"), db)
	routes.ContactsRoutes(api.Group("/"), db)
	routes.BannerSortRoutes(api.Group("/"), db)
	routes.TopCategorySortRoutes(api.Group("/"), db)
	routes.CategorySortRoutes(api.Group("/"), db)

	// File upload route
	routes.FileRoutes(api.Group("/"), db)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
