CREATE TABLE IF NOT EXISTS users(
   id serial PRIMARY KEY,
   first_name VARCHAR (50) NOT NULL,
   last_name VARCHAR (50) NOT NULL,
   email VARCHAR (300) UNIQUE NOT NULL,
   password VARCHAR (255) NOT NULL
);