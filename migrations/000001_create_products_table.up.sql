CREATE TABLE IF NOT EXISTS products (
  id serial PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  price DOUBLE PRECISION NOT NULL,
  brand TEXT,
  promotion_price DOUBLE PRECISION NULL
);