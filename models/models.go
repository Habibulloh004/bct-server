package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User model for authentication
type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Phone     string             `json:"phone" bson:"phone"`
	Password  string             `json:"password" bson:"password"`
	IsActive  bool               `json:"is_active" bson:"is_active"`
	LastLogin *time.Time         `json:"last_login,omitempty" bson:"last_login,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// Updated login request to use phone instead of email
type UserLoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// Updated register request - phone is required, email is optional
type UserRegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type UserAuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// Client model (from schema diagram)
type Client struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	Phone     string             `json:"phone" bson:"phone"`
	Image     string             `json:"image" bson:"image"`
	URL       string             `json:"url" bson:"url"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// TopCategory model (updated)
type TopCategory struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Image     string             `json:"image" bson:"image"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// Category model (updated with top_category_name)
type Category struct {
	ID              primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Name            string              `json:"name" bson:"name"`
	Image           string              `json:"image" bson:"image"`
	TopCategoryID   *primitive.ObjectID `json:"top_category_id" bson:"top_category_id"`
	TopCategoryName string              `json:"top_category_name,omitempty" bson:"top_category_name,omitempty"`
	CreatedAt       time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at" bson:"updated_at"`
}

// Product model (updated with all required fields)
type Product struct {
	ID              primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Name            string              `json:"name" bson:"name"`
	AdsTitle        string              `json:"ads_title" bson:"ads_title"`
	Image           []string            `json:"image" bson:"image"`
	Description     string              `json:"description" bson:"description"`
	Guarantee       string              `json:"guarantee" bson:"guarantee"`
	SerialNumber    string              `json:"serial_number" bson:"serial_number"`
	Price           string              `json:"price" bson:"price"`
	Discount        string              `json:"discount,omitempty" bson:"discount,omitempty"`
	Active          bool                `json:"active" bson:"active"`
	Index           int                 `json:"index" bson:"index"`
	CategoryID      *primitive.ObjectID `json:"category_id" bson:"category_id"`
	TopCategoryID   *primitive.ObjectID `json:"top_category_id" bson:"top_category_id"`
	CategoryName    string              `json:"category_name,omitempty" bson:"category_name,omitempty"`
	TopCategoryName string              `json:"top_category_name,omitempty" bson:"top_category_name,omitempty"`
	CreatedAt       time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at" bson:"updated_at"`
}

// Order model (from schema diagram)
type Order struct {
	ID                primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Phone             string              `json:"phone" bson:"phone"`
	PayType           string              `json:"pay_type" bson:"pay_type"`
	ProductsWithCount []ProductWithCount  `json:"products" bson:"products"`
	ClientID          *primitive.ObjectID `json:"client_id" bson:"client_id"`
	CreatedAt         time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at" bson:"updated_at"`
}

type ProductWithCount struct {
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	Count     int                `json:"count" bson:"count"`
}

// About model (from schema diagram)
type About struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Creation     string             `json:"creation" bson:"creation"`
	Clients      string             `json:"clients" bson:"clients"`
	Partners     string             `json:"partners" bson:"partners"`
	Technologies string             `json:"technologies" bson:"technologies"`
	Scaners      string             `json:"scaners" bson:"scaners"`
	Scales       string             `json:"scales" bson:"scales"`
	Printers     string             `json:"printers" bson:"printers"`
	Cashiers     string             `json:"cashiers" bson:"cashiers"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

// Vendor model (from schema diagram)
type Vendor struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Image     string             `json:"image" bson:"image"`
	URL       string             `json:"url" bson:"url"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// Project model (from schema diagram)
type Project struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Image     string             `json:"image" bson:"image"`
	URL       string             `json:"url" bson:"url"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// Links model (from schema diagram)
type Links struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Facebook  string             `json:"facebook" bson:"facebook"`
	Instagram string             `json:"instagram" bson:"instagram"`
	LinkedIn  string             `json:"linkedin" bson:"linkedin"`
	YouTube   string             `json:"youtube" bson:"youtube"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// Existing models from your code (keeping as they are)

// Reviews model
type Reviews struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Phone     string             `json:"phone" bson:"phone"`
	Email     string             `json:"email" bson:"email"`
	Message   string             `json:"message" bson:"message"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
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
	CreatedAt time.Time          `json:"created_at" bsom:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// Banner model
type Banner struct {
	ID            primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Image         string              `json:"image" bson:"image"`
	Title         string              `json:"title" bson:"title"`
	Description   string              `json:"description" bson:"description"`
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
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CompanyName    string             `json:"company_name" bson:"company_name"`
	Phone1         string             `json:"phone1" bson:"phone1"`
	Phone2         string             `json:"phone2" bson:"phone2"`
	WorkHours      string             `json:"work_hours" bson:"work_hours"`
	Email          string             `json:"email" bson:"email"`
	Address        string             `json:"address" bson:"address"`
	Telegram       string             `json:"telegram" bson:"telegram"`
	TelegramBot    string             `json:"telegram_bot" bson:"telegram_bot"`
	Facebook       string             `json:"facebook" bson:"facebook"`
	Instagram      string             `json:"instagram" bson:"instagram"`
	Youtube        string             `json:"youtube" bson:"youtube"`
	FooterInfo     string             `json:"footer_info" bson:"footer_info"`
	ExperienceInfo string             `json:"experience_info" bson:"experience_info"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
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

// Discount model (singleton)
type Discount struct {
    ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Title     string             `json:"title" bson:"title"`
    ProductID string             `json:"product_id" bson:"product_id"`
    CreatedAt time.Time          `json:"created_at" bson:"created_at"`
    UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// Vendors_about model
type Vendors_about struct {
    ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name        string             `json:"name" bson:"name"`
    Description string             `json:"description" bson:"description"`
    Image       string             `json:"image" bson:"image"`
    CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
    UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// Official_partner model (singleton)
type Official_partner struct {
    ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name        string             `json:"name" bson:"name"`
    Image       string             `json:"image" bson:"image"`
    Description string             `json:"description" bson:"description"`
    CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
    UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// Experiments model
type Experiments struct {
    ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Count       string             `json:"count" bson:"count"`
    Title       string             `json:"title" bson:"title"`
    Description string             `json:"description" bson:"description"`
    CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
    UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// Company_stats model
type Company_stats struct {
    ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Count       string             `json:"count" bson:"count"`
    Title       string             `json:"title" bson:"title"`
    Description string             `json:"description" bson:"description"`
    Image       string             `json:"image" bson:"image"`
    CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
    UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

