// MongoDB initialization script
db = db.getSiblingDB('ecommerce');

// Create collections for existing models
db.createCollection('users');
db.createCollection('reviews');
db.createCollection('topcategories');
db.createCollection('categories');
db.createCollection('products');
db.createCollection('sertificates');
db.createCollection('licenses');
db.createCollection('news');
db.createCollection('partners');
db.createCollection('admins');
db.createCollection('currencies');
db.createCollection('banners');
db.createCollection('select_reviews');
db.createCollection('backgrounds');
db.createCollection('contacts');
db.createCollection('banner_sorts');
db.createCollection('top_category_sorts');
db.createCollection('category_sorts');

// Create collections for new models from schema diagram
db.createCollection('clients');
db.createCollection('orders');
db.createCollection('about');
db.createCollection('vendors');
db.createCollection('projects');
db.createCollection('links');

// Create indexes for better performance

// User indexes
db.users.createIndex({ "email": 1 }, { unique: true });
db.users.createIndex({ "phone": 1 }, { unique: true });

// Client indexes
db.clients.createIndex({ "email": 1 });
db.clients.createIndex({ "phone": 1 });

// Product search indexes
db.products.createIndex({ "name": "text", "description": "text", "ads_title": "text" });

// Relationship indexes
db.categories.createIndex({ "top_category_id": 1 });
db.products.createIndex({ "category_id": 1 });
db.orders.createIndex({ "client_id": 1 });
db.orders.createIndex({ "created_at": -1 });

// Banner indexes
db.banners.createIndex({ "top_category_id": 1 });
db.banners.createIndex({ "category_id": 1 });
db.banners.createIndex({ "product_id": 1 });
db.banners.createIndex({ "title": "text", "description": "text" });

// Sort indexes
db.banner_sorts.createIndex({ "unique_id": 1 });
db.banner_sorts.createIndex({ "banner_id": 1 });
db.top_category_sorts.createIndex({ "unique_id": 1 });
db.top_category_sorts.createIndex({ "top_category_id": 1 });
db.category_sorts.createIndex({ "unique_id": 1 });
db.category_sorts.createIndex({ "category_id": 1 });
db.category_sorts.createIndex({ "top_category_sort_id": 1 });

// Review indexes
db.select_reviews.createIndex({ "review_id": 1 });

// Time-based indexes
db.reviews.createIndex({ "created_at": -1 });
db.news.createIndex({ "created_at": -1 });
db.partners.createIndex({ "created_at": -1 });

// Clear existing admin and create the single default admin
db.admins.deleteMany({});

// Create the single admin user
// Username: "admin", Password: "123"
// This hash is generated with bcrypt cost 10 for password "123"
db.admins.insertOne({
    name: "admin",
    password: "$2a$10$HLjC0Amd/oTcHvQdhwzyguApEnT2n9XThdJXW.Ib1cBZveRAUe6T2", // bcrypt hash for "123"
    created_at: new Date(),
    updated_at: new Date()
});

print('Database initialized successfully!');
print('Collections created for all models from schema diagram:');
print('  - clients, orders, about, vendors, projects, links');
print('  - topcategories, categories, products');
print('  - reviews, sertificates, licenses, news, partners');
print('  - admins, currencies, banners, backgrounds, contacts');
print('  - banner_sorts, top_category_sorts, category_sorts, select_reviews');
print('');
print('Single admin created:');
print('  Username: admin');
print('  Password: 123');
