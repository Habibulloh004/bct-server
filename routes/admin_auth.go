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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Admin authentication models
type AdminLoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type AdminRegisterRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type AdminAuthResponse struct {
	Token string       `json:"token"`
	Admin models.Admin `json:"admin"`
}

func AdminAuthRoutes(app fiber.Router, db *mongo.Client) {
	adminAuth := app.Group("/admin")

	// Admin Register
	adminAuth.Post("/register", func(c *fiber.Ctx) error {
		var req AdminRegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
		}

		// Validate required fields
		if req.Name == "" || req.Password == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Name and password are required",
			})
		}

		collection := config.GetCollection(db, "admins")

		// Check if admin already exists by name
		var existingAdmin models.Admin
		err := collection.FindOne(context.TODO(), bson.M{"name": req.Name}).Decode(&existingAdmin)
		if err == nil {
			return c.Status(400).JSON(fiber.Map{"error": "Admin with this name already exists"})
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":   "Failed to hash password",
				"details": err.Error(),
			})
		}

		// Create admin
		admin := models.Admin{
			Name:      req.Name,
			Password:  string(hashedPassword),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		result, err := collection.InsertOne(context.TODO(), admin)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":   "Failed to create admin",
				"details": err.Error(),
			})
		}

		admin.ID = result.InsertedID.(primitive.ObjectID)
		admin.Password = "" // Don't return password

		// Generate JWT token
		token, err := generateAdminJWT(admin.ID.Hex(), admin.Name)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":   "Failed to generate token",
				"details": err.Error(),
			})
		}

		return c.Status(201).JSON(AdminAuthResponse{
			Token: token,
			Admin: admin,
		})
	})

	// Admin Login
	adminAuth.Post("/login", func(c *fiber.Ctx) error {
		var req AdminLoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		if req.Name == "" || req.Password == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Name and password are required"})
		}

		collection := config.GetCollection(db, "admins")

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

	// Admin Profile (Protected)
	adminAuth.Get("/profile", func(c *fiber.Ctx) error {
		// Extract admin ID from JWT token
		adminID := c.Locals("admin_id")
		if adminID == nil {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
		}

		id, err := primitive.ObjectIDFromHex(adminID.(string))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid admin ID"})
		}

		collection := config.GetCollection(db, "admins")
		var admin models.Admin
		err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&admin)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Admin not found"})
		}

		admin.Password = "" // Don't return password
		return c.JSON(admin)
	})

	// Admin Logout (Optional - mainly for client-side token removal)
	adminAuth.Post("/logout", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Logged out successfully"})
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
