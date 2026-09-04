package routes

import (
	"context"
	"html"
	"regexp"
	"strconv"
	"strings"
	"time"

	"fiber-ecommerce/config"
	"fiber-ecommerce/middleware"
	"fiber-ecommerce/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const blogLanguageSeparator = "***"

var blogHTMLTagPattern = regexp.MustCompile(`<[^>]*>`)

func localizedValues(value string) []string {
	return strings.SplitN(value, blogLanguageSeparator, 3)
}

func localizedFieldIsComplete(value string, richText bool) bool {
	values := localizedValues(value)
	if len(values) != 3 {
		return false
	}

	for _, localizedValue := range values {
		content := localizedValue
		if richText {
			content = blogHTMLTagPattern.ReplaceAllString(content, " ")
			content = html.UnescapeString(content)
			content = strings.ReplaceAll(content, "\u00a0", " ")
		}
		if strings.TrimSpace(content) == "" {
			return false
		}
	}

	return true
}

func validateBlog(blog models.Blog) error {
	if !localizedFieldIsComplete(blog.Image, false) {
		return fiber.NewError(fiber.StatusBadRequest, "Blog image is required in English, Russian and Uzbek")
	}
	if !localizedFieldIsComplete(blog.Title, false) {
		return fiber.NewError(fiber.StatusBadRequest, "Blog title is required in English, Russian and Uzbek")
	}
	if !localizedFieldIsComplete(blog.Text, true) {
		return fiber.NewError(fiber.StatusBadRequest, "Blog text is required in English, Russian and Uzbek")
	}
	return nil
}

// BlogRoutes exposes public reads for the storefront and authenticated writes
// for the admin panel.
func BlogRoutes(app fiber.Router, db *mongo.Client) {
	blogs := app.Group("/blogs")
	collection := config.GetCollection(db, "blogs")

	blogs.Get("/", func(c *fiber.Ctx) error {
		page, err := strconv.Atoi(c.Query("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}

		limit, err := strconv.Atoi(c.Query("limit", "12"))
		if err != nil || limit < 1 {
			limit = 12
		}
		if limit > 100 {
			limit = 100
		}

		skip := int64((page - 1) * limit)
		opts := options.Find().
			SetSkip(skip).
			SetLimit(int64(limit)).
			SetSort(bson.D{{Key: "created_at", Value: -1}})

		cursor, err := collection.Find(c.Context(), bson.M{}, opts)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch blogs"})
		}
		defer cursor.Close(context.Background())

		blogList := make([]models.Blog, 0)
		if err := cursor.All(c.Context(), &blogList); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode blogs"})
		}

		total, err := collection.CountDocuments(c.Context(), bson.M{})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to count blogs"})
		}

		return c.JSON(fiber.Map{
			"data":  blogList,
			"total": total,
			"page":  page,
			"limit": limit,
		})
	})

	blogs.Get("/:id", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid blog ID"})
		}

		var blog models.Blog
		if err := collection.FindOne(c.Context(), bson.M{"_id": id}).Decode(&blog); err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Blog not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch blog"})
		}

		return c.JSON(blog)
	})

	protected := blogs.Group("", middleware.AdminJWTMiddleware())

	protected.Post("/", func(c *fiber.Ctx) error {
		var blog models.Blog
		if err := c.BodyParser(&blog); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}
		if err := validateBlog(blog); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		now := time.Now().UTC()
		blog.ID = primitive.NewObjectID()
		blog.CreatedAt = now
		blog.UpdatedAt = now

		if _, err := collection.InsertOne(c.Context(), blog); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create blog"})
		}

		return c.Status(fiber.StatusCreated).JSON(blog)
	})

	protected.Put("/:id", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid blog ID"})
		}

		var blog models.Blog
		if err := c.BodyParser(&blog); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}
		if err := validateBlog(blog); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		update := bson.M{"$set": bson.M{
			"title":      blog.Title,
			"text":       blog.Text,
			"image":      blog.Image,
			"updated_at": time.Now().UTC(),
		}}
		result, err := collection.UpdateOne(c.Context(), bson.M{"_id": id}, update)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update blog"})
		}
		if result.MatchedCount == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Blog not found"})
		}

		var updatedBlog models.Blog
		if err := collection.FindOne(c.Context(), bson.M{"_id": id}).Decode(&updatedBlog); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch updated blog"})
		}

		return c.JSON(updatedBlog)
	})

	protected.Delete("/:id", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid blog ID"})
		}

		result, err := collection.DeleteOne(c.Context(), bson.M{"_id": id})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete blog"})
		}
		if result.DeletedCount == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Blog not found"})
		}

		return c.JSON(fiber.Map{"message": "Blog deleted successfully"})
	})
}
