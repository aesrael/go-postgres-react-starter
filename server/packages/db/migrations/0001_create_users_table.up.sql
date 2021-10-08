CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    name TEXT NOT NULL,
    password TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    created_on TIMESTAMP NOT NULL default current_timestamp,
    updated_at TIMESTAMP NOT NULL default current_timestamp
);