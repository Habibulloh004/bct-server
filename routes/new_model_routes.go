package routes

import (
	"context"
	"strconv"
	"time"

	"fiber-ecommerce/config"
	"fiber-ecommerce/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// Client CRUD
func ClientRoutes(app fiber.Router, db *mongo.Client) {
	clients := app.Group("/clients")

	// Get all clients
	clients.Get("/", func(c *fiber.Ctx) error {
		collection := config.GetCollection(db, "clients")

		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		skip := (page - 1) * limit

		opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "created_at", Value: -1}})

		cursor, err := collection.Find(context.TODO(), bson.M{}, opts)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch clients"})
		}
		defer cursor.Close(context.TODO())

		var clients []models.Client
		if err = cursor.All(context.TODO(), &clients); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to decode clients"})
		}

		// Remove passwords from response
		for i := range clients {
			clients[i].Password = ""
		}

		total, _ := collection.CountDocuments(context.TODO(), bson.M{})

		return c.JSON(fiber.Map{
			"data":  clients,
			"total": total,
			"page":  page,
			"limit": limit,
		})
	})

	// Get single client
	clients.Get("/:id", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
		}

		collection := config.GetCollection(db, "clients")
		var client models.Client
		err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&client)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Client not found"})
		}

		client.Password = "" // Don't return password
		return c.JSON(client)
	})

	// Create client
	clients.Post("/", func(c *fiber.Ctx) error {
		var client models.Client
		if err := c.BodyParser(&client); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Hash password if provided
		if client.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(client.Password), bcrypt.DefaultCost)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
			}
			client.Password = string(hashedPassword)
		}

		client.CreatedAt = time.Now()
		client.UpdatedAt = time.Now()
		collection := config.GetCollection(db, "clients")

		result, err := collection.InsertOne(context.TODO(), client)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create client"})
		}

		client.ID = result.InsertedID.(primitive.ObjectID)
		client.Password = "" // Don't return password
		return c.Status(201).JSON(client)
	})

	// Update client
	clients.Put("/:id", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
		}

		var updateData bson.M
		if err := c.BodyParser(&updateData); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Hash password if provided
		if password, exists := updateData["password"]; exists && password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password.(string)), bcrypt.DefaultCost)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
			}
			updateData["password"] = string(hashedPassword)
		}

		delete(updateData, "_id")
		updateData["updated_at"] = time.Now()

		collection := config.GetCollection(db, "clients")
		update := bson.M{"$set": updateData}

		result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update client"})
		}

		if result.MatchedCount == 0 {
			return c.Status(404).JSON(fiber.Map{"error": "Client not found"})
		}

		var client models.Client
		collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&client)
		client.Password = "" // Don't return password
		return c.JSON(client)
	})

	// Delete client
	clients.Delete("/:id", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
		}

		collection := config.GetCollection(db, "clients")
		result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete client"})
		}

		if result.DeletedCount == 0 {
			return c.Status(404).JSON(fiber.Map{"error": "Client not found"})
		}

		return c.JSON(fiber.Map{"message": "Client deleted successfully"})
	})
}

// Order CRUD
func OrderRoutes(app fiber.Router, db *mongo.Client) {
	orders := app.Group("/orders")

	// Get all orders
	orders.Get("/", func(c *fiber.Ctx) error {
		collection := config.GetCollection(db, "orders")

		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		skip := (page - 1) * limit

		filter := bson.M{}
		if clientID := c.Query("client_id"); clientID != "" {
			if id, err := primitive.ObjectIDFromHex(clientID); err == nil {
				filter["client_id"] = id
			}
		}

		opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "created_at", Value: -1}})

		cursor, err := collection.Find(context.TODO(), filter, opts)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch orders"})
		}
		defer cursor.Close(context.TODO())

		var orders []models.Order
		if err = cursor.All(context.TODO(), &orders); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to decode orders"})
		}

		total, _ := collection.CountDocuments(context.TODO(), filter)

		return c.JSON(fiber.Map{
			"data":  orders,
			"total": total,
			"page":  page,
			"limit": limit,
		})
	})

	// Get single order
	orders.Get("/:id", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
		}

		collection := config.GetCollection(db, "orders")
		var order models.Order
		err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&order)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Order not found"})
		}

		return c.JSON(order)
	})

	// Create order
	orders.Post("/", func(c *fiber.Ctx) error {
		var order models.Order
		if err := c.BodyParser(&order); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		order.CreatedAt = time.Now()
		order.UpdatedAt = time.Now()
		collection := config.GetCollection(db, "orders")

		result, err := collection.InsertOne(context.TODO(), order)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create order"})
		}

		order.ID = result.InsertedID.(primitive.ObjectID)
		return c.Status(201).JSON(order)
	})

	// Update order
	orders.Put("/:id", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
		}

		var updateData bson.M
		if err := c.BodyParser(&updateData); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		delete(updateData, "_id")
		updateData["updated_at"] = time.Now()

		collection := config.GetCollection(db, "orders")
		update := bson.M{"$set": updateData}

		result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update order"})
		}

		if result.MatchedCount == 0 {
			return c.Status(404).JSON(fiber.Map{"error": "Order not found"})
		}

		var order models.Order
		collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&order)
		return c.JSON(order)
	})

	// Delete order
	orders.Delete("/:id", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
		}

		collection := config.GetCollection(db, "orders")
		result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete order"})
		}

		if result.DeletedCount == 0 {
			return c.Status(404).JSON(fiber.Map{"error": "Order not found"})
		}

		return c.JSON(fiber.Map{"message": "Order deleted successfully"})
	})
}

// About CRUD
func AboutRoutes(app fiber.Router, db *mongo.Client) {
	about := app.Group("/about")

	// Get about info (single record)
	about.Get("/", func(c *fiber.Ctx) error {
		collection := config.GetCollection(db, "about")
		var aboutInfo models.About
		err := collection.FindOne(context.TODO(), bson.M{}).Decode(&aboutInfo)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(404).JSON(fiber.Map{"error": "About information not found"})
			}
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch about information"})
		}

		return c.JSON(aboutInfo)
	})

	// Create about info
	about.Post("/", func(c *fiber.Ctx) error {
		var aboutInfo models.About
		if err := c.BodyParser(&aboutInfo); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		aboutInfo.CreatedAt = time.Now()
		aboutInfo.UpdatedAt = time.Now()
		collection := config.GetCollection(db, "about")

		result, err := collection.InsertOne(context.TODO(), aboutInfo)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create about information"})
		}

		aboutInfo.ID = result.InsertedID.(primitive.ObjectID)
		return c.Status(201).JSON(aboutInfo)
	})

	// Update about info
	about.Put("/", func(c *fiber.Ctx) error {
		var updateData bson.M
		if err := c.BodyParser(&updateData); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		delete(updateData, "_id")
		updateData["updated_at"] = time.Now()

		collection := config.GetCollection(db, "about")
		update := bson.M{"$set": updateData}

		// Find and update the first (and should be only) record
		var aboutInfo models.About
		err := collection.FindOneAndUpdate(context.TODO(), bson.M{}, update).Decode(&aboutInfo)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(404).JSON(fiber.Map{"error": "About information not found"})
			}
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update about information"})
		}

		// Get updated record
		collection.FindOne(context.TODO(), bson.M{"_id": aboutInfo.ID}).Decode(&aboutInfo)
		return c.JSON(aboutInfo)
	})

	// Delete about info
	about.Delete("/", func(c *fiber.Ctx) error {
		collection := config.GetCollection(db, "about")
		result, err := collection.DeleteOne(context.TODO(), bson.M{})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete about information"})
		}

		if result.DeletedCount == 0 {
			return c.Status(404).JSON(fiber.Map{"error": "About information not found"})
		}

		return c.JSON(fiber.Map{"message": "About information deleted successfully"})
	})
}

// Vendor CRUD
func VendorRoutes(app fiber.Router, db *mongo.Client) {
	genericCRUD(app, db, "vendors", "vendors", models.Vendor{})
}

// Project CRUD
func ProjectRoutes(app fiber.Router, db *mongo.Client) {
	genericCRUD(app, db, "projects", "projects", models.Project{})
}

// Links CRUD
func LinksRoutes(app fiber.Router, db *mongo.Client) {
	links := app.Group("/links")

	// Get links (single record)
	links.Get("/", func(c *fiber.Ctx) error {
		collection := config.GetCollection(db, "links")
		var linksInfo models.Links
		err := collection.FindOne(context.TODO(), bson.M{}).Decode(&linksInfo)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(404).JSON(fiber.Map{"error": "Links information not found"})
			}
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch links information"})
		}

		return c.JSON(linksInfo)
	})

	// Create links
	links.Post("/", func(c *fiber.Ctx) error {
		var linksInfo models.Links
		if err := c.BodyParser(&linksInfo); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		linksInfo.CreatedAt = time.Now()
		linksInfo.UpdatedAt = time.Now()
		collection := config.GetCollection(db, "links")

		result, err := collection.InsertOne(context.TODO(), linksInfo)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create links information"})
		}

		linksInfo.ID = result.InsertedID.(primitive.ObjectID)
		return c.Status(201).JSON(linksInfo)
	})

	// Update links
	links.Put("/", func(c *fiber.Ctx) error {
		var updateData bson.M
		if err := c.BodyParser(&updateData); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		delete(updateData, "_id")
		updateData["updated_at"] = time.Now()

		collection := config.GetCollection(db, "links")
		update := bson.M{"$set": updateData}

		// Find and update the first (and should be only) record
		var linksInfo models.Links
		err := collection.FindOneAndUpdate(context.TODO(), bson.M{}, update).Decode(&linksInfo)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(404).JSON(fiber.Map{"error": "Links information not found"})
			}
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update links information"})
		}

		// Get updated record
		collection.FindOne(context.TODO(), bson.M{"_id": linksInfo.ID}).Decode(&linksInfo)
		return c.JSON(linksInfo)
	})

	// Delete links
	links.Delete("/", func(c *fiber.Ctx) error {
		collection := config.GetCollection(db, "links")
		result, err := collection.DeleteOne(context.TODO(), bson.M{})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete links information"})
		}

		if result.DeletedCount == 0 {
			return c.Status(404).JSON(fiber.Map{"error": "Links information not found"})
		}

		return c.JSON(fiber.Map{"message": "Links information deleted successfully"})
	})
}