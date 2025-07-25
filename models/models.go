package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User model for authentication
type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Phone     string             `json:"phone" bson:"phone"`
	Password  string             `json:"password" bson:"password"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type AuthResponse struct {
	// Token string `json:"token"`
	User  User   `json:"user"`
}

// Reviews model
type Reviews struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Phone     string             `json:"phone" bson:"phone"`
	Email     string             `json:"email" bson:"email"`
	Message   string             `json:"message" bson:"message"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

// TopCategory model
type TopCategory struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
	Categories []Category         `json:"categories,omitempty" bson:"categories,omitempty"`
}

// Category model
type Category struct {
	ID                primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Name              string              `json:"name" bson:"name"`
	Image             string              `json:"image" bson:"image"`
	CreatedAt         time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at" bson:"updated_at"`
	TopCategoryID     *primitive.ObjectID `json:"top_category_id" bson:"top_category_id"`
	Products          []Product           `json:"products,omitempty" bson:"products,omitempty"`
	TopCategorySortID *int                `json:"top_category_sort_id" bson:"top_category_sort_id"`
}

// Product model
type Product struct {
	ID          primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Name        string              `json:"name" bson:"name"`
	Description string              `json:"description" bson:"description"`
	Feature     string              `json:"feature" bson:"feature"`
	Price       string              `json:"price" bson:"price"`
	Brand       string              `json:"brand" bson:"brand"`
	Image       []string            `json:"image" bson:"image"`
	CreatedAt   time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at" bson:"updated_at"`
	CategoryID  *primitive.ObjectID `json:"category_id" bson:"category_id"`
}

// Sertificate model
type Sertificate struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Image     string             `json:"image" bson:"image"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// License model
type License struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Image     string             `json:"image" bson:"image"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// News model
type News struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Image     string             `json:"image" bson:"image"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// Partner model
type Partner struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Image     string             `json:"image" bson:"image"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// Admin model
type Admin struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Password  string             `json:"password" bson:"password"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// Currency model
type Currency struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Sum       string             `json:"sum" bson:"sum"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// Banner model
type Banner struct {
	ID            primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Image         string              `json:"image" bson:"image"`
	TopCategoryID *primitive.ObjectID `json:"top_category_id" bson:"top_category_id"`
	CategoryID    *primitive.ObjectID `json:"category_id" bson:"category_id"`
	ProductID     *primitive.ObjectID `json:"product_id" bson:"product_id"`
	CreatedAt     time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at" bson:"updated_at"`
}

// SelectReview model
type SelectReview struct {
	ID        primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	ReviewID  *primitive.ObjectID `json:"review_id" bson:"review_id"`
	Name      string              `json:"name" bson:"name"`
	Phone     string              `json:"phone" bson:"phone"`
	Email     string              `json:"email" bson:"email"`
	Message   string              `json:"message" bson:"message"`
	CreatedAt time.Time           `json:"created_at" bson:"created_at"`
}

// Background model
type Background struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Image     string             `json:"image" bson:"image"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// Contacts model
type Contacts struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CompanyName     string             `json:"company_name" bson:"company_name"`
	Phone1          string             `json:"phone1" bson:"phone1"`
	Phone2          string             `json:"phone2" bson:"phone2"`
	WorkHours       string             `json:"work_hours" bson:"work_hours"`
	Email           string             `json:"email" bson:"email"`
	Address         string             `json:"address" bson:"address"`
	Telegram        string             `json:"telegram" bson:"telegram"`
	TelegramBot     string             `json:"telegram_bot" bson:"telegram_bot"`
	Facebook        string             `json:"facebook" bson:"facebook"`
	Instagram       string             `json:"instagram" bson:"instagram"`
	Youtube         string             `json:"youtube" bson:"youtube"`
	FooterInfo      string             `json:"footer_info" bson:"footer_info"`
	ExperienceInfo  string             `json:"experience_info" bson:"experience_info"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}

// BannerSort model
type BannerSort struct {
	ID            primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	UniqueID      *int                `json:"unique_id" bson:"unique_id"`
	BannerID      *primitive.ObjectID `json:"banner_id" bson:"banner_id"`
	Image         string              `json:"image" bson:"image"`
	TopCategoryID *primitive.ObjectID `json:"top_category_id" bson:"top_category_id"`
	CategoryID    *primitive.ObjectID `json:"category_id" bson:"category_id"`
	ProductID     *primitive.ObjectID `json:"product_id" bson:"product_id"`
	CreatedAt     time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at" bson:"updated_at"`
}

// TopCategorySort model
type TopCategorySort struct {
	ID            primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Name          string              `json:"name" bson:"name"`
	TopCategoryID *primitive.ObjectID `json:"top_category_id" bson:"top_category_id"`
	UniqueID      *int                `json:"unique_id" bson:"unique_id"`
	CreatedAt     time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at" bson:"updated_at"`
}

// CategorySort model
type CategorySort struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UniqueID          int                `json:"unique_id" bson:"unique_id"`
	CategoryID        primitive.ObjectID `json:"category_id" bson:"category_id"`
	TopCategorySortID int                `json:"top_category_sort_id" bson:"top_category_sort_id"`
	Name              string             `json:"name" bson:"name"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
}

// File upload response
type FileUploadResponse struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}