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
db.products.createIndex({ "name": "text", "description": "text" });
db.categories.createIndex({ "top_category_id": 1 });
db.products.createIndex({ "category_id": 1 });
db.banners.createIndex({ "top_category_id": 1 });
db.banners.createIndex({ "category_id": 1 });
db.banners.createIndex({ "product_id": 1 });

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
print('Single admin created:');
print('  Username: admin');
print('  Password: 123');