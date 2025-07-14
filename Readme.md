# Fiber E-commerce Backend

A complete e-commerce backend built with Go Fiber, MongoDB, and Docker. This backend provides comprehensive CRUD operations for all e-commerce entities, user authentication, and file upload functionality.

## Features

- 🚀 **Fast & Lightweight**: Built with Go Fiber framework
- 🗄️ **MongoDB Database**: NoSQL database with proper indexing
- 🔐 **JWT Authentication**: Secure user authentication system
- 📁 **File Upload**: Image upload with URL response
- 🐳 **Docker Support**: Complete containerization
- 📊 **CRUD Operations**: Full CRUD for all models
- 📄 **Pagination**: Built-in pagination for all list endpoints
- 🔍 **Search**: Text search functionality for products
- 📋 **Validation**: Input validation and error handling

## Models Included

- Users (Authentication)
- Reviews
- TopCategory & Category
- Products
- Certificates & Licenses
- News & Partners
- Admin Management
- Currency
- Banners & Banner Sorting
- Contacts & Backgrounds
- Category Sorting

## Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.21+ (for local development)

### Using Docker (Recommended)

1. Clone the repository
2. Copy environment variables:
```bash
cp .env.example .env
```

3. Start the application:
```bash
make docker-up
```

The API will be available at `http://localhost:3000`

### Local Development

1. Install dependencies:
```bash
make deps
```

2. Start MongoDB (via Docker):
```bash
docker-compose up mongodb -d
```

3. Run in development mode:
```bash
make dev
```

## API Endpoints

### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login

### File Upload
- `POST /api/files/upload` - Upload single file
- `POST /api/files/upload-multiple` - Upload multiple files

### CRUD Endpoints (All Protected)
All CRUD endpoints follow the pattern: `/api/{resource}`

- `GET /api/{resource}` - Get all (with pagination)
- `GET /api/{resource}/:id` - Get single item
- `POST /api/{resource}` - Create new item
- `PUT /api/{resource}/:id` - Update item
- `DELETE /api/{resource}/:id` - Delete item

#### Available Resources:
- `reviews`
- `top-categories`
- `categories`
- `products`
- `sertificates`
- `licenses`
- `news`
- `partners`
- `admins`
- `currencies`
- `banners`
- `select-reviews`
- `backgrounds`
- `contacts`
- `banner-sorts`
- `top-category-sorts`
- `category-sorts`

### Query Parameters

#### Pagination
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10)

#### Filtering
- `category_id` - Filter products by category (products endpoint)
- `top_category_id` - Filter categories by top category (categories endpoint)

#### Search
- `search` - Text search in product names and descriptions (products endpoint)

## Authentication

The API uses JWT tokens for authentication. After login/register, include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## File Upload

Files are uploaded to the `/uploads` directory and served statically. The upload endpoint returns:

```json
{
  "url": "/uploads/filename.jpg",
  "filename": "unique-filename.jpg",
  "size": 12345
}
```

Supported file types: `.jpg`, `.jpeg`, `.png`, `.gif`, `.webp`, `.svg`

## Environment Variables

```bash
# Database
MONGODB_URI=mongodb://admin:password123@localhost:27017/ecommerce?authSource=admin

# Security
JWT_SECRET=your-super-secret-jwt-key-here

# Server
PORT=3000
APP_ENV=development
```

## Project Structure

```
.
├── config/           # Database configuration
├── middleware/       # Auth middleware
├── models/          # Data models
├── routes/          # API routes
├── uploads/         # File uploads directory
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── main.go
└── README.md
```

## Development Commands

```bash
# Start development server with auto-reload
make dev

# Build application
make build

# Run tests
make test

# Format code
make fmt

# Docker commands
make docker-up      # Start containers
make docker-down    # Stop containers
make docker-build   # Build images
make logs          # View logs

# Clean build artifacts
make clean
```

## API Examples

### Register User
```bash
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+998901234567",
    "password": "password123"
  }'
```

### Login
```bash
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "+998901234567",
    "password": "password123"
  }'
```

### Create Product (Protected)
```bash
curl -X POST http://localhost:3000/api/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>" \
  -d '{
    "name": "Product Name",
    "description": "Product description",
    "feature": "Product features",
    "price": "100.00",
    "brand": "Brand Name",
    "image": ["/uploads/image1.jpg"],
    "category_id": "category-object-id"
  }'
```

### Upload File (Protected)
```bash
curl -X POST http://localhost:3000/api/files/upload \
  -H "Authorization: Bearer <your-token>" \
  -F "file=@/path/to/your/image.jpg"
```

## Database Schema

The application automatically creates indexes for optimal performance:
- Unique indexes on user email and phone
- Text indexes for product search
- Foreign key indexes for relationships

## Security Features

- Password hashing with bcrypt
- JWT token authentication
- Input validation
- CORS enabled
- File type validation
- Request size limits

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License.