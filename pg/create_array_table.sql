-- Create products table with array types for testing LoadTableSchema
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    tags TEXT[], -- Array of text
    scores INTEGER[], -- Array of integers
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert some test data with arrays
INSERT INTO products (name, tags, scores, metadata) VALUES
    ('Product A', ARRAY['electronics', 'gadgets'], ARRAY[95, 87, 92], '{"brand": "TechCorp", "category": "electronics"}'),
    ('Product B', ARRAY['books', 'fiction'], ARRAY[85, 90], '{"author": "John Author", "genre": "fiction"}'),
    ('Product C', ARRAY['clothing', 'apparel'], ARRAY[88, 91, 89], '{"size": "M", "color": "blue"}');
