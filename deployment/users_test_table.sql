CREATE TABLE users (
  id serial PRIMARY KEY,
  email VARCHAR ( 255 ) UNIQUE NOT NULL,
  username VARCHAR ( 50 ) UNIQUE NOT NULL,
  password VARCHAR ( 50 ) NOT NULL,
  created_at TIMESTAMP NOT NULL
);