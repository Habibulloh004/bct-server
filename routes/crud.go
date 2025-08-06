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
)

// Reviews CRUD
func ReviewRoutes(app fiber.Router, db *mongo.Client) {
	reviews := app.Group("/reviews")

	// Get all reviews
	reviews.Get("/", func(c *fiber.Ctx) error {
		collection := config.GetCollection(db, "reviews")

		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		skip := (page - 1) * limit

		opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "created_at", Value: -1}})

		cursor, err := collection.Find(context.TODO(), bson.M{}, opts)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch reviews"})
		}
		defer cursor.Close(context.TODO())

		var reviews []models.Reviews
		if err = cursor.All(context.TODO(), &reviews); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to decode reviews"})
		}

		// Get total count
		total, _ := collection.CountDocuments(context.TODO(), bson.M{})

		return c.JSON(fiber.Map{
			"data":  reviews,
			"total": total,
			"page":  page,
			"limit": limit,
		})
	})

	// Get single review
	reviews.Get("/:id", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
		}

		collection := config.GetCollection(db, "reviews")
		var review models.Reviews
		err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&review)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Review not found"})
		}

		return c.JSON(review)
	})

	// Create review
	reviews.Post("/", func(c *fiber.Ctx) error {
		var review models.Reviews
		if err := c.BodyParser(&review); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		review.CreatedAt = time.Now()
		collection := config.GetCollection(db, "reviews")

		result, err := collection.InsertOne(context.TODO(), review)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create review"})
		}

		review.ID = result.InsertedID.(primitive.ObjectID)
		return c.Status(201).JSON(review)
	})

	// Update review
	reviews.Put("/:id", func(c *fiber.Ctx) error {
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

		collection := config.GetCollection(db, "reviews")
		update := bson.M{"$set": updateData}

		result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update review"})
		}

		if result.MatchedCount == 0 {
			return c.Status(404).JSON(fiber.Map{"error": "Review not found"})
		}

		var review models.Reviews
		collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&review)
		return c.JSON(review)
	})

	// Delete review
	reviews.Delete("/:id", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
		}

		collection := config.GetCollection(db, "reviews")
		result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete review"})
		}

		if result.DeletedCount == 0 {
			return c.Status(404).JSON(fiber.Map{"error": "Review not found"})
		}

		return c.JSON(fiber.Map{"message": "Review deleted successfully"})
	})
}

// TopCategory CRUD
func TopCategoryRoutes(app fiber.Router, db *mongo.Client) {
	topCategories := app.Group("/top-categories")

	// Get all top categories
	topCategories.Get("/", func(c *fiber.Ctx) error {
		collection := config.GetCollection(db, "topcategories")

		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		skip := (page - 1) * limit

		opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "created_at", Value: -1}})

		cursor, err := collection.Find(context.TODO(), bson.M{}, opts)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch top categories"})
		}
		defer cursor.Close(context.TODO())

		var topCategories []models.TopCategory
		if err = cursor.All(context.TODO(), &topCategories); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to decode top categories"})
		}

		total, _ := collection.CountDocuments(context.TODO(), bson.M{})

		return c.JSON(fiber.Map{
			"data":  topCategories,
			"total": total,
			"page":  page,
			"limit": limit,
		})
	})

	// Get single top category
	topCategories.Get("/:id", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
		}

		collection := config.GetCollection(db, "topcategories")
		var topCategory models.TopCategory
		err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&topCategory)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Top category not found"})
		}

		return c.JSON(topCategory)
	})

	// Create top category
	topCategories.Post("/", func(c *fiber.Ctx) error {
		var topCategory models.TopCategory
		if err := c.BodyParser(&topCategory); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		topCategory.CreatedAt = time.Now()
		topCategory.UpdatedAt = time.Now()
		collection := config.GetCollection(db, "topcategories")

		result, err := collection.InsertOne(context.TODO(), topCategory)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create top category"})
		}

		topCategory.ID = result.InsertedID.(primitive.ObjectID)
		return c.Status(201).JSON(topCategory)
	})

	// Update top category
	topCategories.Put("/:id", func(c *fiber.Ctx) error {
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

		collection := config.GetCollection(db, "topcategories")
		update := bson.M{"$set": updateData}

		result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update top category"})
		}

		if result.MatchedCount == 0 {
			return c.Status(404).JSON(fiber.Map{"error": "Top category not found"})
		}

		var topCategory models.TopCategory
		collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&topCategory)
		return c.JSON(topCategory)
	})

	// Delete top category
	topCategories.Delete("/:id", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
		}

		collection := config.GetCollection(db, "topcategories")
		result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete top category"})
		}

		if result.DeletedCount == 0 {
			return c.Status(404).JSON(fiber.Map{"error": "Top category not found"})
		}

		return c.JSON(fiber.Map{"message": "Top category deleted successfully"})
	})
}

// Category CRUD
func CategoryRoutes(app fiber.Router, db *mongo.Client) {
	categories := app.Group("/categories")

	// Get all categories
	categories.Get("/", func(c *fiber.Ctx) error {
		collection := config.GetCollection(db, "categories")

		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		skip := (page - 1) * limit

		filter := bson.M{}
		if topCategoryID := c.Query("top_category_id"); topCategoryID != "" {
			if id, err := primitive.ObjectIDFromHex(topCategoryID); err == nil {
				filter["top_category_id"] = id
			}
		}

		opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "created_at", Value: -1}})

		cursor, err := collection.Find(context.TODO(), filter, opts)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch categories"})
		}
		defer cursor.Close(context.TODO())

		var categories []models.Category
		if err = cursor.All(context.TODO(), &categories); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to decode categories"})
		}

		total, _ := collection.CountDocuments(context.TODO(), filter)

		return c.JSON(fiber.Map{
			"data":  categories,
			"total": total,
			"page":  page,
			"limit": limit,
		})
	})

	// Get single category
	categories.Get("/:id", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
		}

		collection := config.GetCollection(db, "categories")
		var category models.Category
		err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&category)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
		}

		return c.JSON(category)
	})

	// Create category
	categories.Post("/", func(c *fiber.Ctx) error {
		var category models.Category
		if err := c.BodyParser(&category); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		category.CreatedAt = time.Now()
		category.UpdatedAt = time.Now()
		collection := config.GetCollection(db, "categories")

		result, err := collection.InsertOne(context.TODO(), category)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create category"})
		}

		category.ID = result.InsertedID.(primitive.ObjectID)
		return c.Status(201).JSON(category)
	})

	// Update category
	categories.Put("/:id", func(c *fiber.Ctx) error {
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

		collection := config.GetCollection(db, "categories")
		update := bson.M{"$set": updateData}

		result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update category"})
		}

		if result.MatchedCount == 0 {
			return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
		}

		var category models.Category
		collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&category)
		return c.JSON(category)
	})

	// Delete category
	categories.Delete("/:id", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
		}

		collection := config.GetCollection(db, "categories")
		result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete category"})
		}

		if result.DeletedCount == 0 {
			return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
		}

		return c.JSON(fiber.Map{"message": "Category deleted successfully"})
	})
}

// Product CRUD
func ProductRoutes(app fiber.Router, db *mongo.Client) {
	products := app.Group("/products")

	// Get all products with category names populated
	products.Get("/", func(c *fiber.Ctx) error {
		collection := config.GetCollection(db, "products")

		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		skip := (page - 1) * limit

		filter := bson.M{}
		if categoryID := c.Query("category_id"); categoryID != "" {
			if id, err := primitive.ObjectIDFromHex(categoryID); err == nil {
				filter["category_id"] = id
			}
		}
		if topCategoryID := c.Query("top_category_id"); topCategoryID != "" {
			if id, err := primitive.ObjectIDFromHex(topCategoryID); err == nil {
				filter["top_category_id"] = id
			}
		}

		// Search by name, ads_title, or description
		if search := c.Query("search"); search != "" {
			filter["$or"] = []bson.M{
				{"name": bson.M{"$regex": search, "$options": "i"}},
				{"ads_title": bson.M{"$regex": search, "$options": "i"}},
				{"description": bson.M{"$regex": search, "$options": "i"}},
			}
		}

		opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "created_at", Value: -1}})

		cursor, err := collection.Find(context.TODO(), filter, opts)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch products"})
		}
		defer cursor.Close(context.TODO())

		var products []models.Product
		if err = cursor.All(context.TODO(), &products); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to decode products"})
		}

		// Populate category and top category names
		for i := range products {
			populateCategoryNames(db, &products[i])
		}

		total, _ := collection.CountDocuments(context.TODO(), filter)

		return c.JSON(fiber.Map{
			"data":  products,
			"total": total,
			"page":  page,
			"limit": limit,
		})
	})

	// Get single product with category names populated
	products.Get("/:id", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
		}

		collection := config.GetCollection(db, "products")
		var product models.Product
		err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&product)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
		}

		// Populate category and top category names
		populateCategoryNames(db, &product)

		return c.JSON(product)
	})

	// Create product (populate top_category_id automatically)
	products.Post("/", func(c *fiber.Ctx) error {
		var product models.Product
		if err := c.BodyParser(&product); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Auto-populate top_category_id from category
		if product.CategoryID != nil {
			categoryCollection := config.GetCollection(db, "categories")
			var category models.Category
			err := categoryCollection.FindOne(context.TODO(), bson.M{"_id": product.CategoryID}).Decode(&category)
			if err == nil && category.TopCategoryID != nil {
				product.TopCategoryID = category.TopCategoryID
			}
		}

		product.CreatedAt = time.Now()
		product.UpdatedAt = time.Now()
		collection := config.GetCollection(db, "products")

		result, err := collection.InsertOne(context.TODO(), product)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create product"})
		}

		product.ID = result.InsertedID.(primitive.ObjectID)

		// Populate category names for response
		populateCategoryNames(db, &product)

		return c.Status(201).JSON(product)
	})

	// Update product
	products.Put("/:id", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
		}

		var updateData bson.M
		if err := c.BodyParser(&updateData); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Auto-populate top_category_id if category_id is being updated
		if categoryIDInterface, exists := updateData["category_id"]; exists {
			if categoryIDStr, ok := categoryIDInterface.(string); ok {
				if categoryID, err := primitive.ObjectIDFromHex(categoryIDStr); err == nil {
					categoryCollection := config.GetCollection(db, "categories")
					var category models.Category
					err := categoryCollection.FindOne(context.TODO(), bson.M{"_id": categoryID}).Decode(&category)
					if err == nil && category.TopCategoryID != nil {
						updateData["top_category_id"] = category.TopCategoryID
					}
				}
			}
		}

		delete(updateData, "_id")
		updateData["updated_at"] = time.Now()

		collection := config.GetCollection(db, "products")
		update := bson.M{"$set": updateData}

		result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update product"})
		}

		if result.MatchedCount == 0 {
			return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
		}

		var product models.Product
		collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&product)

		// Populate category names for response
		populateCategoryNames(db, &product)

		return c.JSON(product)
	})

	// Delete product
	products.Delete("/:id", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
		}

		collection := config.GetCollection(db, "products")
		result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete product"})
		}

		if result.DeletedCount == 0 {
			return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
		}

		return c.JSON(fiber.Map{"message": "Product deleted successfully"})
	})

	// Get products by top category
	products.Get("/by-top-category/:top_category_id", func(c *fiber.Ctx) error {
		topCategoryID, err := primitive.ObjectIDFromHex(c.Params("top_category_id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid top category ID"})
		}

		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		skip := (page - 1) * limit

		collection := config.GetCollection(db, "products")
		filter := bson.M{"top_category_id": topCategoryID}

		opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "created_at", Value: -1}})

		cursor, err := collection.Find(context.TODO(), filter, opts)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch products"})
		}
		defer cursor.Close(context.TODO())

		var products []models.Product
		if err = cursor.All(context.TODO(), &products); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to decode products"})
		}

		// Populate category names
		for i := range products {
			populateCategoryNames(db, &products[i])
		}

		total, _ := collection.CountDocuments(context.TODO(), filter)

		return c.JSON(fiber.Map{
			"data":  products,
			"total": total,
			"page":  page,
			"limit": limit,
		})
	})

	// Get discounted products
	products.Get("/discounted", func(c *fiber.Ctx) error {
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		skip := (page - 1) * limit

		collection := config.GetCollection(db, "products")
		filter := bson.M{
			"discount": bson.M{
				"$ne":     nil,
				"$exists": true,
				// "$ne":     "",
			},
		}

		opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "created_at", Value: -1}})

		cursor, err := collection.Find(context.TODO(), filter, opts)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch discounted products"})
		}
		defer cursor.Close(context.TODO())

		var products []models.Product
		if err = cursor.All(context.TODO(), &products); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to decode products"})
		}

		// Populate category names
		for i := range products {
			populateCategoryNames(db, &products[i])
		}

		total, _ := collection.CountDocuments(context.TODO(), filter)

		return c.JSON(fiber.Map{
			"data":  products,
			"total": total,
			"page":  page,
			"limit": limit,
		})
	})
}

// Helper function to populate category and top category names
func populateCategoryNames(db *mongo.Client, product *models.Product) {
	// Populate category name
	if product.CategoryID != nil {
		categoryCollection := config.GetCollection(db, "categories")
		var category models.Category
		err := categoryCollection.FindOne(context.TODO(), bson.M{"_id": product.CategoryID}).Decode(&category)
		if err == nil {
			product.CategoryName = category.Name
		}
	}

	// Populate top category name
	if product.TopCategoryID != nil {
		topCategoryCollection := config.GetCollection(db, "topcategories")
		var topCategory models.TopCategory
		err := topCategoryCollection.FindOne(context.TODO(), bson.M{"_id": product.TopCategoryID}).Decode(&topCategory)
		if err == nil {
			product.TopCategoryName = topCategory.Name
		}
	}
}
