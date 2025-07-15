package routes

import (
	"context"
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
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		if req.Name == "" || req.Password == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Name and password are required"})
		}

		collection := config.GetCollection(db, "admins")

		// Find the single admin (there should only be one)
		var admin models.Admin
		err := collection.FindOne(context.TODO(), bson.M{"name": req.Name}).Decode(&admin)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		// Check password
		err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password))
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		admin.Password = "" // Don't return password

		// Generate JWT token
		token, err := generateAdminJWT(admin.ID.Hex(), admin.Name)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
		}

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