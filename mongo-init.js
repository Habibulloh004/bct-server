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

// Create indexes for better performance (without unique constraints for now)
db.users.createIndex({ "email": 1 });
db.users.createIndex({ "phone": 1 });
db.products.createIndex({ "name": "text", "description": "text" });
db.categories.createIndex({ "top_category_id": 1 });
db.products.createIndex({ "category_id": 1 });
db.banners.createIndex({ "top_category_id": 1 });
db.banners.createIndex({ "category_id": 1 });
db.banners.createIndex({ "product_id": 1 });

// Create a default admin user (optional)
db.admins.insertOne({
    name: "Super Admin",
    password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password: "password"
    created_at: new Date(),
    updated_at: new Date()
});

print('Database initialized successfully!');