-- remove role_id column on users table
ALTER TABLE users DROP COLUMN IF EXISTS role_id;

-- drop roles table
DROP TABLE IF EXISTS roles; 
