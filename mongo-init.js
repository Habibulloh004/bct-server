// MongoDB initialization script
db = db.getSiblingDB('ecommerce');

// Create collections
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

// Create indexes for better performance
db.users.createIndex({ "email": 1 });
db.users.createIndex({ "phone": 1 });
db.admins.createIndex({ "name": 1 }, { unique: true });
db.products.createIndex({ "name": "text", "description": "text" });
db.categories.createIndex({ "top_category_id": 1 });
db.products.createIndex({ "category_id": 1 });
db.banners.createIndex({ "top_category_id": 1 });
db.banners.createIndex({ "category_id": 1 });
db.banners.createIndex({ "product_id": 1 });

// Create a default admin user
// Username: "admin", Password: "admin123"
db.admins.insertOne({
    name: "admin",
    password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // bcrypt hash for "admin123"
    created_at: new Date(),
    updated_at: new Date()
});

print('Database initialized successfully!');
print('Default admin created:');
print('  Username: admin');
print('  Password: admin123');