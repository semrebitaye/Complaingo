-- create role table
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- add role_id to users table if not exists
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS role_id INT REFERENCES roles(id);

-- insert default role
INSERT INTO roles (name)
VALUES ('admin'), ('user')
ON CONFLICT DO NOTHING;