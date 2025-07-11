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
    is_active BOOLEAN DEFAULT TRUE
);

-- Create products table with arrays
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    tags TEXT[], -- Array of text
    scores INTEGER[], -- Array of integers
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert comprehensive test data for users
INSERT INTO users (name, email, age, created_at, is_active) VALUES
    ('John Doe', 'john@example.com', 30, '2024-07-01 10:00:00+00', TRUE),
    ('Jane Smith', 'jane@example.com', 25, '2024-06-15 15:30:00+00', TRUE),
    ('Bob Johnson', 'bob@example.com', 35, '2024-05-20 09:45:00+00', FALSE),
    ('Alice Brown', 'alice@example.com', 28, '2024-08-10 14:20:00+00', TRUE),
    ('Charlie Wilson', 'charlie@example.com', 42, '2024-03-12 11:15:00+00', TRUE),
    ('Diana Davis', 'diana@example.com', 23, '2024-09-01 16:45:00+00', TRUE),
    ('Eve Martinez', 'eve@example.com', 31, '2024-04-18 08:30:00+00', FALSE);

-- Insert comprehensive test data for products with arrays
INSERT INTO products (name, tags, scores, metadata, created_at) VALUES
    ('Smartphone Pro', ARRAY['electronics', 'mobile', 'gadgets'], ARRAY[95, 87, 92], '{"brand": "TechCorp", "category": "electronics"}', '2024-06-01 10:00:00+00'),
    ('Laptop Ultra', ARRAY['electronics', 'computers'], ARRAY[92, 95, 88], '{"brand": "CompuTech", "category": "electronics"}', '2024-05-15 14:30:00+00'),
    ('Mystery Novel', ARRAY['books', 'fiction', 'mystery'], ARRAY[85, 90, 87], '{"author": "John Author", "genre": "fiction"}', '2024-07-20 09:15:00+00'),
    ('Sports T-Shirt', ARRAY['clothing', 'apparel'], ARRAY[88, 91, 89], '{"size": "M", "color": "blue"}', '2024-08-05 12:45:00+00'),
    ('Gaming Console', ARRAY['electronics', 'gaming'], ARRAY[96, 94, 93], '{"brand": "GameTech", "category": "gaming"}', '2024-04-22 16:20:00+00');

-- Add some additional test data for edge cases
INSERT INTO users (name, email, age, created_at, is_active) VALUES
    ('Test User Old', 'old@example.com', 65, '2023-12-01 10:00:00+00', FALSE),
    ('Test User Young', 'young@example.com', 18, '2024-10-01 10:00:00+00', TRUE);

-- Add product with empty arrays for testing
INSERT INTO products (name, tags, scores, metadata, created_at) VALUES
    ('Empty Product', ARRAY[]::TEXT[], ARRAY[]::INTEGER[], '{"empty": true}', '2024-01-01 00:00:00+00');
