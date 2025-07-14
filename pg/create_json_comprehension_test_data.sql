-- Create test tables for JSON/JSONB comprehension testing
-- This includes tables with JSON arrays and complex nested structures

-- Create users table with JSON arrays and nested objects
CREATE TABLE json_users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    settings JSONB DEFAULT '{}',
    metadata JSON DEFAULT '{}',
    tags JSONB DEFAULT '[]', -- JSON array of strings
    scores JSONB DEFAULT '[]', -- JSON array of numbers
    attributes JSON DEFAULT '[]', -- JSON array of objects
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create products table with complex JSON structures
CREATE TABLE json_products (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    features JSONB DEFAULT '[]', -- Array of feature objects
    reviews JSONB DEFAULT '[]', -- Array of review objects
    categories JSON DEFAULT '[]', -- Array of category strings
    properties JSONB DEFAULT '{}', -- Object with dynamic properties
    analytics JSON DEFAULT '{}', -- Analytics object with nested arrays
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert test data for json_users
INSERT INTO json_users (name, email, settings, metadata, tags, scores, attributes) VALUES
    ('Alice Johnson', 'alice.j@example.com',
     '{"theme": "dark", "notifications": {"email": true, "sms": false}, "permissions": ["read", "write"]}',
     '{"department": "Engineering", "level": "senior", "projects": ["project-a", "project-b"]}',
     '["developer", "team-lead", "mentor"]',
     '[85, 92, 88, 95]',
     '[{"skill": "JavaScript", "level": 9}, {"skill": "Python", "level": 8}, {"skill": "Go", "level": 7}]'),
    
    ('Bob Smith', 'bob.s@example.com',
     '{"theme": "light", "notifications": {"email": false, "sms": true}, "permissions": ["read"]}',
     '{"department": "Marketing", "level": "junior", "projects": ["campaign-x", "campaign-y"]}',
     '["marketer", "analyst"]',
     '[75, 80, 78, 82]',
     '[{"skill": "Analytics", "level": 8}, {"skill": "Design", "level": 6}]'),
    
    ('Carol Davis', 'carol.d@example.com',
     '{"theme": "auto", "notifications": {"email": true, "sms": true}, "permissions": ["read", "write", "admin"]}',
     '{"department": "Sales", "level": "manager", "projects": ["deal-1", "deal-2", "deal-3"]}',
     '["sales", "manager", "closer"]',
     '[90, 94, 91, 89, 93]',
     '[{"skill": "Negotiation", "level": 9}, {"skill": "Communication", "level": 10}, {"skill": "Leadership", "level": 8}]'),
    
    ('David Wilson', 'david.w@example.com',
     '{"theme": "dark", "notifications": {"email": true, "sms": false}, "permissions": ["read", "write"]}',
     '{"department": "Engineering", "level": "senior", "projects": ["infrastructure", "security"]}',
     '["developer", "architect", "security"]',
     '[88, 91, 87, 90]',
     '[{"skill": "DevOps", "level": 9}, {"skill": "Security", "level": 8}, {"skill": "Architecture", "level": 9}]'),
    
    ('Eva Brown', 'eva.b@example.com',
     '{"theme": "light", "notifications": {"email": false, "sms": false}, "permissions": ["read"]}',
     '{"department": "Design", "level": "mid", "projects": ["ui-redesign"]}',
     '["designer", "ux"]',
     '[83, 86, 84]',
     '[{"skill": "UI Design", "level": 8}, {"skill": "UX Research", "level": 7}]');

-- Insert test data for json_products
INSERT INTO json_products (name, features, reviews, categories, properties, analytics) VALUES
    ('Smart Phone X1',
     '[
        {"name": "Camera", "specs": {"megapixels": 108, "features": ["night-mode", "portrait"]}},
        {"name": "Display", "specs": {"size": 6.7, "resolution": "4K"}},
        {"name": "Battery", "specs": {"capacity": 5000, "fast_charge": true}}
     ]',
     '[
        {"rating": 5, "comment": "Excellent phone!", "author": "user1", "helpful": 15},
        {"rating": 4, "comment": "Good camera quality", "author": "user2", "helpful": 8},
        {"rating": 5, "comment": "Long battery life", "author": "user3", "helpful": 12}
     ]',
     '["electronics", "smartphones", "mobile", "premium"]',
     '{"brand": "TechCorp", "warranty": 24, "colors": ["black", "white", "blue"], "price": 999.99}',
     '{"views": [1500, 1200, 1800, 2100], "sales": {"daily": [15, 23, 18, 31], "weekly": [89, 102, 95]}}'),
    
    ('Laptop Pro 15',
     '[
        {"name": "Processor", "specs": {"model": "Intel i9", "cores": 8, "threads": 16}},
        {"name": "Memory", "specs": {"ram": 32, "storage": {"type": "SSD", "capacity": 1024}}},
        {"name": "Graphics", "specs": {"gpu": "RTX 4080", "vram": 16}}
     ]',
     '[
        {"rating": 5, "comment": "Perfect for development", "author": "dev1", "helpful": 25},
        {"rating": 4, "comment": "Great performance", "author": "user4", "helpful": 18},
        {"rating": 5, "comment": "Excellent build quality", "author": "pro1", "helpful": 22}
     ]',
     '["electronics", "computers", "laptops", "professional"]',
     '{"brand": "CompuTech", "warranty": 36, "colors": ["silver", "space-gray"], "price": 2499.99}',
     '{"views": [2200, 1800, 2500, 2800], "sales": {"daily": [8, 12, 6, 15], "weekly": [45, 52, 38]}}'),
    
    ('Coffee Maker Deluxe',
     '[
        {"name": "Brewing", "specs": {"capacity": "12 cups", "types": ["espresso", "americano", "latte"]}},
        {"name": "Smart Features", "specs": {"wifi": true, "app_control": true, "scheduling": true}},
        {"name": "Design", "specs": {"material": "stainless steel", "color": "black"}}
     ]',
     '[
        {"rating": 4, "comment": "Great coffee taste", "author": "coffee1", "helpful": 7},
        {"rating": 3, "comment": "App is buggy", "author": "user5", "helpful": 4},
        {"rating": 5, "comment": "Love the scheduling feature", "author": "morning1", "helpful": 11}
     ]',
     '["home", "kitchen", "appliances", "smart"]',
     '{"brand": "BrewMaster", "warranty": 12, "colors": ["black", "white"], "price": 299.99}',
     '{"views": [850, 920, 780, 1100], "sales": {"daily": [5, 8, 3, 12], "weekly": [28, 31, 24]}}'),
    
    ('Gaming Headset Pro',
     '[
        {"name": "Audio", "specs": {"drivers": "50mm", "frequency": "20Hz-20kHz", "surround": "7.1"}},
        {"name": "Microphone", "specs": {"type": "noise-canceling", "retractable": true}},
        {"name": "Comfort", "specs": {"padding": "memory foam", "weight": "320g"}}
     ]',
     '[
        {"rating": 5, "comment": "Amazing sound quality", "author": "gamer1", "helpful": 19},
        {"rating": 4, "comment": "Comfortable for long sessions", "author": "streamer1", "helpful": 14},
        {"rating": 5, "comment": "Crystal clear mic", "author": "pro_gamer", "helpful": 16}
     ]',
     '["gaming", "audio", "accessories", "pro"]',
     '{"brand": "GameAudio", "warranty": 24, "colors": ["black", "red"], "price": 199.99}',
     '{"views": [1100, 1350, 1200, 1450], "sales": {"daily": [18, 22, 16, 25], "weekly": [78, 85, 72]}}');

-- Add some test data with empty JSON arrays and nulls for edge cases
INSERT INTO json_users (name, email, settings, metadata, tags, scores, attributes) VALUES
    ('Test Empty', 'empty@example.com',
     '{"theme": "dark"}',
     '{"department": "Test"}',
     '[]',
     '[]',
     '[]'),
    ('Test Null', 'null@example.com',
     '{}',
     '{}',
     'null',
     'null',
     'null');

INSERT INTO json_products (name, features, reviews, categories, properties, analytics) VALUES
    ('Empty Product', '[]', '[]', '[]', '{}', '{}'),
    ('Null Product', 'null', 'null', 'null', '{}', '{}');
