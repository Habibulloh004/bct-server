package routes

import (
	"context"
	"log"
	"os"
	"time"

	"fiber-ecommerce/config"
	"fiber-ecommerce/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Admin authentication models
type AdminLoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type AdminUpdateRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type AdminAuthResponse struct {
	Token string       `json:"token"`
	Admin models.Admin `json:"admin"`
}

func AdminAuthRoutes(app fiber.Router, db *mongo.Client) {
	adminAuth := app.Group("/admin")

	// Admin Login - only login with existing admin
	adminAuth.Post("/login", func(c *fiber.Ctx) error {
		var req AdminLoginRequest
		if err := c.BodyParser(&req); err != nil {
			log.Printf("Error parsing request body: %v", err)
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		log.Printf("Login attempt for admin: %s", req.Name)

		if req.Name == "" || req.Password == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Name and password are required"})
		}

		collection := config.GetCollection(db, "admins")

		// Find the admin by name
		var admin models.Admin
		err := collection.FindOne(context.TODO(), bson.M{"name": req.Name}).Decode(&admin)
		if err != nil {
			log.Printf("Admin not found: %s, error: %v", req.Name, err)
			return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		log.Printf("Found admin: %s, stored password hash: %s", admin.Name, admin.Password)
		log.Printf("Input password: %s", req.Password)

		// Check password
		err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password))
		if err != nil {
			log.Printf("Password comparison failed for admin %s: %v", admin.Name, err)
			
			// Debug: Generate hash for the input password to compare
			testHash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			log.Printf("Generated hash for input password '%s': %s", req.Password, string(testHash))
			
			return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		log.Printf("Password verification successful for admin: %s", admin.Name)

		admin.Password = "" // Don't return password

		// Generate JWT token
		token, err := generateAdminJWT(admin.ID.Hex(), admin.Name)
		if err != nil {
			log.Printf("Failed to generate JWT token: %v", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
		}

		log.Printf("Login successful for admin: %s", admin.Name)

		return c.JSON(AdminAuthResponse{
			Token: token,
			Admin: admin,
		})
	})

	// Admin Update - update existing admin (protected route)
	adminAuth.Put("/update", func(c *fiber.Ctx) error {
		var req AdminUpdateRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		if req.Name == "" || req.Password == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Name and password are required"})
		}

		collection := config.GetCollection(db, "admins")

		// Find the single admin
		var admin models.Admin
		err := collection.FindOne(context.TODO(), bson.M{}).Decode(&admin)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Admin not found"})
		}

		// Hash new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
		}

		// Update admin
		updateData := bson.M{
			"name":       req.Name,
			"password":   string(hashedPassword),
			"updated_at": time.Now(),
		}

		_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": admin.ID}, bson.M{"$set": updateData})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update admin"})
		}

		// Get updated admin
		err = collection.FindOne(context.TODO(), bson.M{"_id": admin.ID}).Decode(&admin)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch updated admin"})
		}

		admin.Password = "" // Don't return password

		// Generate new JWT token with updated info
		token, err := generateAdminJWT(admin.ID.Hex(), admin.Name)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
		}

		return c.JSON(AdminAuthResponse{
			Token: token,
			Admin: admin,
		})
	})

	// Admin Profile - get current admin info (protected route)
	adminAuth.Get("/profile", func(c *fiber.Ctx) error {
		collection := config.GetCollection(db, "admins")

		// Find the single admin
		var admin models.Admin
		err := collection.FindOne(context.TODO(), bson.M{}).Decode(&admin)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Admin not found"})
		}

		admin.Password = "" // Don't return password
		return c.JSON(admin)
	})

	// Debug endpoint to check admin existence and password hash
	adminAuth.Get("/debug", func(c *fiber.Ctx) error {
		collection := config.GetCollection(db, "admins")

		var admin models.Admin
		err := collection.FindOne(context.TODO(), bson.M{"name": "admin"}).Decode(&admin)
		if err != nil {
			return c.JSON(fiber.Map{
				"error":       "Admin not found",
				"admin_count": getAdminCount(collection),
			})
		}

		// Test password hash
		testPassword := "bct123"
		err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(testPassword))
		passwordValid := err == nil

		return c.JSON(fiber.Map{
			"admin_exists":    true,
			"admin_name":      admin.Name,
			"password_hash":   admin.Password,
			"test_password":   testPassword,
			"password_valid":  passwordValid,
			"admin_count":     getAdminCount(collection),
		})
	})
}

func getAdminCount(collection *mongo.Collection) int64 {
	count, _ := collection.CountDocuments(context.TODO(), bson.M{})
	return count
}

func generateAdminJWT(adminID, adminName string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-super-secret-jwt-key-here"
	}

	claims := jwt.MapClaims{
		"admin_id":   adminID,
		"admin_name": adminName,
		"type":       "admin",
		"exp":        time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}