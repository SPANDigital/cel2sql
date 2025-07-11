-- Create comprehensive test data for CEL-to-SQL integration tests
-- This includes both users and products tables with rich data for testing
-- date arithmetic, array manipulation, and complex conditions

-- Create users table with comprehensive data
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    age INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    profile JSON,
    preferences JSONB
);

-- Create products table with arrays
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    tags TEXT[], -- Array of text
    scores INTEGER[], -- Array of integers
    metadata JSONB,
    details JSON,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert comprehensive test data for users
INSERT INTO users (name, email, age, created_at, is_active, profile, preferences) VALUES
    ('John Doe', 'john@example.com', 30, '2024-07-01 10:00:00+00', TRUE, 
     '{"bio": "Software developer", "location": "New York", "skills": ["JavaScript", "Python"]}',
     '{"theme": "dark", "notifications": true, "language": "en"}'),
    ('Jane Smith', 'jane@example.com', 25, '2024-06-15 15:30:00+00', TRUE,
     '{"bio": "Designer", "location": "San Francisco", "skills": ["UI/UX", "Photoshop"]}',
     '{"theme": "light", "notifications": false, "language": "en"}'),
    ('Bob Johnson', 'bob@example.com', 35, '2024-05-20 09:45:00+00', FALSE,
     '{"bio": "Manager", "location": "Chicago", "skills": ["Leadership", "Strategy"]}',
     '{"theme": "auto", "notifications": true, "language": "es"}'),
    ('Alice Brown', 'alice@example.com', 28, '2024-08-10 14:20:00+00', TRUE,
     '{"bio": "Data scientist", "location": "Seattle", "skills": ["Python", "R", "SQL"]}',
     '{"theme": "dark", "notifications": true, "language": "en"}'),
    ('Charlie Wilson', 'charlie@example.com', 42, '2024-03-12 11:15:00+00', TRUE,
     '{"bio": "Architect", "location": "Austin", "skills": ["Design", "Planning"]}',
     '{"theme": "light", "notifications": false, "language": "fr"}'),
    ('Diana Davis', 'diana@example.com', 23, '2024-09-01 16:45:00+00', TRUE,
     '{"bio": "Student", "location": "Boston", "skills": ["Learning", "Research"]}',
     '{"theme": "dark", "notifications": true, "language": "en"}'),
    ('Eve Martinez', 'eve@example.com', 31, '2024-04-18 08:30:00+00', FALSE,
     '{"bio": "Marketing", "location": "Miami", "skills": ["Marketing", "Social Media"]}',
     '{"theme": "auto", "notifications": false, "language": "es"}');

-- Insert comprehensive test data for products with arrays
INSERT INTO products (name, tags, scores, metadata, details, created_at) VALUES
    ('Smartphone Pro', ARRAY['electronics', 'mobile', 'gadgets'], ARRAY[95, 87, 92], 
     '{"brand": "TechCorp", "category": "electronics", "price": 999.99, "features": ["5G", "Camera", "GPS"]}',
     '{"warranty": "2 years", "color": "black", "storage": "256GB"}',
     '2024-06-01 10:00:00+00'),
    ('Laptop Ultra', ARRAY['electronics', 'computers'], ARRAY[92, 95, 88], 
     '{"brand": "CompuTech", "category": "electronics", "price": 1499.99, "features": ["SSD", "16GB RAM", "Intel i7"]}',
     '{"warranty": "3 years", "color": "silver", "screen": "15.6 inch"}',
     '2024-05-15 14:30:00+00'),
    ('Mystery Novel', ARRAY['books', 'fiction', 'mystery'], ARRAY[85, 90, 87], 
     '{"author": "John Author", "genre": "fiction", "price": 12.99, "pages": 320}',
     '{"publisher": "BookHouse", "language": "English", "isbn": "978-0123456789"}',
     '2024-07-20 09:15:00+00'),
    ('Sports T-Shirt', ARRAY['clothing', 'apparel'], ARRAY[88, 91, 89], 
     '{"size": "M", "color": "blue", "price": 29.99, "material": "cotton"}',
     '{"brand": "SportWear", "care": "machine wash", "fit": "regular"}',
     '2024-08-05 12:45:00+00'),
    ('Gaming Console', ARRAY['electronics', 'gaming'], ARRAY[96, 94, 93], 
     '{"brand": "GameTech", "category": "gaming", "price": 499.99, "features": ["4K", "HDR", "VR Ready"]}',
     '{"warranty": "1 year", "color": "white", "storage": "1TB"}',
     '2024-04-22 16:20:00+00');

-- Add some additional test data for edge cases
INSERT INTO users (name, email, age, created_at, is_active, profile, preferences) VALUES
    ('Test User Old', 'old@example.com', 65, '2023-12-01 10:00:00+00', FALSE,
     '{"bio": "Retired", "location": "Florida", "skills": ["Wisdom", "Experience"]}',
     '{"theme": "light", "notifications": false, "language": "en"}'),
    ('Test User Young', 'young@example.com', 18, '2024-10-01 10:00:00+00', TRUE,
     '{"bio": "Student", "location": "California", "skills": ["Learning", "Gaming"]}',
     '{"theme": "dark", "notifications": true, "language": "en"}');

-- Add product with empty arrays for testing
INSERT INTO products (name, tags, scores, metadata, details, created_at) VALUES
    ('Empty Product', ARRAY[]::TEXT[], ARRAY[]::INTEGER[], 
     '{"empty": true, "test": "data"}',
     '{"description": "Test product with empty arrays"}',
     '2024-01-01 00:00:00+00');
