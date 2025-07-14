package routes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fiber-ecommerce/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

func FileRoutes(app fiber.Router, db *mongo.Client) {
	files := app.Group("/files")

	// Upload single file
	files.Post("/upload", func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "No file uploaded"})
		}

		// Validate file type
		allowedTypes := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".gif":  true,
			".webp": true,
			".svg":  true,
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !allowedTypes[ext] {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid file type"})
		}

		// Generate unique filename
		filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		uploadPath := filepath.Join("uploads", filename)

		// Create uploads directory if it doesn't exist
		if err := os.MkdirAll("uploads", 0755); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create upload directory"})
		}

		// Save file
		if err := c.SaveFile(file, uploadPath); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save file"})
		}

		// Return file URL
		fileURL := fmt.Sprintf("/uploads/%s", filename)
		
		return c.JSON(models.FileUploadResponse{
			URL:      fileURL,
			Filename: filename,
			Size:     file.Size,
		})
	})

	// Upload multiple files
	files.Post("/upload-multiple", func(c *fiber.Ctx) error {
		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Failed to parse form"})
		}

		files := form.File["files"]
		if len(files) == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "No files uploaded"})
		}

		var uploadedFiles []models.FileUploadResponse
		allowedTypes := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".gif":  true,
			".webp": true,
			".svg":  true,
		}

		// Create uploads directory if it doesn't exist
		if err := os.MkdirAll("uploads", 0755); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create upload directory"})
		}

		for _, file := range files {
			// Validate file type
			ext := strings.ToLower(filepath.Ext(file.Filename))
			if !allowedTypes[ext] {
				continue // Skip invalid files
			}

			// Generate unique filename
			filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
			uploadPath := filepath.Join("uploads", filename)

			// Save file
			if err := c.SaveFile(file, uploadPath); err != nil {
				continue // Skip files that failed to save
			}

			// Add to uploaded files list
			fileURL := fmt.Sprintf("/uploads/%s", filename)
			uploadedFiles = append(uploadedFiles, models.FileUploadResponse{
				URL:      fileURL,
				Filename: filename,
				Size:     file.Size,
			})
		}

		if len(uploadedFiles) == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "No valid files were uploaded"})
		}

		return c.JSON(fiber.Map{
			"files": uploadedFiles,
			"count": len(uploadedFiles),
		})
	})
}