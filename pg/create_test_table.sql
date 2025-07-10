-- Create users table for testing LoadTableSchema
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    age INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE
);

-- Insert some test data
INSERT INTO users (name, email, age, is_active) VALUES
    ('John Doe', 'john@example.com', 30, TRUE),
    ('Jane Smith', 'jane@example.com', 25, TRUE),
    ('Bob Johnson', 'bob@example.com', 35, FALSE);
