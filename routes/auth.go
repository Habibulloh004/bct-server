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

func AuthRoutes(app fiber.Router, db *mongo.Client) {
	auth := app.Group("/auth")

	// Register
	auth.Post("/register", func(c *fiber.Ctx) error {
		var req models.RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request body",
				"details": err.Error(),
			})
		}

		// Log the received data for debugging
		println("Registration attempt:", req.Name, req.Email, req.Phone)

		// Validate required fields
		if req.Name == "" || req.Email == "" || req.Phone == "" || req.Password == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "All fields are required",
				"received": fiber.Map{
					"name": req.Name,
					"email": req.Email,
					"phone": req.Phone,
					"password_provided": req.Password != "",
				},
			})
		}

		collection := config.GetCollection(db, "users")

		// Check if user already exists by phone
		var existingUser models.User
		err := collection.FindOne(context.TODO(), bson.M{"phone": req.Phone}).Decode(&existingUser)
		if err == nil {
			return c.Status(400).JSON(fiber.Map{"error": "User with this phone already exists"})
		}

		// Check if user already exists by email
		err = collection.FindOne(context.TODO(), bson.M{"email": req.Email}).Decode(&existingUser)
		if err == nil {
			return c.Status(400).JSON(fiber.Map{"error": "User with this email already exists"})
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to hash password",
				"details": err.Error(),
			})
		}

		// Create user
		user := models.User{
			Name:      req.Name,
			Email:     req.Email,
			Phone:     req.Phone,
			Password:  string(hashedPassword),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		result, err := collection.InsertOne(context.TODO(), user)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to create user",
				"details": err.Error(),
			})
		}

		user.ID = result.InsertedID.(primitive.ObjectID)
		user.Password = "" // Don't return password

		// Generate JWT token
		// token, err := generateJWT(user.ID.Hex())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to generate token",
				"details": err.Error(),
			})
		}

		return c.Status(201).JSON(models.AuthResponse{
			// Token: token,
			User:  user,
		})
	})

	// Login
	auth.Post("/login", func(c *fiber.Ctx) error {
		var req models.LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		if req.Phone == "" || req.Password == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Phone and password are required"})
		}

		collection := config.GetCollection(db, "users")

		var user models.User
		err := collection.FindOne(context.TODO(), bson.M{"phone": req.Phone}).Decode(&user)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		// Check password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		user.Password = "" // Don't return password

		// Generate JWT token
		// token, err := generateJWT(user.ID.Hex())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
		}

		return c.JSON(models.AuthResponse{
			// Token: token,
			User:  user,
		})
	})
}

func generateJWT(userID string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-super-secret-jwt-key-here"
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}